package server

import (
	"fmt"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"user-center/pkg/util"
)

func Grpc(add string) {
	listen, _ := net.Listen("tcp", add)
	//创建grpc服务
	grpcServer := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	pb.RegisterUserServiceServer(grpcServer, &UserRPCServer{})
	pb.RegisterRelationServiceServer(grpcServer, &RelationService{})
	fmt.Println("startPPC: ", listen)
	//启动服务
	err := grpcServer.Serve(listen)
	if err != nil {
		fmt.Printf("failed to serve%v", err)
		util.LogrusObj.Error("Service startup error ", err)
		return
	}
}
