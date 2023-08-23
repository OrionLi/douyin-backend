package main

import (
	"chat-center/Web/controller"
	"chat-center/Web/middleware"
	"chat-center/conf"
	"chat-center/dao"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	// 初始化配置
	conf.InitConf()
	// 初始化数据库
	dao.Init()
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
