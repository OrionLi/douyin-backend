package server

import (
	"context"
	"fmt"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"google.golang.org/appengine/log"
	"strconv"
	"user-center/constant"
	"user-center/dao"
	"user-center/model"
	"user-center/pkg/e"
	"user-center/pkg/util"
)

type RelationService struct {
	pb.UnimplementedRelationServiceServer
}

// IsFollowDict 判断用户是否关注了字典中的用户
func (c RelationService) IsFollowDict(ctx context.Context, req *pb.IsFollowDictReq) (*pb.IsFollowDictRsp, error) {
	var isFollowDict = make(map[string]bool)

	for _, unit := range req.FollowUintList {
		// UserIdList 可能是我关注的人
		isFollow, err := dao.NewRelationDao(ctx).IsFollow(unit.SelfUserId, unit.UserIdList)
		if err != nil {
			util.LogrusObj.Errorf("IsFollowDict err: %v", err)
			return nil, e.NewError(e.ErrorAborted)
		}
		isFollowKey := strconv.FormatInt(unit.SelfUserId, 10) + "_" + strconv.FormatInt(unit.UserIdList, 10)
		isFollowDict[isFollowKey] = isFollow
	}

	return &pb.IsFollowDictRsp{IsFollowDict: isFollowDict}, nil
}

// RelationAction  关注/取消关注
func (c RelationService) RelationAction(ctx context.Context, req *pb.RelationActionReq) (*pb.RelationActionRsp, error) {
	if req.SelfUserId == req.ToUserId {
		return nil, fmt.Errorf("you can't follow yourself")
	}

	if req.ActionType == 1 {
		log.Infof(ctx, "follow action id:%v,toid:%v", req.SelfUserId, req.ToUserId)
		// 关注
		err := dao.NewRelationDao(ctx).FollowAction(req.SelfUserId, req.ToUserId)
		if err != nil {
			return nil, e.NewError(e.ErrorAborted)
		}
	} else {
		log.Infof(ctx, "unfollow action id:%v,toid:%v", req.SelfUserId, req.ToUserId)
		err := dao.NewRelationDao(ctx).UnFollowAction(req.SelfUserId, req.ToUserId)
		if err != nil {
			return nil, e.NewError(e.ErrorAborted)
		}
	}
	return &pb.RelationActionRsp{
		CommonRsp: &pb.CommonResponse{
			Code: constant.SuccessCode,
			Msg:  constant.SuccessMsg,
		},
	}, nil
}

// GetRelationFollowList 获取被关注者列表
func (c RelationService) GetRelationFollowList(ctx context.Context, req *pb.GetRelationFollowListReq) (*pb.GetRelationFollowListRsp, error) {
	list, err := RelationFollowList(ctx, req.UserId, 1)
	if err != nil {
		log.Errorf(ctx, "GetRelationFollowList error, err:%v", err)
		return nil, e.NewError(e.ErrorAborted)
	}
	return &pb.GetRelationFollowListRsp{
		FollowList: list,
	}, nil
}

// GetRelationFollowerList 获取粉丝列表
func (c RelationService) GetRelationFollowerList(ctx context.Context, req *pb.GetRelationFollowerListReq) (*pb.GetRelationFollowerListRsp, error) {
	list, err := RelationFollowList(ctx, req.UserId, 2)
	if err != nil {
		return nil, e.NewError(e.ErrorAborted)
	}
	return &pb.GetRelationFollowerListRsp{
		FollowerList: list,
	}, nil
}

// RelationFollowList  获取关注列表
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
