package grpcClient

import (
	"context"
	"github.com/OrionLi/douyin-backend/pkg/pb"
)

func RelationAction(ctx context.Context, token string, toUserId int64, actionType int64) (*pb.RelationActionResponse, error) {
	return RelationClient.RelationAction(ctx, &pb.RelationActionRequest{
		Token:      token,
		ToUserId:   toUserId,
		ActionType: int32(actionType),
	})
}

func GetFollowList(ctx context.Context, userId int64, token string) (*pb.GetFollowListResponse, error) {
	return RelationClient.GetFollowList(ctx, &pb.GetFollowListRequest{
		UserId: userId,
		Token:  token,
	})
}

func GetFollowerList(ctx context.Context, userId int64, token string) (*pb.GetFollowerListResponse, error) {
	return RelationClient.GetFollowerList(ctx, &pb.GetFollowerListRequest{
		UserId: userId,
		Token:  token,
	})
}

func GetFriendList(ctx context.Context, userId int64, token string) (*pb.GetFriendListResponse, error) {
	return RelationClient.GetFriendList(ctx, &pb.GetFriendListRequest{
		UserId: userId,
		Token:  token,
	})
}
