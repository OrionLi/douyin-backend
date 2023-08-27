package grpc

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

var (
	UserConn *grpc.ClientConn
)

// UserClientInit grpc初始化
func UserClientInit() {
	addr := "127.0.0.1"
	port := "3001"
	addr = addr + ":" + port
	fmt.Println(addr)

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: #{err}")
	}
	UserConn = conn

}

func GetUserConn() *grpc.ClientConn {
	return UserConn
}
