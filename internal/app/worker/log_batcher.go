package worker

import (
	"context"
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

func NewBatcher(apiHitRepo *postgres.ApiHitRepo) *Batcher {
	return &Batcher{
		apiHitRepo: apiHitRepo,
		in:         make(chan usecase.HitIn, 5000),
		batchSize:  200,
		flushEvery: 1 * time.Second,
	}
}

func (b *Batcher) Enqueue(hit usecase.HitIn) {
	select {
	case b.in <- hit:
	default:
		// channel full: drop OR block depending requirement.
		// for test, better block a tiny bit:
		b.in <- hit
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
		_ = b.apiHitRepo.BulkInsert(ctx, rows)

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
