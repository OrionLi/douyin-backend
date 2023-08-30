package conf

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"log"
	"testing"
)

func TestUnReg(t *testing.T) {
	clientConfig := constant.ClientConfig{
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogLevel:            "debug",
	}

	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: "127.0.0.1",
			Port:   8848,
		},
	}

	nacosClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		})

	if err != nil {
		log.Fatalf("failed to create Nacos client: %v", err)
	}

	instance := &vo.DeregisterInstanceParam{
		Ip:          "192.168.1.12",
		Port:        9421,
		ServiceName: "chat-center",
	}

	_, err = nacosClient.DeregisterInstance(*instance)

	if err != nil {
		t.Log("9421不存在或已注销")
	}

	instance = &vo.DeregisterInstanceParam{
		Ip:          "192.168.1.12",
		Port:        9422,
		ServiceName: "chat-center",
	}

	_, err = nacosClient.DeregisterInstance(*instance)

	if err != nil {
		t.Log("9422不存在或已注销")
	}

	t.Log("持久化实例注销成功")
}
