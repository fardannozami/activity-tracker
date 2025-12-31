package worker

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/fardannozami/activity-tracker/internal/repo/postgres"
	"github.com/fardannozami/activity-tracker/internal/usecase"
)

type Batcher struct {
	apiHitRepo *postgres.ApiHitRepo
	in         chan usecase.HitIn
	batchSize  int
	flushEvery time.Duration
}

var ErrQueueFull = errors.New("log queue full")

func NewBatcher(apiHitRepo *postgres.ApiHitRepo) *Batcher {
	return &Batcher{
		apiHitRepo: apiHitRepo,
		in:         make(chan usecase.HitIn, 5000),
		batchSize:  200,
		flushEvery: 30 * time.Second,
	}
}

func (b *Batcher) Enqueue(hit usecase.HitIn) error {
	select {
	case b.in <- hit:
		return nil
	default:
		return ErrQueueFull
	}
}

func (b *Batcher) Run(ctx context.Context) {
	ticker := time.NewTicker(b.flushEvery)
	defer ticker.Stop()

	buf := make([]usecase.HitIn, 0, b.batchSize)

	flush := func() {
		if len(buf) == 0 {
			return
		}

		rows := make([]struct {
			ClientID string
			IP       string
			Endpoint string
			TS       string
		}, 0, len(buf))

		for _, h := range buf {
			rows = append(rows, struct {
				ClientID string
				IP       string
				Endpoint string
				TS       string
			}{
				ClientID: h.ClientID,
				IP:       h.IP,
				Endpoint: h.Endpoint,
				TS:       h.Timestamp.Format(time.RFC3339),
			})
		}

		// 1) raw insert
		if err := b.apiHitRepo.BulkInsert(ctx, rows); err != nil {
			log.Printf("bulk insert api_hits failed: %v", err)
		}

		buf = buf[:0]
	}

	for {
		select {
		case <-ctx.Done():
			flush()
			return
		case h := <-b.in:
			buf = append(buf, h)
			if len(buf) >= b.batchSize {
				flush()
			}
		case <-ticker.C:
			flush()
		}
	}
}
