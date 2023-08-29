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
	GRPCAddress string
	GRPCPort    int
)

var (
	NacosAddress    string
	NacosPort       int
	NacosServerName string
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

	// 解析 gRPC 配置
	gRPCConfig := viper.Sub("gRPC")
	if gRPCConfig == nil {
		log.Fatal("Missing 'gRPC' configuration section")
	}
	GRPCAddress = gRPCConfig.GetString("address")
	GRPCPort = gRPCConfig.GetInt("port")

	// 解析 Nacos 配置
	nacosConfig := viper.Sub("nacos")
	if nacosConfig == nil {
		log.Fatal("Missing 'nacos' configuration section")
	}
	NacosAddress = nacosConfig.GetString("address")
	NacosPort = nacosConfig.GetInt("port")
	NacosServerName = nacosConfig.GetString("server-name")
}
