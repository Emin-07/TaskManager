package redis

import (
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClientRepo struct {
	Rdb *redis.Client
}

func NewRedisClientRepo() *RedisClientRepo {
	return &RedisClientRepo{redis.NewClient(&redis.Options{
		Addr:               os.Getenv("REDIS_ADDR"),
		Password:           os.Getenv("REDIS_PASSWORD"),
		DB:                 0,
		DialerRetries:      5,
		DialerRetryTimeout: 100 * time.Millisecond,
	})}
}

func taskKeyMaker(name, id, userId string) string {
	if userId == "" {
		return fmt.Sprintf("%s:%s", name, id)
	}
	return fmt.Sprintf("%s:%s-%s", name, id, userId)
}

func userKeyMaker(id string) string {
	return fmt.Sprintf("user:%v", id)
}
