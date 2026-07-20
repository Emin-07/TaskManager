package redis

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/Emin-07/TaskManager/internal/core/domain"
)

func (rs *RedisClientRepo) Increment(ctx context.Context, id string) (int, error) {
	key := userKeyMaker(id)
	val, err := rs.Rdb.Get(ctx, key).Result()
	if errors.Is(err, domain.ErrKeyNotFound) {
		return 1, rs.Rdb.Set(ctx, key, 1, time.Minute).Err()
	} else if err != nil {
		return 0, err
	}

	valN, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}
	valN++
	err = rs.Rdb.Set(ctx, id, valN, time.Minute).Err()
	if err != nil {
		return 0, err
	}
	return valN, nil
}

func (rs *RedisClientRepo) DecrementBy(ctx context.Context, id string, n int, timeToLive time.Duration) (int, error) {
	key := userKeyMaker(id)
	val, err := rs.Rdb.Get(ctx, key).Result()
	if errors.Is(err, domain.ErrKeyNotFound) {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	valN, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}
	if valN == 0 || valN <= n {
		rs.Rdb.Del(ctx, key)
		return 0, nil
	}
	valN -= n
	err = rs.Rdb.Set(ctx, key, val, timeToLive).Err()
	if err != nil {
		return 0, err
	}
	return valN, nil
}
