package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"video-center/Web/controller"
	"video-center/Web/middleware"
	"video-center/Web/rpc"
	"video-center/cache"
	"video-center/conf"
	"video-center/dao"
	"video-center/service"
)

func main() {
	conf.InitConfig()
	dao.Init()
	cache.Init()
	rpc.Init()
	host := conf.Viper.GetString("http.host")
	r := gin.New()
	douyin := r.Group("/douyin")
	douyin.Use(middleware.LogMiddleware())
	fmt.Println(time.Now().Unix())
	// 视频相关请求
	publish := douyin.Group("/publish")
	publish.POST("/action/", middleware.ValidationPublishActionRequest(), controller.PublishAction)
	publish.GET("/list/", controller.PublishList)
	douyin.GET("feed", controller.Feed)
	// 评论相关请求
	comment := douyin.Group("/comment")
	comment.POST("/action", controller.CommentAction)
	comment.GET("/list", controller.CommentList)
	// 点赞相关请求
	favoriteController := controller.NewFavoriteController(service.NewFavoriteService(context.Background()))
	favorite := douyin.Group("/favorite")
	favorite.POST("/action", favoriteController.ActionFav)
	favorite.GET("/list", favoriteController.ListFav)
	//开启端口监听
	if err := http.ListenAndServe(host, r); err != nil {
		fmt.Println(err)
	}
}
