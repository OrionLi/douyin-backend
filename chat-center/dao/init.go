package dao

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/spf13/viper"
	"log"
)

var ESClient *elasticsearch.Client

// 初始化 Elasticsearch 连接
func init() {
	// 初始化 Viper 配置
	viper.SetConfigName("config")
	viper.AddConfigPath("../conf")
	err := viper.ReadInConfig() // 读取配置文件
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	// 解析 Elasticsearch 配置
	elasticsearchConfig := viper.Sub("elasticsearch")
	if elasticsearchConfig == nil {
		log.Fatal("Missing 'elasticsearch' configuration section")
	}

	// 获取配置值
	address := elasticsearchConfig.GetString("address")
	port := elasticsearchConfig.GetInt("port")
	user := elasticsearchConfig.GetString("user")
	password := elasticsearchConfig.GetString("password")

	// 构建 Elasticsearch 连接配置
	cfg := elasticsearch.Config{
		Addresses: []string{fmt.Sprintf("http://%s:%d", address, port)},
		Username:  user,
		Password:  password,
	}

	// 创建 Elasticsearch 连接
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// 测试连接
	res, err := client.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	// 将连接赋值给全局变量
	ESClient = client

	// 打印连接信息
	log.Printf("[%s] Elasticsearch connection OK", res.Status())
}
