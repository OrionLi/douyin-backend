package main

import (
	"chat-center/conf"
	"chat-center/controller"
	"chat-center/dao"
	"chat-center/service"
	"context"
	"fmt"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"github.com/gin-gonic/gin"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
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

	//go testGRPC()

	select {}
}

func initGin() {
	// Gin服务
	chatService := service.NewChatService()
	diaryHandler := controller.NewChatHandler(chatService)

	r := gin.Default()
	api := r.Group("/douyin/message")
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

// FIXME Nacos实例健康数为0，但是RPC服务可以正常访问
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

	// 至少一个ServerConfig
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
		Ip:          conf.GRPCAddress, // 设置你的服务器的 IP 地址
		Port:        uint64(conf.GRPCPort),
		Metadata:    map[string]string{"protocol": "grpc"}, // 设置元数据
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

// HACK 测试方法，随后删除
func testGRPC() {
	// 连接 gRPC 服务
	conn, err := grpc.Dial("localhost:9422", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// 创建 gRPC 客户端
	client := pb.NewDouyinMessageServiceClient(conn)

	// 调用 gRPC 方法
	getMessageRequest := &pb.DouyinMessageChatRequest{
		SelfUserId: 1,
		ToUserId:   123,
		PreMsgTime: 0, // 设置合适的时间戳
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	getMessageResponse, err := client.GetMessage(ctx, getMessageRequest)
	if err != nil {
		log.Fatalf("Failed to call GetMessage: %v", err)
	}

	// 打印响应
	fmt.Printf("GetMessage Response: %+v\n", getMessageResponse)

	fmt.Println("=====================================")

	// 调用 gRPC 方法
	sendMessageRequest := &pb.DouyinMessageActionRequest{
		SelfUserId: 1,
		ToUserId:   123,
		Content:    "Hello World!",
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	sendMessageResponse, err := client.SendMessage(ctx, sendMessageRequest)
	if err != nil {
		log.Fatalf("Failed to call SendMessage: %v", err)
	}

	// 打印响应
	fmt.Printf("SendMessage Response: %+v\n", sendMessageResponse)
}
