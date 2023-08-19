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
	// 测试
	/*	isFollow := service.IsFollowService{
			UserId:       3,
			FollowUserId: 5,
		}
		b, err := isFollow.IsFollow(ctx)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(b)*/
	id := service.GetUserByIdService{
		Id: 2,
	}
	id.GetUserById(ctx)
}
