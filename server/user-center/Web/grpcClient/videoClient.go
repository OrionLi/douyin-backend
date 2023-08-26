package grpcClient

import (
	"context"
	"github.com/OrionLi/douyin-backend/pkg/pb"
)

func GetFavCount(ctx context.Context, uId uint) (*pb.DouyinFavoriteCountResponse, error) {
	conn := GetConn()
	clientVideo := pb.NewDouyinVideoInteractionServiceClient(conn)
	resp, err := clientVideo.CountFavorite(ctx, &pb.DouyinFavoriteCountRequest{UserId: int64(uId)})
	return resp, err
}

func GetPublishList(ctx context.Context, uId uint, token string) (*pb.DouyinPublishListResponse, error) {
	conn := GetConn()
	clientVideo := pb.NewVideoCenterClient(conn)
	resp, err := clientVideo.PublishList(ctx, &pb.DouyinPublishListRequest{
		UserId: int64(uId),
		Token:  token,
	})
	return resp, err
}
