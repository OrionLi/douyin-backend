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
		&pb.DouyinPublishListRequest{Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6NiwidXNlcl9uYW1lIjoib3JpIiwiYXV0aG9yaXR5IjowLCJleHAiOjE2OTI0MDk3NTIsImlzcyI6IkZhbk9uZS1naW4tbWFsbCJ9.ihgXqU_IdnzAkUIYwg6GzVwRWmtQBDmdVXhwqHdiaJY",
			UserId: 6})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(list.VideoList[0].Title)
}
