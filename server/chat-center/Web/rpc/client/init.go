package client

import (
	"chat-center/conf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"strconv"
)

var (
	Conn *grpc.ClientConn
)

func InitGRPCClient() {
	// 建立到 gRPC 服务器的连接。
	target := ":" + strconv.Itoa(conf.GRPCPort)
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("connect failed: %v", err)
	}
	defer conn.Close()

	Conn = conn
}

func GetMsgConn() *grpc.ClientConn {
	return Conn
}
