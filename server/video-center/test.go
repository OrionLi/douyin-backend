package main

import (
	"context"
	"fmt"
	"video-center/cache"
	"video-center/conf"
	"video-center/dao"
	"video-center/handler"
	"video-center/oss"
	"video-center/pkg/pb"
)

func main() {
	conf.InitConfig()
	cache.Init()
	dao.Init()
	oss.Init("D://d", "OssConf.yaml")
	videoServer := handler.VideoServer{}
	list, err := videoServer.PublishList(context.Background(),
		&pb.DouyinPublishListRequest{Token: "",
			UserId: 6})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(list.VideoList[0].Title)
}
