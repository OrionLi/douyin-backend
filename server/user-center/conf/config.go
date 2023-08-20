package conf

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
	"user-center/cache"
	"user-center/dao"
	"user-center/server"
)

var (
	Address string

	Db         string
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassword string
	DbName     string

	RedisDb     string
	redisAddr   string
	RedisPw     string
	RedisDbName string
)

func Init() {

	// 设置配置文件的名称和路径
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	viper.AddConfigPath("./conf/")

	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	// 获取配置项的值并赋值给变量
	Address = viper.GetString("server.address")

	Db = viper.GetString("mysql.DB")
	DbHost = viper.GetString("mysql.DbHost")
	DbPort = viper.GetString("mysql.DbPort")
	DbUser = viper.GetString("mysql.DbUser")
	DbPassword = viper.GetString("mysql.DbPassword")
	DbName = viper.GetString("mysql.DbName")

	RedisDb = viper.GetString("redis.RedisDb")
	redisAddr = viper.GetString("redis.RedisAddr")

	RedisPw = viper.GetString("redis.RedisPw")
	RedisDbName = viper.GetString("redis.RedisDbName")
	fmt.Println(Db)
	//mysql

	conn := strings.Join([]string{DbUser, ":", DbPassword, "@tcp(", DbHost, ":", DbPort, ")/", DbName, "?charset=utf8mb4&parseTime=true"}, "")

	dao.Database(conn)
	cache.Redis(RedisDb, redisAddr, RedisPw, RedisDbName)
	server.Grpc(Address)

}
