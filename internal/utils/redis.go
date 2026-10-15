package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func GetFromRedis[T any](ctx context.Context,rc *redis.Client, key string) (T, error) {
	var data T
	str, err := rc.Get(ctx, key).Result()
	if err != nil {
		return data, err
	}
	
	if err := json.Unmarshal([]byte(str), &data); err != nil {
		return data, err
	}

	return data, nil
}

func SetToRedis[T any](ctx context.Context,rc *redis.Client, key string, data T, expiration ...time.Duration) error {
	var TTL time.Duration
	if len(expiration) != 0 {
		TTL = expiration[0]
	}else {
		TTL = 0
	}

	str, err  := json.Marshal(&data)
	if err != nil {
		return err
	}

	if err := rc.Set(ctx, key, str, TTL).Err(); err != nil {
		return err
	}

	return nil
}

func DelFromRedis(ctx context.Context, rc *redis.Client, key string) (string) {
	res, err := rc.Del(ctx, key).Result()
	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf("%d key berhasil di hapus", res)
}