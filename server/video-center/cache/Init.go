package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"video-center/conf"
)

var RedisClient *redis.Client

func Init() {
	host := conf.Viper.GetString("db.redis.host")
	db := conf.Viper.GetInt("db.redis.db")
	passwd := conf.Viper.GetString("db.redis.passwd")
	fmt.Println(host)
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
