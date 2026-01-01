package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fardannozami/activity-tracker/internal/repo/cache"
)

type RecordHitUC struct {
	Cache     cache.Cache
	EnqueueFn func(hit HitIn) error
}

type HitIn struct {
	ClientID  string
	IP        string
	Endpoint  string
	Timestamp time.Time
}

func (uc *RecordHitUC) Execute(ctx context.Context, in HitIn) error {
	day := in.Timestamp.Format("2006-01-02")
	hour := in.Timestamp.Truncate(time.Hour).Format(time.RFC3339)

	// Redis keys
	dailyKey := fmt.Sprintf("counter:daily:%s:%s", in.ClientID, day)
	hourlyKey := fmt.Sprintf("counter:hourly:%s:%s", in.ClientID, hour)

	// increment counters (ttl enough long)
	_, _ = uc.Cache.Incr(ctx, dailyKey, 1, 72*time.Hour)
	_, _ = uc.Cache.Incr(ctx, hourlyKey, 1, 48*time.Hour)

	// bump version for cache invalidation
	verKey := fmt.Sprintf("usage:ver:%s", in.ClientID)
	_, _ = uc.Cache.Incr(ctx, verKey, 1, 7*24*time.Hour)

	// publish update (optional)
	msg, _ := json.Marshal(map[string]any{
		"client_id": in.ClientID,
		"ts":        time.Now().Format(time.RFC3339),
	})
	_ = uc.Cache.Publish(ctx, "usage.updated", string(msg))

	if uc.EnqueueFn != nil {
		return uc.EnqueueFn(in)
	}

	return nil
}
