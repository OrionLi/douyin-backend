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

// Redis 引擎初始化
func Redis(redisDb, redisAddr, redisPw, redisDbName string) error {
	RedisDb = redisDb
	RedisAddr = redisAddr
	RedisPw = redisPw
	RedisDbName = redisDbName
	db, _ := strconv.ParseUint(RedisDbName, 10, 64)

	client := redis.NewClient(&redis.Options{
		Addr:     RedisAddr,
		Password: RedisPw,
		DB:       int(db),
	})
	//创建上下文
	ctx := context.Background()
	//测试连接
	err := client.Ping(ctx).Err()
	if err != nil {
		return err
	}
	_redis = client
	return nil
}

func NewRedisClient(ctx context.Context) *redis.Client {
	client := _redis
	return client.WithContext(ctx)
}
