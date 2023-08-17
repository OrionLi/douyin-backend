package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

func RedisSetKey(ctx context.Context, key string, value interface{}) error {
	if _, err := RedisClient.Get(ctx, key).Result(); err != redis.Nil {
		fmt.Println("Redis client set successfully...")
	}
	err := RedisClient.Set(ctx, key, value, 3*time.Minute).Err()
	if err != nil {
		panic(err)
	}
	if _, ok := value.(string); ok {
		err = RedisClient.Set(ctx, value.(string), key, 3*time.Minute).Err()
	}
	if err != nil {
		panic(err)
	}

	fmt.Println("Redis client set successfully...")
	return nil
}

func RedisGetKey(ctx context.Context, key string) (string, error) {
	value, err := RedisClient.Get(ctx, key).Result()

	if err == redis.Nil {
		return value, err
	}

	if err != nil {
		panic(err)
	}

	fmt.Println("Redis client get successfully...")
	return value, nil
}
