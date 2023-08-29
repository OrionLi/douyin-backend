package dao

import (
	"chat-center/conf"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"io"
	"log"
)

var ESClient *elasticsearch.Client

// Init 初始化 Elasticsearch 连接
func Init() {
	// 构建 Elasticsearch 连接配置
	cfg := elasticsearch.Config{
		Addresses: []string{fmt.Sprintf("http://%s:%d", conf.ESAddress, conf.ESport)},
		Username:  conf.ESUser,
		Password:  conf.ESPassword,
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
	if res.StatusCode != 200 {
		log.Fatalf("ElsaticSearch connect fail: %s", res.Status())
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalf("Error closing the response body: %s", err)
		}
	}(res.Body)

	// 将连接赋值给全局变量
	ESClient = client

	// 打印连接信息
	log.Printf("[%s] Elasticsearch connection OK", res.Status())
}
