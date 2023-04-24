package modules

import (
	"Scythe/core/common"
	"Scythe/core/utility"
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

type PreCheck struct {
	Media   []common.Media
	Timeout time.Duration
}

type Resource struct {
	URL           string
	Status        string
	Filename      string
	ContentLength int64
}

func Prepare(ctx context.Context, p *PreCheck) ([]*Resource, error) {
	ctx, cancel := context.WithTimeout(ctx, p.Timeout)
	defer cancel()

	infos, err := parseResults(ctx, p.Media)
	if err != nil {
		return nil, err
	}

	return infos, nil
}

func parseResults(ctx context.Context, media []common.Media) ([]*Resource, error) {
	var mutable sync.Mutex
	eg, ctx := errgroup.WithContext(ctx)

	results := []*Resource{}
	for _, m := range media {
		m := m
		eg.Go(func() error {
			result, err := resultParser(ctx, m)
			if err != nil {
				return err
			}

			mutable.Lock()
			results = append(results, result)
			mutable.Unlock()

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return results, nil
}

func resultParser(ctx context.Context, m common.Media) (*Resource, error) {
	resource := &Resource{
		URL:    m.SourceURL,
		Status: "success",
	}

	safeFilename := fmt.Sprintf("%s.%s", utility.FilenameSafe(m.Name), m.Extension)

	request, _, err := common.MakeRequest(ctx, common.Request{
		URL:         resource.URL,
		Method:      "GET",
		ParseBody:   false,
		RandomAgent: true,
	})

	if err != nil {
		resource.Status = "error"
		return resource, err
	}

	if request.Response.ContentLength <= 0 {
		resource.Status = "no_content"
		return resource, err
	}

	if request.Response.Header.Get("Accept-Ranges") != "bytes" {
		resource.Status = "no_range"
	}

	possibleUrl := request.Response.Request.URL.String()
	if utility.NotPrev(possibleUrl, resource.URL) {
		resource.Filename = safeFilename
		resource.ContentLength = request.Response.ContentLength

		return resource, nil
	}

	resource.Filename = safeFilename
	resource.ContentLength = request.Response.ContentLength

	return resource, nil
}
