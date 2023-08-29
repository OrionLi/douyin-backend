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
			// 聊天模块
			authed.POST("/message/action/", controller.SendMessage)
			authed.GET("/message/chat/", controller.GetMessage)

			// 用户模块
			authed.GET("/user/", controller.GetUser)

			// 视频模块
			authed.POST("/favorite/action/", controller.ActionFav)
			authed.GET("/favorite/list/", controller.ListFav)
			authed.POST("/comment/action/", controller.CommentAction)
			authed.GET("/comment/list/", controller.CommentList)
		}
		// user模块路由
		group.POST("/user/register/", controller.UserRegister)
		group.POST("/user/login/", controller.UserLogin)
		group.POST("/relation/action/", controller.RelationAction)
		group.GET("/relation/follow/list/", controller.GetFollowList)
		group.GET("/relation/follower/list/", controller.GetFollowerList)
		group.GET("/relation/friend/list/", controller.GetFriendList)

		group.POST("/publish/action/", controller.PublishAction)
		group.GET("/feed/", controller.Feed)
		group.GET("/publish/list/", controller.PublishList)
	}

	return router
}
