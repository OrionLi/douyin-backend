package server

import (
	"fmt"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"user-center/pkg/util"
)

// Grpc 启动grpc服务的
func Grpc(addr string) error {

	listen, _ := net.Listen("tcp", addr)
	// 创建一个grpc服务，并设置不安全的证书 todo: 后期改用STL
	grpcServer := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

	pb.RegisterUserServiceServer(grpcServer, &UserRPCServer{})       //注册用户服务
	pb.RegisterRelationServiceServer(grpcServer, &RelationService{}) //注册关注服务

	//启动服务
	err := grpcServer.Serve(listen)
	if err != nil {
		fmt.Printf("failed to serve%v", err)
		util.LogrusObj.Error("Service startup error ", err)
		return err
	}
	return nil
}
