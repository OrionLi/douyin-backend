package main

import (
	"chat-center/conf"
	"chat-center/dao"
	pb "chat-center/generated/message"
	"chat-center/rpc"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	conf.InitConf()

	// 初始化ES连接
	dao.Init()

	// 创建ChatService
	//chatService := service.NewChatService()

	// 创建grpc服务
	grpcServer := grpc.NewServer()

	// 注册ChatService
	pb.RegisterDouyinMessageServiceServer(grpcServer, rpc.NewChatRPCService())

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
