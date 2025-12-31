package cache

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) (string, bool, error)
	Set(ctx context.Context, key string, val string, ttl time.Duration) error
	Del(ctx context.Context, key string) error

	// Counter ops
	Incr(ctx context.Context, key string, by int64, ttl time.Duration) (int64, error)

	// Pubsub (optional)
	Publish(ctx context.Context, channel, msg string) error
}
