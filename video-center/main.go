package main

import (
	"douyin-backend/video-center/kitex_gen/videoCenter/videocenter"
	"log"
)

func main() {
	svr := videocenter.NewServer(new(VideoCenterImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
