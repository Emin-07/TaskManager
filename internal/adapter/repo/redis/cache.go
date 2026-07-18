package redis

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/Emin-07/TaskManager/internal/core/domain"
)

func (rs *RedisClientRepo) Set(ctx context.Context, name, id, userId string, val any, duration time.Duration) error {
	return rs.Rdb.Set(ctx, keyMaker(name, id, userId), val, duration).Err()
}

func (rs *RedisClientRepo) Get(ctx context.Context, name, id, userId string) (string, error) {
	val, err := rs.Rdb.Get(ctx, keyMaker(name, id, userId)).Result()
	if errors.Is(err, redis.Nil) {
		return "", domain.ErrKeyNotFound
	} else if err != nil {
		return "", err
	}
	return val, nil
}

func (rs *RedisClientRepo) Del(ctx context.Context, name, id, userId string) error {
	return rs.Rdb.Del(ctx, keyMaker(name, id, userId)).Err()
}
