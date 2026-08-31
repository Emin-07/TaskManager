package redis

import (
	"context"
	"sync"
	"time"
)

// rateLimitMu serializes Increment/Decrement so the read-modify-write of the
// Redis counter plus its TTL assignment stays atomic across goroutines.
var rateLimitMu sync.Mutex

func (rs *RedisClientRepo) Increment(ctx context.Context, id string) (int, error) {
	key := userKeyMaker(id)

	rateLimitMu.Lock()
	defer rateLimitMu.Unlock()

	val, err := rs.Rdb.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	if val == 1 {
		rs.Rdb.Expire(ctx, key, time.Minute)
	}

	return int(val), nil
}

func (rs *RedisClientRepo) Decrement(ctx context.Context, id string) (int, error) {
	key := userKeyMaker(id)

	rateLimitMu.Lock()
	defer rateLimitMu.Unlock()

	val, err := rs.Rdb.Decr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return int(val), nil
}
