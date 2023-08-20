package main

import (
	"chat-center/conf"
	"chat-center/dao"
	"chat-center/handler"
	"chat-center/service"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	// 初始化配置文件
	conf.InitConf()

	// 初始化ES连接
	dao.Init()

	// Gin服务
	chatService := service.NewChatService()
	diaryHandler := handler.NewDiaryHandler(chatService)

	r := gin.Default()
	api := r.Group("/douyin/message")
	{
		api.GET("/chat", diaryHandler.GetMessage)
		api.POST("/action", diaryHandler.SendMessage)
	}

	if err := r.Run(":" + conf.WebPort); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

	// TODO gRPC服务
	// 创建grpc服务
	grpcServer := grpc.NewServer()

	// 注册ChatService
	pb.RegisterDouyinMessageServiceServer(grpcServer, service.NewChatRPCService())

	// 监听指定端口
	// FIXME: 修复端口号
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// 启动gRPC服务器
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
