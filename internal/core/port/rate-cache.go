package port

import (
	"context"
	"time"
)

type RateAndCacheRepo interface {
	Increment(ctx context.Context, id string) (int, error)
	Decrement(ctx context.Context, id string) (int, error)
	Set(ctx context.Context, name, id, userId string, val any, duration time.Duration) error
	Get(ctx context.Context, name, id, userId string) (string, error)
	Del(ctx context.Context, name, id, userId string) error
}

type RateAndCacheService interface {
	Increment(ctx context.Context, id string) (int, error)
	Decrement(ctx context.Context, id string) (int, error)

	Set(ctx context.Context, name, id, userId string, val any, duration time.Duration) error
	Get(ctx context.Context, name, id, userId string) (string, error)
	Del(ctx context.Context, name, id, userId string) error
}
