package service

import (
	"context"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	userCache "user-center/cache"
	"user-center/dao"
	e2 "user-center/pkg/e"
	userUtil "user-center/pkg/util"
)

// GetUserByIdService 获取用户信息服务请求
type GetUserByIdService struct {
	Id uint `json:"id" form:"id"`
}

// GetUserById 获取用户基本信息 返回Id,Name,关注总数和粉丝总数
func (service *GetUserByIdService) GetUserById(ctx context.Context) (*pb.DouyinUserResponse, error) { //todo: 添加返回结构体
	cache := userCache.NewUserCache(ctx)
	relationDao := dao.NewRelationDao(ctx)
	var err error

	defer func() {
		//返回时若err!=nil则写入日志
		if err != nil {
			userUtil.LogrusObj.Error("<getUserById> : ", err, " [be from req]:", service)
		}
	}()

	cacheData, err := cache.HasUser(ctx, service.Id)
	if err != nil {
		return nil, e2.NewError(e2.Error)
	}
	if len(cacheData) != 0 { //若缓存存在该记录
		id := service.Id
		name := cacheData["Name"]

		followCount, _ := GetFollowCount(relationDao, cache, int64(id))
		fanCount, _ := GetFollowerCount(relationDao, cache, int64(id))

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

		return nil, e2.NewError(e2.Error)
	}
	if user.IsCelebrity() == true {
		m := map[string]interface{}{
			"Id":          user.ID,
			"Name":        user.Username,
			"FollowCount": user.FollowCount,
			"FanCount":    user.FanCount,
		}
		//将用户信息添加至缓存
		err = cache.AddUser(ctx, user.ID, m)
		if err != nil {
			return nil, e2.NewError(e2.Error)
		}
	}
	return &pb.DouyinUserResponse{User: &pb.User{
		Id:            int64(user.ID),
		Name:          user.Username,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FanCount,
	}}, nil

}
