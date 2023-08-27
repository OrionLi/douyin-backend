package server

import (
	"chat-center/conf"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
)

func InitGRPCServer() {
	// 创建grpc服务
	grpcServer := grpc.NewServer()

	// 注册ChatService
	pb.RegisterDouyinMessageServiceServer(grpcServer, NewChatServer())

	// 监听指定端口
	address := ":" + strconv.Itoa(conf.GRPCPort)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("gRPC server is listening on port : 9422")

	// 启动gRPC服务器
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func RegisterNacos() {
	// 创建clientConfig
	clientConfig := constant.ClientConfig{
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogLevel:            "debug",
	}

	// 创建ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: conf.NacosAddress,
			Port:   uint64(conf.NacosPort),
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

	// 注册 gRPC 服务到 Nacos
	instance := &vo.RegisterInstanceParam{
		Ip:          conf.GRPCAddress, // 设置你的服务器的 IP 地址
		Port:        uint64(conf.GRPCPort),
		ServiceName: conf.NacosServerName, // 设置服务的名称
		GroupName:   conf.NacosGroup,      // 设置服务的分组
		Weight:      10,                   // 权重为10
		Enable:      true,                 // 设置实例为可用状态
		Healthy:     true,                 // 设置实例为健康状态
	}

	_, err = nacosClient.RegisterInstance(*instance)

	if err != nil {
		log.Fatalf("Error registering instance: %v", err)
	}

	log.Println("gRPC server is registered on Nacos")
}
