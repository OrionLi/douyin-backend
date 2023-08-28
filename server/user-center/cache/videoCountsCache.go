package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

func RedisSetKey(ctx context.Context, key string, value interface{}) error {
	if _, err := NewRedisClient(ctx).Get(ctx, key).Result(); err != redis.Nil {
		fmt.Printf("Key is existed %s\n", key)
	}
	err := NewRedisClient(ctx).Set(ctx, key, value, 3*time.Minute).Err()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Redis client set successfully...%s\n", key)
	return nil
}

func RedisGetKey(ctx context.Context, key string) (string, error) {
	value, err := NewRedisClient(ctx).Get(ctx, key).Result()

	if err == redis.Nil {
		return value, err
	}

	if err != nil {
		panic(err)
	}

	fmt.Printf("Redis client get successfully... %s\n", key)
	return value, nil
}
