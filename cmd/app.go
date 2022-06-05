package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/hyonosake/HTTP-Multiplexer/internal/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	a, err := app.NewService(ctx)
	if err != nil {
		log.Fatalf("Unable to craete app: %v\n", err)
	}
	a.Run(ctx)

}
