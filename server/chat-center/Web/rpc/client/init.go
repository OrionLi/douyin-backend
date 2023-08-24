package client

import (
	"chat-center/conf"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

var (
	Conn *grpc.ClientConn
)

func Init() {
	InitChatRPC()
}

func InitChatRPC() {
	// 创建clientConfig
	clientConfig := constant.ClientConfig{
		NamespaceId:         conf.NacosNamespaceId, // 如果需要支持多namespace，我们可以创建多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
	}

	// 创建ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      conf.NacosAddress,
			ContextPath: "/nacos",
			Port:        uint64(conf.NacosPort),
			Scheme:      "http",
		},
	}

	// 创建服务发现客户端
	nacosClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		})

	if err != nil {
		log.Fatalf("failed to create Nacos client: %v", err)
	}

	// 获取 gRPC 服务实例信息
	serviceName := conf.NacosServerName
	instances, err := nacosClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: serviceName,
		GroupName:   conf.NacosGroup,
	})
	if err != nil {
		log.Fatalf("failed to get service instances: %v", err)
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", instances.Ip, instances.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("connect failed: %v", err)
	}
	Conn = conn
}

func GetMsgConn() *grpc.ClientConn {
	return Conn
}

func ClosConn() {
	if Conn != nil {
		err := Conn.Close()
		if err != nil {
			panic(err)
		}
	}
}
