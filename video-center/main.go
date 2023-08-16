package main

import (
	"context"
	"douyin-backend/video-center/dao"
	"douyin-backend/video-center/service"
	"fmt"
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
	dao.Init()
	list, err := service.NewVideoService(context.Background()).PublishList(1)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(list[0].PlayUrl)
	fmt.Println(list[0].CoverUrl)
}
