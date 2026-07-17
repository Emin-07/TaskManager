package service

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/Emin-07/TaskManager/internal/core/port"
)

type RateAndCacheServ struct {
	repo port.RateAndCacheRepo
}

func NewRateAndCacheService(repo port.RateAndCacheRepo) RateAndCacheServ {
	return RateAndCacheServ{repo: repo}
}

func (rc RateAndCacheServ) Increment(ctx context.Context, name, id, userId string) (int, error) {
	return rc.repo.Increment(ctx, name, id, userId)
}
func (rc RateAndCacheServ) DecrementBy(ctx context.Context, name, id, userId string, n int, timeToLive time.Duration) (int, error) {
	return rc.repo.DecrementBy(ctx, name, id, userId, n, timeToLive)
}

func (rc RateAndCacheServ) Set(ctx context.Context, name, id, userId string, val any, duration time.Duration) error {
	if reflect.ValueOf(val).Kind() == reflect.Struct {
		res, err := json.Marshal(val)
		if err != nil {
			fmt.Println(val)
			fmt.Println(res)
			return err
		}
		return rc.repo.Set(ctx, name, id, userId, res, duration)
	}
	return rc.repo.Set(ctx, name, id, userId, val, duration)
}
func (rc RateAndCacheServ) Get(ctx context.Context, name, id, userId string) (string, error) {
	return rc.repo.Get(ctx, name, id, userId)
}
