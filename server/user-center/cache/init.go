package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
	"sync"
)

var (
	_redis      *redis.Client
	RedisDb     string
	RedisAddr   string
	RedisPw     string
	RedisDbName string

	redisOnce sync.Once
)

func Redis(redisDb, redisAddr, redisPw, redisDbName string) error {
	RedisDb = redisDb
	RedisAddr = redisAddr
	RedisPw = redisPw
	RedisDbName = redisDbName
	db, _ := strconv.ParseUint(RedisDbName, 10, 64)

	redisOnce.Do(func() {
		client := redis.NewClient(&redis.Options{
			Addr:     RedisAddr,
			Password: RedisPw,
			DB:       int(db),
			PoolSize: 10, // 配置连接池容量
		})

		// 创建上下文
		ctx := context.Background()

		// 测试连接
		if err := client.Ping(ctx).Err(); err != nil {
			panic("Failed to connect to Redis: " + err.Error())
		}

		_redis = client
	})

	return nil
}

func NewRedisClient(ctx context.Context) *redis.Client {
	client := _redis
	return client.WithContext(ctx)
}
