package common

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

var (
	Conn *grpc.ClientConn
)

// grpc连接
func Grpc_conn() *grpc.ClientConn {
	addr := "127.0.0.1"
	port := "3001"
	addr = addr + ":" + port
	fmt.Println(addr)

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: #{err}")
	}
	Conn = conn
	return conn
}

func GetConn() *grpc.ClientConn {
	return Conn
}
