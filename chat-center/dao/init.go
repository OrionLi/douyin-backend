package dao

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var Rdb *redis.Client

func init() {
	viper.SetConfigFile("../conf/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("error reading config file: %s", err))
	}

	// 读取配置信息
	redisAddr := viper.GetString("redis.address")
	redisPassword := viper.GetString("redis.password")
	redisDB := viper.GetInt("redis.db")

	// 创建一个 Redis 客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})
	Rdb = rdb
}
