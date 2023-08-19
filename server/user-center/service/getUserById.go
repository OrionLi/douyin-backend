package service

import (
	"context"
	"fmt"
	"time"
	cache2 "user-center/cache"
	"user-center/dao"
	util2 "user-center/pkg/util"
)

type GetUserByIdService struct {
	Id uint `json:"id" form:"id"`
}

func (service *GetUserByIdService) GetUserById(ctx context.Context) { //todo: 添加返回结构体
	cache := cache2.NewRedisCache(ctx)
	util2.LogrusObj.WithTime(time.Now()).Info("requestId: ", service.Id)
	//todo: 需添加缓存，并添加逻辑：粉丝数大于等于300为网红
	cacheData, err := cache.HasUser(ctx, service.Id)
	if err != nil {
		util2.LogrusObj.Info("err: ", err)
		return
	}
	if len(cacheData) != 0 {
		Id := service.Id
		Name := cacheData["Name"]

		FollowCount := util2.StrToUint(cacheData["FollowCount"])
		FanCount := util2.StrToUint(cacheData["FanCount"])
		fmt.Println("缓存读取：", Id,
			Name,
			FollowCount,
			FanCount)
		return
	}
	userDao := dao.NewUserDao(ctx)
	//获取用户基本信息
	user, err := userDao.GetUserById(service.Id)
	if err != nil {
		util2.LogrusObj.Info("err: ", err)
		return
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
			util2.LogrusObj.Info("err: ", err)
			return
		}
	}

	fmt.Println("user:", user)
}
