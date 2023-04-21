package handlers

import (
	"Scythe/core/utility"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Handler struct {
	URLs          []string
	Timeout       int
	Workers       int
	SaveDirectory string
	Context       context.Context
}

func New(ctx context.Context) *Handler {
	return &Handler{
		Timeout: 10,
		Context: ctx,
		Workers: runtime.NumCPU(),
	}
}

func (h *Handler) CreateSaveDirectory() (string, error) {
	if len(h.SaveDirectory) > 0 {
		location, err := os.Stat(h.SaveDirectory)
		if err == nil && location.IsDir() {
			return h.SaveDirectory, nil
		} else {
			directory, _ := filepath.Split(h.SaveDirectory)
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
	case "pictures":
		fmt.Println("pictures")
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
