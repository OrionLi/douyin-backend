package grpc

import (
	"context"
	"github.com/OrionLi/douyin-backend/pkg/pb"
)

func GetFavCount(ctx context.Context, uId uint) (*pb.DouyinFavoriteCountResponse, error) {

	resp, err := VideoInteractionClient.CountFavorite(ctx, &pb.DouyinFavoriteCountRequest{UserId: int64(uId)})
	return resp, err
}

func GetPublishList(ctx context.Context, uId uint, token string) (*pb.DouyinPublishListResponse, error) {
	resp, err := VideoClient.PublishList(ctx, &pb.DouyinPublishListRequest{
		UserId: int64(uId),
		Token:  token,
	})
	return resp, err
}
