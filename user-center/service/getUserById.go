package service

import (
	"context"
	"fmt"
	"user-center/dao"
)

type GetUserByIdService struct {
	Id uint `json:"id" form:"id"`
}

func (service *GetUserByIdService) GetUserById(ctx context.Context) { //todo: 添加返回结构体

	//todo: 需添加缓存，并添加逻辑：粉丝数大于等于300为网红

	userDao := dao.NewUserDao(ctx)
	user, err := userDao.GetUserById(service.Id)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println("user:", user)
}
