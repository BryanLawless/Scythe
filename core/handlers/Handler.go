package handlers

import (
	"Scythe/core/common"
	"Scythe/core/utility"
	"context"
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

func (h *Handler) HandleCommands(ctx context.Context) {
	action := utility.StartActionPrompt()

	var err error
	switch strings.ToLower(action) {
	case "browse":
		err = h.BrowseHandler(ctx)
	case "download":
		err = h.DownloadHandler(ctx)
	case "exit":
		os.Exit(0)
	}

	if err != nil {
		utility.WriteToConsole(err.Error(), "error")
	}
}
