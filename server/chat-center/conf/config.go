package conf

import (
	"github.com/spf13/viper"
	"log"
)

var (
	ESAddress  string
	ESport     int
	ESUser     string
	ESPassword string
)

var (
	RedisAddress  string
	RedisPort     int
	RedisPassword string
	RedisDB       int
)

var (
	WebPort string
)

func InitConf() {
	// 初始化 Viper 配置
	viper.SetConfigName("config")
	viper.AddConfigPath("./conf")
	err := viper.ReadInConfig() // 读取配置文件
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	// 解析 Elasticsearch 配置
	elasticsearchConfig := viper.Sub("elasticsearch")
	if elasticsearchConfig == nil {
		log.Fatal("Missing 'elasticsearch' configuration section")
	}
	ESAddress = elasticsearchConfig.GetString("address")
	ESport = elasticsearchConfig.GetInt("port")
	ESUser = elasticsearchConfig.GetString("user")
	ESPassword = elasticsearchConfig.GetString("password")

	redisConfig := viper.Sub("redis")
	if redisConfig == nil {
		log.Fatal("Missing 'redis' configuration section")
	}
	RedisAddress = redisConfig.GetString("address")
	RedisPort = redisConfig.GetInt("port")
	RedisPassword = redisConfig.GetString("password")
	RedisDB = redisConfig.GetInt("db")

	ginConfig := viper.Sub("gin")
	if ginConfig == nil {
		log.Fatal("Missing 'gin' configuration section")
	}
	WebPort = ginConfig.GetString("port")
}
