package service

import (
	"context"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"user-center/cache"
	"user-center/dao"
	e2 "user-center/pkg/e"
	userUtil "user-center/pkg/util"
)

// IsFollowService 关注否服务请求
type IsFollowService struct {
	UserId       uint `json:"user_id"  form:"user_id"`
	FollowUserId uint `json:"follow_user_id" form:"follow_user_id"`
}

// IsFollow 判断是否关注 返回 true OR false
func (service *IsFollowService) IsFollow(ctx context.Context) (*pb.IsFollowResponse, error) { //todo: 添加返回结构体
	userDao := dao.NewUserDao(ctx)
	userCache := cache.NewUserCache(ctx)
	var err error

	defer func() {
		//返回时若err!=nil则写入日志
		if err != nil {
			userUtil.LogrusObj.Error("<IsFollow> ", err, " [be from req]:", service)
		}
	}()
	//查找缓存中是否存在
	if userCache.IsFollow(ctx, service.UserId, service.FollowUserId) == true {
		return &pb.IsFollowResponse{IsFollow: true}, nil
	}
	//查找缓存中是否存在
	exist, err := userDao.IsFollow(service.UserId, service.FollowUserId)
	if err != nil {
		return nil, e2.NewError(e2.Error)
	}
	if exist == true {
		//将关系存入缓存
		err = userCache.AddFollow(ctx, service.UserId, service.FollowUserId)
		if err != nil {
			return nil, e2.NewError(e2.Error)
		}
		return &pb.IsFollowResponse{IsFollow: true}, nil
	}
	return &pb.IsFollowResponse{IsFollow: false}, nil
}
