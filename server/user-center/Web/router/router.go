package router

import (
	"github.com/gin-gonic/gin"
	"web/api"
	"web/middleware"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	g := r.Group("/douyin")
	{
		g.POST("user/register/", api.UserRegister)
		g.POST("user/login/", api.UserLogin)
		authed := g.Group("/") //需要登录保护
		authed.Use(middleware.JWT())
		{
			authed.GET("user/", api.GetUserById)
		}

	}
	return r
}
