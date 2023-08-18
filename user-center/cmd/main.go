package main

import (
	"context"
	"fmt"
	"user-center/conf"
	"user-center/service"
)

func main() {
	ctx := context.Background()
	//初始化配置文件
	conf.Init()
	isFollow := service.IsFollowService{
		UserId:       3,
		FollowUserId: 5,
	}
	b, err := isFollow.IsFollow(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(b)
}
