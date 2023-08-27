package service

import (
	"context"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"golang.org/x/crypto/bcrypt"
	"user-center/dao"
	e2 "user-center/pkg/e"
	"user-center/pkg/util"
)

// LoginUserService 用户登录服务请求
type LoginUserService struct {
	UserName string `json:"user_name" form:"user_name"`
	Password string `json:"password" form:"password"`
}

// Login 登录函数，返回用户ID和token
func (service *LoginUserService) Login(ctx context.Context) (*pb.DouyinUserLoginResponse, error) {

	var err error

	defer func() {
		// 返回时若err!=nil则写入日志
		if err != nil {
			util.LogrusObj.Error("<login> ", err, " [be from req]:", service)
		}
	}()
	// 数据验证
	if err = util.ValidateUser(service.UserName, service.Password); err != nil {
		return nil, e2.NewError(e2.InvalidParams)
	}

	// 查询用户是否存在
	userDao := dao.NewUserDao(ctx)
	user, exist, err := userDao.ExistOrNotByUserName(service.UserName)
	if err != nil {
		return nil, e2.NewError(e2.Error)
	}
	if exist == false {
		return nil, e2.NewError(e2.ErrorExistUserNotFound)
	}

	// 比较密码是否匹配
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(service.Password)); err != nil {
		return nil, e2.NewError(e2.ErrorNotCompare)
	}

	// 签发token
	token, err := util.GenerateToken(user.ID, service.UserName, 0)
	if err != nil {
		return nil, e2.NewError(e2.Error)
	}

	return &pb.DouyinUserLoginResponse{
		UserId: int64(user.ID),
		Token:  token,
	}, nil
}
