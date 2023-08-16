package main

import (
	videoCenter "douyin-backend/vedio-center/kitex_gen/videoCenter/videocenter"
	"log"
)

func main() {
	svr := videoCenter.NewServer(new(VideoCenterImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
