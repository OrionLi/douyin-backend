package main

import (
	conf2 "douyin-backend/server/user-center/conf"
	"douyin-backend/server/video-center/cache"
	"douyin-backend/server/video-center/conf"
	"douyin-backend/server/video-center/dao"
	"douyin-backend/server/video-center/oss"
)

func main() {
	conf.InitConfig()
	cache.Init()
	dao.Init()
	oss.Init("D://d", "OssConf.yaml")
	conf2.Init()
}
