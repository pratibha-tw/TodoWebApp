package redisclient

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Del(ctx context.Context, key string) *redis.IntCmd
}

type redisClientImpl struct {
	redisClient *redis.Client
}

func (r redisClientImpl) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return r.redisClient.Set(ctx, key, value, expiration)
}

func (r redisClientImpl) Get(ctx context.Context, key string) *redis.StringCmd {
	return r.redisClient.Get(ctx, key)
}
func (r redisClientImpl) Del(ctx context.Context, key string) *redis.IntCmd {
	return r.redisClient.Del(ctx, key)
}

func NewRedisClient(rc *redis.Client) RedisClient {
	return &redisClientImpl{rc}
}
