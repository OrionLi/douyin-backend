package conf

import (
	"github.com/spf13/viper"
	"os"
	"strings"
	"user-center/cache"
	"user-center/dao"
	"user-center/server"
)

var (
	Address string

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

// nacos配置项
var (
	NacosAddress string
	NacosPort    int
)

// 各个微服务的服务名
var (
	ChatCenterServiceName  string
	UserCenterServiceName  string
	VideoCenterServiceName string
)

// Init 初始化配置文件与引擎
func Init() error {

	// 设置配置文件的名称和路径
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./conf/")

	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	// 获取配置项的值并赋值给变量
	Address = viper.GetString("server.address")

	DbHost = viper.GetString("mysql.DbHost")
	DbPort = viper.GetString("mysql.DbPort")
	DbUser = viper.GetString("mysql.DbUser")
	DbPassword, _ = os.LookupEnv("MYSQL_PASSWORD")
	DbName = viper.GetString("mysql.DbName")
	RedisDb = viper.GetString("redis.RedisDb")
	redisAddr = viper.GetString("redis.RedisAddr")
	RedisPw, _ = os.LookupEnv("REDIS_PASSWORD")
	RedisDbName = viper.GetString("redis.RedisDbName")

	// 解析 Nacos 配置
	NacosAddress = viper.GetString("nacos.Ip")
	NacosPort = viper.GetInt("nacos.Port")

	// 解析各个微服务的服务名
	ChatCenterServiceName = viper.GetString("application.chat-center.ServiceName")
	UserCenterServiceName = viper.GetString("application.user-center.ServiceName")
	VideoCenterServiceName = viper.GetString("application.video-center.ServiceName")

	//mysql连接信息
	conn := strings.Join([]string{DbUser, ":", DbPassword, "@tcp(", DbHost, ":", DbPort, ")/", DbName, "?charset=utf8mb4&parseTime=true"}, "")
	// gorm引擎初始化
	err = dao.Database(conn)
	if err != nil {
		return err
	}
	// redis引擎初始化
	err = cache.Redis(RedisDb, redisAddr, RedisPw, RedisDbName)
	if err != nil {
		return err
	}
	// grpc初始化
	err = server.Grpc(Address)
	if err != nil {
		return err
	}
	return nil
}
