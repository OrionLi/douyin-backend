package main

import (
	"video-center/cache"
	"video-center/conf"
	"video-center/dao"
	"video-center/oss"
)

func main() {
	conf.InitConfig()
	cache.Init()
	dao.Init()
	oss.Init("D://d", "OssConf.yaml")
}
