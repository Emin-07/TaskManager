package port

import (
	"context"
	"time"
)

type RateAndCacheService interface {
	Increment(ctx context.Context, name, id, userId string) (int, error)
	DecrementBy(ctx context.Context, name, id, userId string, n int, timeToLive time.Duration) (int, error)

	Set(ctx context.Context, name, id, userId string, val any, duration time.Duration) error
	Get(ctx context.Context, name, id, userId string) (string, error)
	Del(ctx context.Context, name, id, userId string) error
}

type RateAndCacheRepo interface {
	Increment(ctx context.Context, name, id, userId string) (int, error)
	DecrementBy(ctx context.Context, name, id, userId string, n int, timeToLive time.Duration) (int, error)

	Set(ctx context.Context, name, id, userId string, val any, duration time.Duration) error
	Get(ctx context.Context, name, id, userId string) (string, error)
	Del(ctx context.Context, name, id, userId string) error
}
