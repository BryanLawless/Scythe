package modules

import (
	"Scythe/core/common"
	"Scythe/core/utility"
	"context"
	"path"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

type PreCheck struct {
	URLs    []string
	Timeout time.Duration
}

type Mirror struct {
	URL           string
	Status        string
	Filename      string
	ContentLength int64
}

func Check(ctx context.Context, p *PreCheck) ([]*Mirror, error) {
	ctx, cancel := context.WithTimeout(ctx, p.Timeout)
	defer cancel()

	infos, err := startAnalyzer(ctx, p.URLs)
	if err != nil {
		return nil, err
	}

	return infos, nil
}

func startAnalyzer(ctx context.Context, urls []string) ([]*Mirror, error) {
	var mutable sync.Mutex
	eg, ctx := errgroup.WithContext(ctx)

	infos := make([]*Mirror, 0, len(urls))

	for _, url := range urls {
		url := url
		eg.Go(func() error {
			info, err := analyze(ctx, url)
			if err != nil {
				return err
			}

			mutable.Lock()
			infos = append(infos, info)
			mutable.Unlock()

			return nil
		})
	}

	eg.Wait()

	return infos, nil
}

func analyze(ctx context.Context, url string) (*Mirror, error) {
	request, _, err := common.MakeRequest(ctx, common.Request{
		URL:    url,
		Method: "HEAD",
	})

	if err != nil {
		return &Mirror{Status: "error"}, err
	}

	/*if request.Response.Header.Get("Accept-Ranges") != "bytes" {
		return &Mirror{Status: "no_range"}, err
	}*/

	if request.Response.ContentLength <= 0 {
		return &Mirror{Status: "no_content"}, err
	}

	possibleUrl := request.Response.Request.URL.String()
	if utility.NotPrev(possibleUrl, url) {
		return &Mirror{
			URL:           possibleUrl,
			Status:        "success",
			Filename:      path.Base(possibleUrl),
			ContentLength: request.Response.ContentLength,
		}, nil
	}

	return &Mirror{
		URL:           url,
		Status:        "success",
		Filename:      path.Base(url),
		ContentLength: request.Response.ContentLength,
	}, nil
}
