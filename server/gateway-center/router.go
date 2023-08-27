package main

import (
	"gateway-center/controller"
	"gateway-center/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	g := r.Group("/douyin")
	{
		g.POST("user/register/", controller.UserRegister)
		g.POST("user/login/", controller.UserLogin)
		authed := g.Group("/") //需要token认证保护
		authed.Use(middleware.JWT())
		{
			/*
				可通过ctx.get("id")来获取user_id
			*/
			authed.GET("user/", controller.GetUser)
		}

	}
	return r
}
