package app

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	db  *pgxpool.Pool
	cfg Config
}

func NewPostgre(ctx context.Context, cfg Config) (*App, error) {
	config, err := pgxpool.ParseConfig(cfg.DatabseUrl)
	if err != nil {
		return nil, err
	}

	config.MinConns = 1
	config.MaxConns = 5
	config.MaxConnIdleTime = 5 * time.Minute
	config.MaxConnLifetime = 30 * time.Minute
	config.HealthCheckPeriod = 1 * time.Minute

	db, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	ctxPing, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := db.Ping(ctxPing); err != nil {
		return nil, err
	}

	return &App{
		db: db,
	}, nil
}
