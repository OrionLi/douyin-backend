package service

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"user-center/dao"
	"user-center/pb"
	"user-center/pkg/e"
	"user-center/pkg/util"
)

type LoginUserService struct {
	UserName string `json:"user_name" form:"user_name"`
	Password string `json:"password" form:"password"`
}

func (service *LoginUserService) Login(ctx context.Context) (*pb.DouyinUserLoginResponse, error) { //todo: 添加返回结构体
	userDao := dao.NewUserDao(ctx)
	user, exist, err := userDao.ExistOrNotByUserName(service.UserName)
	if err != nil {
		return nil, e.NewError(e.Error)
	}
	if exist == false {
		fmt.Println("该用户不存在")
		return nil, e.NewError(e.ErrorExistUserNotFound)
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(service.Password)); err != nil {
		fmt.Println("密码错误")
		return nil, e.NewError(e.Error)
	}
	//签发token
	token, err := util.GenerateToken(user.ID, service.UserName, 0)
	if err != nil {
		fmt.Println("token签发失败")
		return nil, e.NewError(e.Error)
	}
	return &pb.DouyinUserLoginResponse{
		UserId: int64(user.ID),
		Token:  token,
	}, nil
}
