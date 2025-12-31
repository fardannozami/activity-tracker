package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/fardannozami/activity-tracker/internal/httpapi"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	db         *pgxpool.Pool
	cfg        Config
	httpServer *http.Server
}

func NewPostgre(ctx context.Context, cfg Config) (*App, error) {
	config, err := pgxpool.ParseConfig(cfg.DatabaseUrl)
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

	log.Printf(
		"database connected: host=%s db=%s minConns=%d maxConns=%d",
		config.ConnConfig.Host,
		config.ConnConfig.Database,
		config.MinConns,
		config.MaxConns,
	)

	return &App{
		db:  db,
		cfg: cfg,
	}, nil
}

func (a *App) RunHttp(ctx context.Context) error {
	r := httpapi.NewRouter()

	a.httpServer = &http.Server{
		Addr:    a.cfg.HTTPAddr,
		Handler: r,
	}

	log.Printf("listening on %s", a.cfg.HTTPAddr)
	err := a.httpServer.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}
