package cache

import (
	"context"
	"github.com/go-redis/redis/v8"

	"strconv"
)

var (
	_redis      *redis.Client
	RedisDb     string
	RedisAddr   string
	RedisPw     string
	RedisDbName string
)

func Redis(redisDb, redisAddr, redisPw, redisDbName string) {
	RedisDb = redisDb
	RedisAddr = redisAddr
	RedisPw = redisPw
	RedisDbName = redisDbName
	db, _ := strconv.ParseUint(RedisDbName, 10, 64)
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:     RedisAddr,
		Password: RedisPw,
		DB:       int(db),
	})
	err := client.Ping(ctx).Err()
	if err != nil {
		panic(err)
	}
	_redis = client

}

func NewRedisClient(ctx context.Context) *redis.Client {
	client := _redis
	return client.WithContext(ctx)
}
