package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"video-center/Web/controller"
	"video-center/Web/rpc"
)

func main() {
	rpc.Init()
	r := gin.New()
	douyin := r.Group("/douyin")
	fmt.Println(time.Now().Unix())
	publish := douyin.Group("/publish")
	publish.POST("/action/", controller.PublishAction)
	publish.GET("/list/", controller.PublishList)
	douyin.GET("feed", controller.Feed)
	if err := http.ListenAndServe(":9999", r); err != nil {
		fmt.Println(err)
	}

}
