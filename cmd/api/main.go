package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fardannozami/activity-tracker/internal/app"
	"github.com/fardannozami/activity-tracker/internal/infra/migrate"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// load config
	cfg, err := app.LoadConfig()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	// run migration
	if err := migrate.RunMigration(cfg.DatabaseUrl); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	// init db pool
	a, err := app.NewPostgre(ctx, cfg)
	if err != nil {
		log.Fatalf("init app: %v", err)
	}

	// init http server
	go func() {
		if err := a.RunHttp(ctx); err != nil {
			log.Printf("http server:  %v", err)
			cancel()
		}
	}()

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
