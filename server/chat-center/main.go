package main

import (
	"chat-center/Web/controller"
	"chat-center/Web/middleware"
	"chat-center/Web/rpc/client"
	"chat-center/conf"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	// 初始化配置
	conf.InitConf()
	// 初始化grpc
	client.InitGRPCClient()
	defer client.ClosConn()
	// 初始化gin
	r := gin.Default()
	api := r.Group("/douyin/message")
	api.Use(middleware.LogMiddleware())
	{
		api.GET("/chat", controller.GetMessage)
		api.POST("/action", controller.SendMessage)
	}

	if err := r.Run(":" + conf.WebPort); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
