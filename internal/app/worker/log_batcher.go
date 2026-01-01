package worker

import (
	"context"
	"errors"
	"time"

	"github.com/fardannozami/activity-tracker/internal/repo/postgres"
	"github.com/fardannozami/activity-tracker/internal/usecase"
)

type Batcher struct {
	hitRepo    *postgres.ApiHitRepo
	usage      postgres.UsageWriter
	in         chan usecase.HitIn
	batchSize  int
	flushEvery time.Duration
}

var ErrQueueFull = errors.New("log queue full")

func NewBatcher(hitRepo *postgres.ApiHitRepo, uw postgres.UsageWriter) *Batcher {
	return &Batcher{
		hitRepo:    hitRepo,
		usage:      uw,
		in:         make(chan usecase.HitIn, 5000),
		batchSize:  200,
		flushEvery: 1 * time.Second,
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
		_ = b.hitRepo.BulkInsert(ctx, rows)

		// 2) aggregated upserts
		if b.usage != nil {
			usageHits := make([]postgres.UsageHit, 0, len(buf))
			for _, h := range buf {
				usageHits = append(usageHits, postgres.UsageHit{
					ClientID:  h.ClientID,
					Timestamp: h.Timestamp,
				})
			}
			_ = b.usage.UpsertAggregates(ctx, usageHits)
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
