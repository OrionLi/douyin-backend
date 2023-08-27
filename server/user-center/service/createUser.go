package service

import (
	"context"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"golang.org/x/crypto/bcrypt"
	"user-center/dao"
	"user-center/model"
	e2 "user-center/pkg/e"
	"user-center/pkg/util"
)

// CreateUserService 用户创建服务请求
type CreateUserService struct {
	UserName string `json:"user_name" form:"user_name"`
	Password string `json:"password" form:"password"`
}

// Register 用户注册 返回用户Id和token
func (service *CreateUserService) Register(ctx context.Context) (*pb.DouyinUserRegisterResponse, error) {
	var err error
	defer func() {
		//返回时若err!=nil则写入日志
		if err != nil {
			util.LogrusObj.Error("<Register>  ", err, " [be from req]:", service)
		}
	}()

	//数据验证
	if err = util.ValidateUser(service.UserName, service.Password); err != nil {
		return nil, e2.NewError(e2.InvalidParams)
	}

	var user model.User
	userDao := dao.NewUserDao(ctx)

	_, exist, err := userDao.ExistOrNotByUserName(service.UserName)
	if err != nil {
		return nil, e2.NewError(e2.Error)
	}
	//若该用户已存在
	if exist {
		return nil, e2.NewError(e2.ErrorExistUser) //todo: 需添加返回
	}
	// 密码加密
	password, err := bcrypt.GenerateFromPassword([]byte(service.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, e2.NewError(e2.Error) //todo: 需添加返回
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
		return nil, e2.NewError(e2.Error)
	}
	token, err := util.GenerateToken(user.ID, service.UserName, 0)
	if err != nil {
		return nil, e2.NewError(e2.ErrorAuthToken)
	}
	return &pb.DouyinUserRegisterResponse{
		UserId: int64(user.ID),
		Token:  token,
	}, nil
}
