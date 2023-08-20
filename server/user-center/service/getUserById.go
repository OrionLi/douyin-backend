package service

import (
	"context"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"time"
	userCache "user-center/cache"
	"user-center/dao"
	"user-center/pkg/e"
	userUtil "user-center/pkg/util"
)

type GetUserByIdService struct {
	Id uint `json:"id" form:"id"`
}

func (service *GetUserByIdService) GetUserById(ctx context.Context) (*pb.DouyinUserResponse, error) { //todo: 添加返回结构体
	cache := userCache.NewRedisCache(ctx)
	userUtil.LogrusObj.WithTime(time.Now()).Info("requestId: ", service.Id)
	//todo: 需添加缓存，并添加逻辑：粉丝数大于等于300为网红
	cacheData, err := cache.HasUser(ctx, service.Id)
	if err != nil {
		userUtil.LogrusObj.Info("err: ", err)
		return nil, e.NewError(e.Error)
	}
	if len(cacheData) != 0 {
		id := service.Id
		name := cacheData["Name"]

		followCount := userUtil.StrToUint(cacheData["FollowCount"])
		fanCount := userUtil.StrToUint(cacheData["FanCount"])
		return &pb.DouyinUserResponse{User: &pb.User{
			Id:            int64(id),
			Name:          name,
			FollowCount:   followCount,
			FollowerCount: fanCount,
		}}, nil
	}
	userDao := dao.NewUserDao(ctx)
	//获取用户基本信息
	user, err := userDao.GetUserById(service.Id)
	if err != nil {
		userUtil.LogrusObj.Info("err: ", err)

		return nil, e.NewError(e.Error)
	}
	if user.IsCelebrity() == true {
		m := map[string]interface{}{
			"Id":          user.ID,
			"Name":        user.Username,
			"FollowCount": user.FollowCount,
			"FanCount":    user.FanCount,
		}
		err = cache.AddUser(ctx, user.ID, m)
		if err != nil {
			userUtil.LogrusObj.Info("err: ", err)
			return nil, e.NewError(e.Error)
		}
	}
	return &pb.DouyinUserResponse{User: &pb.User{
		Id:            int64(user.ID),
		Name:          user.Username,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FanCount,
	}}, nil

}
