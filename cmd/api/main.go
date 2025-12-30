package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fardannozami/activity-tracker/internal/app"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := app.LoadConfig()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	_, err = app.NewPostgre(ctx, cfg)
	if err != nil {
		log.Fatalf("init app: %v", err)
	}

	// Graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigCh:
		log.Println("shutdown: signal received")
	case <-ctx.Done():
		log.Println("shutdown: context canceled")
	}

	log.Println("bye 👋")

}
