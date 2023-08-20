package cache

import (
	"chat-center/conf"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
)

var RedisClient *redis.Client

func Init() {
	// 初始化 Redis 客户端
	redisClient := redis.NewClient(&redis.Options{
		Addr:     conf.RedisAddress + ":" + strconv.Itoa(conf.RedisPort),
		Password: conf.RedisPassword,
		DB:       conf.RedisDB,
	})

	// 测试连接
	_, err := redisClient.Ping(redisClient.Context()).Result()
	if err != nil {
		log.Printf("Failed to connect to Redis: %v", err)
	}

	RedisClient = redisClient

	log.Printf("Redis connection OK")
}
