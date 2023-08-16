package service

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"log"
	"user-center/dao"
	"user-center/model"
)

type CreateUserService struct {
	UserName string `json:"user_name" form:"user_name"`
	Password string `json:"password" form:"password"`
}

func (service *CreateUserService) Register(ctx context.Context) {
	var user model.User

	userDao := dao.NewUserDao(ctx)

	_, exist, err := userDao.ExistOrNotByUserName(service.UserName)
	if err != nil {
		log.Println(err)
		return //todo: 需添加返回
	}
	//若该用户已存在
	if exist {
		log.Println(err)
		return //todo: 需添加返回
	}
	// 密码加密
	password, err := bcrypt.GenerateFromPassword([]byte(service.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)

		return //todo: 需添加返回
	}

	user = model.User{
		Username:    service.UserName,
		Password:    string(password),
		FollowCount: 0,
		FanCount:    0,
		Follows:     nil,
		Fans:        nil,
	}
	err = userDao.CreateUser(&user)
	if err != nil {
		log.Println(err)

		return //todo: 需添加返回
	}

}
