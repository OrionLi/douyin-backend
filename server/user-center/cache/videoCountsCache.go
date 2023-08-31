package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
	"user-center/pkg/util"
)

func RedisSetKey(ctx context.Context, key string, value interface{}) error {
	var err error
	defer func() {
		if err != nil {
			util.LogrusObj.Error("<videoCache>", err)
		}
	}()
	if _, err = NewRedisClient(ctx).Get(ctx, key).Result(); err != redis.Nil {
		fmt.Printf("Key is existed %s\n", key)
	}

	err = NewRedisClient(ctx).Set(ctx, key, value, 3*time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}

func RedisGetKey(ctx context.Context, key string) (string, error) {
	var err error
	defer func() {
		if err != nil {
			util.LogrusObj.Error("<videoCache>", err)
		}
	}()
	value, err := NewRedisClient(ctx).Get(ctx, key).Result()
	// 当该key不存在
	if err == redis.Nil {
		return "", nil
	}

	if err != nil {
		return "", err
	}

	return value, nil
}
