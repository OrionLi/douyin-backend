package rpc

import (
	"context"
	"github.com/OrionLi/douyin-backend/pkg/pb"
)

func ActionFavorite(ctx context.Context, selfUserId int64, videoId int64, actionType int32) (*pb.DouyinFavoriteActionResponse, error) {
	return VideoInteractionClient.ActionFavorite(ctx, &pb.DouyinFavoriteActionRequest{
		SelfUserId: selfUserId,
		VideoId:    videoId,
		ActionType: actionType,
	})
}

func GetFavoriteCount(ctx context.Context, userId int64) (*pb.DouyinFavoriteCountResponse, error) {
	return VideoInteractionClient.CountFavorite(ctx, &pb.DouyinFavoriteCountRequest{
		UserId: userId,
	})
}

func GetFavoriteList(ctx context.Context, request *pb.DouyinFavoriteListRequest) (*pb.DouyinFavoriteListResponse, error) {
	return VideoInteractionClient.ListFavorite(ctx, request)
}
