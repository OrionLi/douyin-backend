package main

import (
	"gateway-center/controller"
	"gateway-center/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.Default()
	group := router.Group("/douyin")
	{
		// JWT认证中间件
		authed := group.Group("/") //需要token认证保护
		authed.Use(middleware.JWT())
		{
			//可通过ctx.get("id")来获取user_id
			authed.GET("user/", controller.GetUser)
		}

		// chat模块路由
		group.POST("/message/action/", controller.SendMessage)
		group.GET("/message/chat/", controller.GetMessage)

		// user模块路由
		group.POST("/user/register/", controller.UserRegister)
		group.POST("/user/login/", controller.UserLogin)
		// TODO 用户信息、用户关系相关请求

		// video模块路由
		// TODO 视频流相关请求
		group.GET("/video/favorite/action/", controller.ActionFav)
		group.GET("/video/favorite/list/", controller.ListFav)
		group.GET("/video/comment/action/", controller.CommentAction)
		group.GET("/video/comment/list/", controller.CommentList)

	}

	return router
}
