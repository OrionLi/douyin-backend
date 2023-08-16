package main

import (
	"context"
	"user-center/conf"
	"user-center/service"
)

func main() {
	ctx := context.Background()
	//初始化配置文件
	conf.Init()
	userRequest := service.LoginUserService{
		UserName: "test",
		Password: "123456",
	}
	userRequest.Login(ctx)

}
