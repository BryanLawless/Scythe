package app

import (
	"Scythe/core/handlers"
	"Scythe/core/utility"
	"context"
	"fmt"
)

func Start(ctx context.Context) error {
	fmt.Println(utility.GetBanner())
	utility.WriteToConsole("Starting Scythe: "+utility.GetQuote(), "success")

	handler := handlers.New(ctx)
	for {
		handler.HandleCommands(ctx)
	}
}
