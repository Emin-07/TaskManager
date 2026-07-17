package redis

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/Emin-07/TaskManager/internal/core/domain"
)

func (rs *RedisClientRepo) Increment(ctx context.Context, name, id, userId string) (int, error) {
	val, err := rs.Get(ctx, name, id, userId)
	if errors.Is(err, domain.ErrKeyNotFound) {
		return 1, rs.Set(ctx, name, id, userId, 1, time.Minute)
	} else if err != nil {
		return 0, err
	}

	valN, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}
	valN++
	err = rs.Set(ctx, name, id, userId, valN, time.Minute)
	if err != nil {
		return 0, err
	}
	return valN, nil
}

func (rs *RedisClientRepo) DecrementBy(ctx context.Context, name, id, userId string, n int, timeToLive time.Duration) (int, error) {
	val, err := rs.Get(ctx, name, id, userId)
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
		rs.Rdb.Del(ctx, keyMaker(name, id, userId))
		return 0, nil
	}
	valN -= n
	err = rs.Set(ctx, name, id, userId, valN, timeToLive)
	if err != nil {
		return 0, err
	}
	return valN, nil
}
