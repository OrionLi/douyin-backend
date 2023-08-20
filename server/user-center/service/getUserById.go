package service

import (
	"context"
	cache2 "user-center/cache"
	"user-center/dao"
	"user-center/pb"
	"user-center/pkg/e"
	util2 "user-center/pkg/util"
)

type GetUserByIdService struct {
	Id uint `json:"id" form:"id"`
}

func (service *GetUserByIdService) GetUserById(ctx context.Context) (*pb.DouyinUserResponse, error) { //todo: 添加返回结构体
	cache := cache2.NewRedisCache(ctx)
	var err error

	defer func() {
		//返回时若err!=nil则写入日志
		if err != nil {
			util2.LogrusObj.Error("<getUserById> : ", err, " [be from req]:", service)
		}
	}()

	cacheData, err := cache.HasUser(ctx, service.Id)
	if err != nil {
		return nil, e.NewError(e.Error)
	}
	if len(cacheData) != 0 { //若缓存存在该记录
		id := service.Id
		name := cacheData["Name"]

		followCount := util2.StrToUint(cacheData["FollowCount"])
		fanCount := util2.StrToUint(cacheData["FanCount"])
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

		return nil, e.NewError(e.Error)
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
