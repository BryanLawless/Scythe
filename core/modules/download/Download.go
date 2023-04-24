package download

import (
	"Scythe/core/common"
	"Scythe/core/modules"
	"context"
	"fmt"
	"net/http"
	"path/filepath"
)

type Task struct {
	ID        int
	URL       string
	Range     modules.Range
	Workers   int
	Filename  string
	Directory string
}

type OperateTasks struct {
	URL           string
	Size          int64
	Workers       int
	Filename      string
	Directory     string
	ContentLength int64
}

type Download struct {
	URL           string
	Workers       int
	Filename      string
	Directory     string
	ContentLength int64
}

type ParallelDownloadConfig struct {
	Tasks         []*Task
	Directory     string
	ContentLength int64
}

func (t *Task) destinationPath() string {
	return filepath.Join(t.Directory, fmt.Sprintf("%s.%d.%d", t.Filename, t.Workers, t.ID))
}

func (t *Task) probeRequest(ctx context.Context) (*http.Request, error) {
	request, _, err := common.MakeRequest(ctx, common.Request{
		URL:     t.URL,
		Method:  "GET",
		Headers: map[string]string{"Range": t.Range.BytesRange()},
		Partial: true,
	})

	if err != nil {
		return nil, err
	}

	return request.Request, nil
}
