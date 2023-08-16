package cache

import (
	"context"
	"douyin-backend/video-center/conf"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var RedisClient *redis.Client

func Init() {
	conf.InitConfig()
	host := viper.GetString("db.redis.host")
	db := viper.GetInt("db.redis.db")
	passwd := viper.GetString("db.redis.passwd")
	RedisClient = redis.NewClient(&redis.Options{
		Password: passwd,
		Addr:     host,
		DB:       db,
	})

	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println(err)
		panic("redis ping error")
	}
}
