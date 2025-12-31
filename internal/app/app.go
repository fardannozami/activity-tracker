package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/fardannozami/activity-tracker/internal/app/httpapi"
	"github.com/fardannozami/activity-tracker/internal/app/httpapi/handler"
	"github.com/fardannozami/activity-tracker/internal/app/worker"
	"github.com/fardannozami/activity-tracker/internal/repo/postgres"
	"github.com/fardannozami/activity-tracker/internal/usecase"
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
	// repos
	clientRepo := postgres.NewClientRepo(a.db)
	apiHitRepo := postgres.NewApiHitRepo(a.db)

	// worker
	batcher := worker.NewBatcher(apiHitRepo)
	go batcher.Run(ctx)

	// usecase
	clientUC := usecase.NewRegisterClientUC(clientRepo)
	recordUC := &usecase.RecordHitUC{
		EnqueueFn: func(hit usecase.HitIn) {
			batcher.Enqueue(hit)
		},
	}

	// handler
	clientHandler := handler.NewClientHandler(clientUC)
	logHandler := handler.NewLogHandler(recordUC)
	r := httpapi.NewRouter(httpapi.Dependency{ClientHandler: clientHandler, LogHandler: logHandler, ClientsRepo: clientRepo})

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
