package main

import "douyin-backend/video-center/dao"

func main() {
	//svr := videocenter.NewServer(new(VideoCenterImpl))
	//
	//err := svr.Run()
	//
	//if err != nil {
	//	log.Println(err.Error())
	//}
	dao.Init()
	dao.Migrate()
}
