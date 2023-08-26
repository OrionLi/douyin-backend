package main

import (
	"gateway/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	g := r.Group("/douyin")
	{
		g.POST("user/register/")
		g.POST("user/login/")
		authed := g.Group("/") //需要token认证保护
		authed.Use(middleware.JWT())
		{
			/*
				可通过ctx.get("id")来获取user_id
			*/
			authed.GET("user/")
		}

	}
	return r
}
