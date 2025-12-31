package usecase

import (
	"context"
	"time"
)

type RecordHitUC struct {
	// Cache     cache.MemoryCache
	EnqueueFn func(hit HitIn)
}

type HitIn struct {
	ClientID  string
	IP        string
	Endpoint  string
	Timestamp time.Time
}

func (uc *RecordHitUC) Execute(ctx context.Context, in HitIn) error {
	if uc.EnqueueFn != nil {
		uc.EnqueueFn(in)
	}

	return nil
}
