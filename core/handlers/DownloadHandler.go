package handlers

import (
	"Scythe/core/modules"
	"Scythe/core/modules/download"
	"Scythe/core/utility"
	"context"
	"fmt"
	"time"
)

func (h *Handler) DownloadHandler(ctx context.Context) error {
	if len(h.URLs) <= 0 {
		utility.WriteToConsole("You must add URLs to download first!", "error")
		return nil
	}

	targets, err := modules.Check(h.Context, &modules.PreCheck{
		URLs:    h.URLs,
		Timeout: time.Duration(h.Timeout) * time.Second,
	})

	if err != nil {
		return err
	}

	/*save, err := h.CreateSaveDirectory()
	if err != nil {
		return err
	}*/

	for _, target := range targets {
		switch target.Status {
		case "success":
			/*utility.WriteToConsole(fmt.Sprintf("Starting multi-threaded download: %s", target.Filename), "success")
			download.StartMultiDownload(h.Context, &download.DownloadConfig{
				URL:           target.URL,
				Workers:       h.Workers,
				Filename:      target.Filename,
				Directory:     save,
				ContentLength: target.ContentLength,
			})*/

			utility.WriteToConsole(fmt.Sprintf("Starting single-threaded download: %s", target.Filename), "warning")
			err := download.SingleDownload(h.Context, &download.DownloadConfig{
				URL:           target.URL,
				Filename:      "heheh.mp4",
				Directory:     "heheheheheheheh",
				ContentLength: target.ContentLength,
			})

			if err != nil {
				return err
			}
		case "no_range":
			utility.WriteToConsole(fmt.Sprintf("Starting single-threaded download: %s", target.Filename), "warning")
			err := download.SingleDownload(h.Context, &download.DownloadConfig{
				URL:           target.URL,
				Filename:      "heheh.mp4",
				Directory:     "heheheheheheh",
				ContentLength: target.ContentLength,
			})

			if err != nil {
				return err
			}

		default:
			utility.WriteToConsole(fmt.Sprintf("Failed to download from %s", target.URL), "error")
		}
	}

	return nil
}
