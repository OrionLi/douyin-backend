package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"video-center/Web/controller"
	"video-center/Web/rpc"
)

// todo 完成gin框架部分，封装user时候调用其他模块获取user以及是否关注，并且封装在其中
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
