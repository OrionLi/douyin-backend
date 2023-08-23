package main

import (
	"chat-center/conf"
	"chat-center/dao"
	"chat-center/handler"
	"chat-center/service"
	"fmt"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"github.com/gin-gonic/gin"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	// 初始化配置
	conf.InitConf()

	// 初始化数据库
	dao.Init()

	// 初始化gin
	go initGin()

	// 初始化gRPC
	go initGRPC()

	go RegisterNacos()

	select {}
}

func initGin() {
	// Gin服务
	chatService := service.NewChatService()
	diaryHandler := handler.NewChatHandler(chatService)

	r := gin.Default()
	api := r.Group("/douyin/message")
	api.Use(handler.LogMiddleware())
	{
		api.GET("/chat", diaryHandler.GetMessage)
		api.POST("/action", diaryHandler.SendMessage)
	}

	if err := r.Run(":" + conf.WebPort); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

func initGRPC() {
	// 创建grpc服务
	grpcServer := grpc.NewServer()

	// 注册ChatService
	pb.RegisterDouyinMessageServiceServer(grpcServer, service.NewChatRPCService())

	// 监听指定端口
	listener, err := net.Listen("tcp", ":9422")
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
	nacosClient, err := clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})

	// 注册 gRPC 服务到 Nacos
	instance := &vo.RegisterInstanceParam{
		Ip:   conf.GRPCAddress, // 设置你的服务器的 IP 地址
		Port: uint64(conf.GRPCPort),
		Metadata: map[string]string{
			"protocol":            "grpc",
			"healthCheckType":     "TCP", // 使用 TCP 健康检查
			"healthCheckPort":     fmt.Sprintf("%d", uint64(conf.GRPCPort)),
			"healthCheckPath":     "",    // 空路径
			"healthCheckInterval": "10s", // 健康检查间隔
		},
		ClusterName: "default",
		ServiceName: conf.NacosServerName,
		GroupName:   conf.NacosGroup,
	}

	_, err = nacosClient.RegisterInstance(*instance)

	if err != nil {
		log.Fatalf("Error registering instance: %v", err)
	}

	log.Println("gRPC server is registered on Nacos")
}
