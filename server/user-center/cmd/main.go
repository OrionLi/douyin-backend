package main

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"user-center/conf"
	"user-center/pb"
	"user-center/server"
)

func main() {

	//初始化配置文件
	conf.Init()
	// 创建grpc服务
	listen, _ := net.Listen("tcp", ":8088")
	//创建grpc服务
	grpcServer := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	pb.RegisterUserServiceServer(grpcServer, &server.UserRPCServer{})

	//启动服务

	err := grpcServer.Serve(listen)
	if err != nil {
		fmt.Printf("failed to serve%v", err)
		return
	}
}
