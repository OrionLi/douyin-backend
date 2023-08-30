package conf

import (
	"github.com/spf13/viper"
)

var (
	ServiceName string
	ServerIp    string
	ServerPort  uint64
	NacosIp     string

	DbHost     string
	DbPort     string
	DbUser     string
	DbPassword string
	DbName     string

	RedisDb     string
	RedisAddr   string
	RedisPw     string
	RedisDbName string
)

// nacos配置项
var (
	NacosAddress string
	NacosPort    uint64
)

// 微服务的服务名
var (
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
	DbHost = viper.GetString("mysql.DbHost")
	DbPort = viper.GetString("mysql.DbPort")
	DbUser = viper.GetString("mysql.DbUser")
	DbPassword = viper.GetString("mysql.DbPassword")
	DbName = viper.GetString("mysql.DbName")
	RedisDb = viper.GetString("redis.RedisDb")
	RedisAddr = viper.GetString("redis.RedisAddr")
	RedisPw = viper.GetString("redis.RedisPw")
	RedisDbName = viper.GetString("redis.RedisDbName")

	// 解析 Nacos 配置
	NacosAddress = viper.GetString("nacos.Ip")
	NacosPort = uint64(viper.GetInt("nacos.Port"))

	// 解析各个微服务的服务名

	VideoCenterServiceName = viper.GetString("application.video-center.ServiceName")
	ServiceName = viper.GetString("application.ServiceName")
	ServerIp = viper.GetString("application.Ip")
	ServerPort = viper.GetUint64("application.Port")
	NacosIp = viper.GetString("nacos.Ip")
	NacosPort = viper.GetUint64("nacos.Port")

	return nil
}
