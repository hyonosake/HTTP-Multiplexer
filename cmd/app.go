package main

import (
	"context"
	"log"

	"github.com/hyonosake/HTTP-Multiplexer/internal/server"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s, err := server.New(ctx)
	if err != nil {
		log.Fatalf("Unable to craete server: %v\n", err)
	}
	s.Run(ctx)

}
