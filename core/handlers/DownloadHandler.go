package handlers

import (
	"Scythe/core/common"
	"Scythe/core/modules"
	"Scythe/core/modules/download"
	"Scythe/core/utility"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func (h *Handler) DownloadHandler(ctx context.Context) error {
	var err error

	if len(h.Media) <= 0 {
		utility.WriteToConsole("You must add URLs to download first!", "error")
		return nil
	}

	mediaResults := utility.MediaToResults(h.Media)
	confirmed := utility.ConfirmSelectDownloadsPrompt(mediaResults)

	confirmedMedia := []common.Media{}
	for _, confirm := range confirmed {
		possibleIndex := strings.Split(confirm, "-")[0]

		mediaIndex, err := strconv.Atoi(possibleIndex)
		if err != nil {
			return err
		}

		confirmedMedia = append(confirmedMedia, h.Media[mediaIndex])
	}

	targets, err := modules.Prepare(h.Context, &modules.PreCheck{
		Media:   confirmedMedia,
		Timeout: time.Duration(h.Timeout) * time.Second,
	})

	if err != nil {
		return err
	}

	save, err := h.CreateSaveDirectory()
	if err != nil {
		return err
	}

	for _, target := range targets {
		utility.WriteToConsole(fmt.Sprintf("Downloading %s", target.Filename), "info")
		switch target.Status {
		case "success":
			utility.WriteToConsole(fmt.Sprintf("Utilizing %d download threads", h.Workers), "success")
			err = download.MultiDownload(h.Context, &download.Download{
				URL:           target.URL,
				Workers:       h.Workers,
				Filename:      target.Filename,
				Directory:     save,
				ContentLength: target.ContentLength,
			})
		case "no_range":
			utility.WriteToConsole("Using one download thread", "warning")
			err = download.SingleDownload(h.Context, &download.Download{
				URL:           target.URL,
				Filename:      target.Filename,
				Directory:     save,
				ContentLength: target.ContentLength,
			})
		default:
			utility.WriteToConsole(fmt.Sprintf("Failed to download from %s", target.URL), "error")
		}
	}

	return err
}
