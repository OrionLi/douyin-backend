package main

import (
	"chat-center/conf"
	"chat-center/dao"
	"chat-center/server"
)

func main() {
	// 初始化配置
	conf.InitConf()
	// 初始化数据库
	dao.Init()
	// 初始化gRPC
	go server.InitGRPCServer()
	// 注册nacos
	go server.RegisterNacos()
	// 堵塞程序
	select {}
}
