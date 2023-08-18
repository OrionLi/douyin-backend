package main

import (
	"context"
	"douyin-backend/video-center/cache"
	"douyin-backend/video-center/conf"
	"douyin-backend/video-center/dao"
	"douyin-backend/video-center/generated/video"
	"douyin-backend/video-center/handler"
	"douyin-backend/video-center/oss"
	"fmt"
	"time"
)

func main() {
	//open, _ := os.Open("D:\\chatfile\\psc.png")
	//defer func(open *os.File) {
	//	err := open.Close()
	//	if err != nil {
	//
	//	}
	//}(open)
	//stat, err := open.Stat()
	//if err != nil {
	//	fmt.Println(err)
	//}
	//size := stat.Size()
	//bytes := make([]byte, size)
	//_, err = open.Read(bytes)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//oss.Init("D://d", "OssConf.yaml")
	//video, s, err := oss.UploadVideo(context.Background(), 1, bytes, "02.jpg")
	//if err != nil {
	//	return
	//}
	//fmt.Println(video)
	//fmt.Println(s)
	//dao.Init()
	//list, err := service.NewVideoService(context.Background()).PublishList(1)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(list[0].PlayUrl)
	//fmt.Println(list[0].CoverUrl)
	//user, err := dao.QueryUserByID(context.Background(), 3)
	//of, err := dao.IsFanOf(1, 3)
	//user.IsFollow = of
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(user.Id)
	//fmt.Println(user.Username)
	//fmt.Println(user.FanCount)
	//fmt.Println(user.FollowCount)
	//fmt.Println(user.IsFollow)
	conf.InitConfig()
	dao.Init()
	oss.Init("D://d", "OssConf.yaml")
	cache.Init()
	unix := time.Now().Unix()
	str := ""
	list, err := handler.FeedVideoList(context.Background(), &video.DouyinFeedRequest{LatestTime: &unix, Token: &str})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(list.VideoList[0].CoverUrl)
	fmt.Println(list.VideoList[0].PlayUrl)
}
