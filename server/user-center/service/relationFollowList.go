package service

import (
	"context"
	"google.golang.org/appengine/log"
	"user-center/dao"
	"user-center/model"
)

// RelationFollowList  获取关注或粉丝列表
func RelationFollowList(ctx context.Context, userID, relationType int64) ([]int64, error) {
	var (
		relationList []*model.User
		err          error
	)
	if relationType == 1 {
		// 获取关注者
		relationList, err = dao.NewRelationDao(ctx).GetFollowList(userID)
	} else {
		// 获取被关注者
		relationList, err = dao.NewRelationDao(ctx).GetFollowerList(userID)
	}
	if err != nil {
		return nil, err
	}
	if len(relationList) == 0 {
		return []int64{}, nil
	}
	log.Infof(nil, "relationList: %v", relationList)
	resp := make([]int64, 0)
	for _, user := range relationList {
		resp = append(resp, int64(user.ID))
	}

	return resp, nil
}
