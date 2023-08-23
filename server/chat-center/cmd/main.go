package main

import (
	"chat-center/Web/rpc/server"
)

func main() {
	// 初始化gRPC
	go server.InitGRPCServer()
	// 注册nacos
	go server.RegisterNacos()
	// 堵塞程序
	select {}
}
