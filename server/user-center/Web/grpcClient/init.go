package grpcClient

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

var (
	Conn *grpc.ClientConn
)

// grpc初始化
func init() {
	addr := "127.0.0.1"
	port := "3001"
	addr = addr + ":" + port
	fmt.Println(addr)

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: #{err}")
	}
	Conn = conn

}

func GetConn() *grpc.ClientConn {
	return Conn
}
