package main

import (
	"log"
	"strings"
	"user-center/cache"
	"user-center/conf"
	"user-center/dao"
	"user-center/pkg/util"
	"user-center/server"
)

func main() {
	//初始化配置文件
	err := conf.Init()
	if err != nil {
		util.LogrusObj.Error("<Main> : ", err)
	}
	//nacos注册
	server.RegisterNacos(conf.ServerIp, conf.ServiceName, conf.NacosIp, conf.NacosPort, conf.ServerPort)
	//mysql连接信息
	conn := strings.Join([]string{conf.DbUser, ":", conf.DbPassword, "@tcp(", conf.DbHost, ":", conf.DbPort, ")/", conf.DbName, "?charset=utf8mb4&parseTime=true"}, "")
	// gorm引擎初始化
	err = dao.Database(conn)
	if err != nil {
		log.Fatal(err)
	}
	// redis引擎初始化
	err = cache.Redis(conf.RedisDb, conf.RedisAddr, conf.RedisPw, conf.RedisDbName)
	if err != nil {
		log.Fatal(err)
	}
	// grpc初始化
	err = server.Grpc(conf.ServerIp)
	if err != nil {
		log.Fatal(err)
	}
}
