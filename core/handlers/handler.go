package handlers

import (
	"Scythe/core/common"
	"Scythe/core/utility"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Handler struct {
	Media     []common.Media
	Timeout   int
	Workers   int
	Directory string
	Context   context.Context
}

func New(ctx context.Context) *Handler {
	return &Handler{
		Timeout: 30,
		Context: ctx,
		Workers: runtime.NumCPU(),
	}
}

func (h *Handler) CreateSaveDirectory() (string, error) {
	if len(h.Directory) > 0 {
		location, err := os.Stat(h.Directory)
		if err == nil && location.IsDir() {
			return h.Directory, nil
		} else {
			directory, _ := filepath.Split(h.Directory)
			if len(directory) > 0 {
				if err := os.MkdirAll(directory, 0755); err != nil {
					return directory, nil
				}
			}
		}
	}

	return "", nil
}

func (h *Handler) BrowseHandler(ctx context.Context) error {
	category := utility.MediaCategoryPrompt()

	var err error
	switch strings.ToLower(category) {
	case "videos":
		err = h.VideoHandler(ctx)
	}

	return err
}

func (h *Handler) AddRawHandler(ctx context.Context) error {
	mediaList := []common.Media{}
	media := utility.AddRawPrompt()

	for _, m := range media {
		request, _, err := common.MakeRequest(ctx, common.Request{
			URL:         m,
			Method:      "GET",
			ParseBody:   false,
			RandomAgent: true,
		})

		if err != nil {
			return err
		}

		contentType := request.Response.Header.Get("Content-Type")
		extension := utility.GetFileExtensionFromMime(contentType)
		category := utility.GetCategoryFromMimeType(contentType)
		filename := utility.FilenamePrompt()

		utility.WriteToConsole(fmt.Sprintf(
			"Resource '%s' category %s added to download list.",
			filename, category), "success")

		mediaList = append(mediaList, common.Media{
			URL:       m,
			Name:      filename,
			Provider:  "raw",
			Category:  category,
			Extension: extension,
			SourceURL: m,
		})
	}

	h.Media = append(h.Media, mediaList...)

	return nil
}

func (h *Handler) HandleCommands(ctx context.Context) {
	action := utility.StartActionPrompt()

	var err error
	switch strings.ToLower(action) {
	case "browse":
		err = h.BrowseHandler(ctx)
	case "download":
		err = h.DownloadHandler(ctx)
	case "add raw":
		err = h.AddRawHandler(ctx)
	case "exit":
		utility.WriteToConsole("Exiting...", "info")
		os.Exit(0)
	}

	if err != nil {
		utility.WriteToConsole(err.Error(), "error")
	}
}
