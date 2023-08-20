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
)

type RelationService struct {
	pb.UnimplementedRelationServiceServer
}

// IsFollowDict 判断用户是否关注了字典中的用户
func (c RelationService) IsFollowDict(ctx context.Context, req *pb.IsFollowDictReq) (*pb.IsFollowDictRsp, error) {
	var isFollowDict = make(map[string]bool)

	for _, unit := range req.FollowUintList {
		// UserIdList 可能是我关注的人
		isFollow, err := dao.IsFollow(unit.SelfUserId, unit.UserIdList)
		if err != nil {
			log.Errorf(ctx, "IsFollowDict err: %v", err)
			return nil, err
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
		err := dao.FollowAction(req.SelfUserId, req.ToUserId)
		if err != nil {
			return nil, err
		}
	} else {
		log.Infof(ctx, "unfollow action id:%v,toid:%v", req.SelfUserId, req.ToUserId)
		err := dao.UnFollowAction(req.SelfUserId, req.ToUserId)
		if err != nil {
			return nil, err
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
	list, err := RelationFollowList(req.UserId, 1)
	if err != nil {
		log.Errorf(ctx, "GetRelationFollowList error, err:%v", err)
		return nil, err
	}
	return &pb.GetRelationFollowListRsp{
		FollowList: list,
	}, nil
}

// GetRelationFollowerList 获取粉丝列表
func (c RelationService) GetRelationFollowerList(ctx context.Context, req *pb.GetRelationFollowerListReq) (*pb.GetRelationFollowerListRsp, error) {
	list, err := RelationFollowList(req.UserId, 2)
	if err != nil {
		return nil, err
	}
	return &pb.GetRelationFollowerListRsp{
		FollowerList: list,
	}, nil
}

// RelationFollowList  获取关注列表
func RelationFollowList(userID, relationType int64) ([]int64, error) {
	var (
		relationList []*model.Relation
		err          error
	)
	if relationType == 1 {
		// 获取关注者
		relationList, err = dao.GetFollowList(userID)
	} else {
		// 获取被关注者
		relationList, err = dao.GetFollowerList(userID)
	}
	if err != nil {
		return nil, err
	}
	if len(relationList) == 0 {
		return []int64{}, nil
	}
	log.Infof(nil, "relationList: %v", relationList)
	resp := make([]int64, 0)
	for _, relation := range relationList {
		// 关注者
		if relationType == 1 {
			resp = append(resp, relation.Follow)
		} else {
			// 被关注者
			resp = append(resp, relation.Follower)
		}
	}

	return resp, nil
}
