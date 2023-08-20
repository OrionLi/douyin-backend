package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"video-center/Web/controller"
	"video-center/Web/rpc"
	"video-center/cache"
	"video-center/conf"
	"video-center/dao"
)

func main() {
	conf.InitConfig()
	cache.Init()
	rpc.Init()
	dao.Init()
	r := gin.New()
	douyin := r.Group("/douyin")
	fmt.Println(time.Now().Unix())
	//视频相关请求
	publish := douyin.Group("/publish")
	publish.POST("/action/", controller.PublishAction)
	publish.GET("/list/", controller.PublishList)
	douyin.GET("feed", controller.Feed)
	//评论相关请求
	comment := douyin.Group("/comment")
	comment.POST("/action", controller.CommentAction)
	comment.GET("/list", controller.CommentList)
	//开启端口监听
	if err := http.ListenAndServe(":9999", r); err != nil {
		fmt.Println(err)
	}
}
