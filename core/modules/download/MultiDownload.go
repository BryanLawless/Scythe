package download

import (
	"Scythe/core/common"
	"Scythe/core/modules"
	"Scythe/core/utility"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb/v3"
	"golang.org/x/sync/errgroup"
)

func StartMultiDownload(ctx context.Context, c *DownloadConfig) error {
	directory := utility.PartialDirectory(c.Directory, c.Filename, c.Workers)

	if err := os.MkdirAll(directory, 0755); err != nil {
		return err
	}

	tasks := createTasks(OperateTasks{
		URL:           c.URL,
		Size:          c.ContentLength / int64(c.Workers),
		Workers:       c.Workers,
		Filename:      c.Filename,
		Directory:     directory,
		ContentLength: c.ContentLength,
	})

	if err := parallelDownload(ctx, &ParallelDownloadConfig{
		Tasks:         tasks,
		Directory:     directory,
		ContentLength: c.ContentLength,
	}); err != nil {
		return err
	}

	return bindChunks(c, directory)
}

func createTasks(c OperateTasks) []*Task {
	tasks := make([]*Task, 0, c.Workers)

	var totalWorkers int
	for i := 0; i < c.Workers; i++ {
		currentRange := modules.MakeRange(i, c.Workers, c.Size, c.ContentLength)

		part := filepath.Join(c.Directory, fmt.Sprintf("%s.%d.%d", c.Filename, c.Workers, i))

		if info, err := os.Stat(part); err == nil {
			infoSize := info.Size()

			if i == c.Workers-1 {
				if infoSize == currentRange.Max-currentRange.Min {
					continue
				}
			} else if infoSize == c.Size {
				continue
			}

			currentRange.Min += infoSize
		}

		tasks = append(tasks, &Task{
			ID:        i,
			URL:       c.URL,
			Range:     currentRange,
			Workers:   c.Workers,
			Filename:  c.Filename,
			Directory: c.Directory,
		})

		totalWorkers++
	}

	return tasks
}

func parallelDownload(ctx context.Context, c *ParallelDownloadConfig) error {
	eg, ctx := errgroup.WithContext(ctx)
	bar := pb.Start64(c.ContentLength).Set(pb.Bytes, true)

	defer bar.Finish()

	size, err := utility.DirectorySize(c.Directory)
	if err != nil {
		return err
	}

	bar.SetCurrent(size)

	for _, task := range c.Tasks {
		task := task

		eg.Go(func() error {
			req, err := task.probeRequest(ctx)
			if err != nil {
				return err
			}

			return task.download(ctx, req, bar)
		})
	}

	return eg.Wait()
}

func (t *Task) download(ctx context.Context, req *http.Request, bar *pb.ProgressBar) error {
	request, _, err := common.MakeRequest(ctx, common.Request{
		ParseBody:           false,
		RandomAgent:         true,
		ContinueFromRequest: true,
		ResumeFromRequest:   req,
	})

	if err != nil {
		return err
	}

	output, err := os.OpenFile(t.destinationPath(), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	defer output.Close()

	rd := bar.NewProxyReader(request.Response.Body)
	if _, err := io.Copy(output, rd); err != nil {
		return err
	}

	return nil
}

func bindChunks(c *DownloadConfig, partialDir string) error {
	destination := filepath.Join(c.Directory, c.Filename)
	file, err := os.Create(destination)
	if err != nil {
		return err
	}

	defer file.Close()

	bar := pb.Start64(c.ContentLength)

	copier := func(name string) error {
		chunk, err := os.Open(name)
		if err != nil {
			return err
		}

		defer chunk.Close()

		proxy := bar.NewProxyReader(chunk)
		if _, err := io.Copy(file, proxy); err != nil {
			return err
		}

		if err := os.Remove(name); err != nil {
			return err
		}

		return nil
	}

	for i := 0; i < c.Workers; i++ {
		name := fmt.Sprintf("%s/%s.%d.%d", partialDir, c.Filename, c.Workers, i)
		if err := copier(name); err != nil {
			return err
		}
	}

	bar.Finish()

	if err := os.RemoveAll(partialDir); err != nil {
		return err
	}

	return nil
}
