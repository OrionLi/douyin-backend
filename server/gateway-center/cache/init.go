package cache

import (
	"context"
	"fmt"
	"gateway-center/conf"
	"github.com/go-redis/redis/v8"
	"time"
)

// todo 后续可加入Feed流缓存、PublishList缓存等等

var RedisClient *redis.Client

func Init() {
	RedisClient = redis.NewClient(&redis.Options{
		Password: conf.RedisPasswd,
		Addr:     conf.RedisHost,
		DB:       conf.RedisDB,
	})

	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println(err)
		panic("redis ping error")
	}
}
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
