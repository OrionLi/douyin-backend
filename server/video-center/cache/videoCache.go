package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
	"video-center/baseResponse"
	"video-center/dao"
)

func RedisSetKey(ctx context.Context, key string, value interface{}) error {
	if _, err := RedisClient.Get(ctx, key).Result(); err != redis.Nil {
		fmt.Printf("Key is existed %s\n", key)
	}
	err := RedisClient.Set(ctx, key, value, 3*time.Minute).Err()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Redis client set successfully...%s\n", key)
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

	fmt.Printf("Redis client get successfully... %s\n", key)
	return value, nil
}

func RedisSetVideoList(ctx context.Context, key string, videoList dao.VideoArray) error {
	if _, err := RedisClient.Get(ctx, key).Result(); err != redis.Nil {
		fmt.Printf("Key is existed %s\n", key)
	}
	err := RedisClient.Set(ctx, key, &videoList, 3*time.Minute).Err()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Redis client set successfully...%s\n", key)
	return nil
}
func RedisGetVideoList(ctx context.Context, key string) (dao.VideoArray, error) {
	videos := dao.VideoArray{}
	err := RedisClient.Get(ctx, key).Scan(&videos)
	if err == redis.Nil {
		return videos, err
	}
	if err != nil {
		panic(err)
	}
	fmt.Printf("Redis client get successfully... %s\n", key)
	return videos, nil
}

func RedisSetHttpVideoList(ctx context.Context, key string, videoList baseResponse.VideoArray) error {
	if _, err := RedisClient.Get(ctx, key).Result(); err != redis.Nil {
		fmt.Printf("Key is existed %s\n", key)
	}
	err := RedisClient.Set(ctx, key, &videoList, 3*time.Minute).Err()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Redis client set successfully...%s\n", key)
	return nil
}
func RedisGetHttpVideoList(ctx context.Context, key string) (baseResponse.VideoArray, error) {
	videos := baseResponse.VideoArray{}
	err := RedisClient.Get(ctx, key).Scan(&videos)
	if err == redis.Nil {
		return videos, err
	}
	if err != nil {
		panic(err)
	}
	fmt.Printf("Redis client get successfully... %s\n", key)
	return videos, nil
}
