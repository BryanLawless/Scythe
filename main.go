package main

import (
	"Scythe/core/app"
	"Scythe/core/utility"
	"context"
	"math/rand"
	"os"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	if err := app.Start(context.Background()); err != nil {
		utility.WriteToConsole(err.Error(), "error")
		os.Exit(1)
	}
}
