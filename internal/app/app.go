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
	postgresinfra "github.com/fardannozami/activity-tracker/internal/infra/postgres"
	redisinfra "github.com/fardannozami/activity-tracker/internal/infra/redis"
	"github.com/fardannozami/activity-tracker/internal/repo/cache"
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

	db, err := postgresinfra.NewPool(ctx, cfg.DatabaseUrl)
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
	// Redis optional -> fallback memory
	mem := cache.NewMemoryCache()
	var ca cache.Cache = mem

	rdb := redisinfra.New(a.cfg.RedisAddr)
	if err := redisinfra.Ping(ctx, rdb); err == nil {
		ca = cache.NewRedis(rdb)
		log.Println("redis: connected")
	} else {
		log.Println("redis: down, using memory fallback")
	}

	// repos
	clientRepo := postgres.NewClientRepo(a.db)
	apiHitRepo := postgres.NewApiHitRepo(a.db)
	usageRepo := postgres.NewUsageRepo(a.db)

	// worker
	batcher := worker.NewBatcher(apiHitRepo, usageRepo)
	go batcher.Run(ctx)

	// usecase
	clientUC := usecase.NewRegisterClientUC(clientRepo)
	recordUC := &usecase.RecordHitUC{
		EnqueueFn: func(hit usecase.HitIn) error {
			return batcher.Enqueue(hit)
		},
		Cache: ca,
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
