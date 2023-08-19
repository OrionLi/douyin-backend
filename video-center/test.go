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
	conf.InitConfig()
	cache.Init()
	dao.Init()
	oss.Init("D://d", "OssConf.yaml")
	server := handler.VideoServer{}
	//Feed测试
	unix := time.Now().Unix()
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6NiwidXNlcl9uYW1lIjoib3JpIiwiYXV0aG9yaXR5IjowLCJleHAiOjE2OTI0MDk3NTIsImlzcyI6IkZhbk9uZS1naW4tbWFsbCJ9.ihgXqU_IdnzAkUIYwg6GzVwRWmtQBDmdVXhwqHdiaJY"
	feed, err := server.Feed(context.Background(), &video.DouyinFeedRequest{
		LatestTime: &unix,
		Token:      &token,
	})
	if err != nil {
		return
	}
	fmt.Println("视频1")
	fmt.Println(feed.VideoList[0].Title)
	fmt.Println(feed.VideoList[0].IsFavorite)
	fmt.Println(feed.VideoList[0].Author.IsFollow)
	fmt.Println("视频2")
	fmt.Println(feed.VideoList[1].Title)
	fmt.Println(feed.VideoList[1].IsFavorite)
	fmt.Println(feed.VideoList[1].Author.IsFollow)

	////Action测试
	//open, err := os.Open("D:\\BaiduNetdiskDownload\\oceans.mp4")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//defer func(open *os.File) {
	//	err := open.Close()
	//	if err != nil {
	//
	//	}
	//}(open)
	//stat, err := open.Stat()
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//bs := make([]byte, stat.Size())
	//_, err = bufio.NewReader(open).Read(bs)
	//if err != nil && err != io.EOF {
	//	fmt.Println(err)
	//	return
	//}
	//resp, err := server.PublishAction(context.Background(), &video.DouyinPublishActionRequest{
	//	Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6NiwidXNlcl9uYW1lIjoib3JpIiwiYXV0aG9yaXR5IjowLCJleHAiOjE2OTI0MDk3NTIsImlzcyI6IkZhbk9uZS1naW4tbWFsbCJ9.ihgXqU_IdnzAkUIYwg6GzVwRWmtQBDmdVXhwqHdiaJY",
	//	Title: "测试视频",
	//	Data:  bs,
	//})
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(resp.StatusCode)
	//fmt.Println(resp.StatusMsg)
}
