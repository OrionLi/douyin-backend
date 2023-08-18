package main

import (
	"douyin-backend/video-center/generated/video"
	"douyin-backend/video-center/handler"
	"google.golang.org/grpc"
)

func main() {
	server := grpc.NewServer()
	video.RegisterVideoCenterServer(server, &handler.VideoServer{})
}
