package download

import (
	"Scythe/core/common"
	"context"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

func SingleDownload(ctx context.Context, c *Download) error {
	request, _, err := common.MakeRequest(ctx, common.Request{
		URL:         c.URL,
		Method:      "GET",
		ParseBody:   false,
		RandomAgent: true,
	})

	if err != nil {
		return err
	}

	output, err := os.OpenFile(c.Filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	defer output.Close()

	bar := pb.Start64(request.Response.ContentLength).Set(pb.Bytes, true)

	defer bar.Finish()

	//bar.SetCurrent(0)

	rd := bar.NewProxyReader(request.Response.Body)
	if _, err := io.Copy(output, rd); err != nil {
		return err
	}

	return nil
}
