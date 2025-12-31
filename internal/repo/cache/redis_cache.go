package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct{ rdb *redis.Client }

func NewRedis(rdb *redis.Client) *RedisCache { return &RedisCache{rdb: rdb} }

func (c *RedisCache) Get(ctx context.Context, key string) (string, bool, error) {
	val, err := c.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return val, true, nil
}

func (c *RedisCache) Set(ctx context.Context, key string, val string, ttl time.Duration) error {
	return c.rdb.Set(ctx, key, val, ttl).Err()
}

func (c *RedisCache) Del(ctx context.Context, key string) error {
	return c.rdb.Del(ctx, key).Err()
}

func (c *RedisCache) Incr(ctx context.Context, key string, by int64, ttl time.Duration) (int64, error) {
	pipe := c.rdb.TxPipeline()
	incr := pipe.IncrBy(ctx, key, by)
	pipe.Expire(ctx, key, ttl)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}
	return incr.Val(), nil
}

func (c *RedisCache) Publish(ctx context.Context, channel, msg string) error {
	return c.rdb.Publish(ctx, channel, msg).Err()
}
