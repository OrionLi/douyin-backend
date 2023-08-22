package rpc

import (
	"context"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"video-center/dao"
)

type FavoriteRPCService struct {
	pb.UnimplementedDouyinMessageServiceServer
}

func NewFavoriteRPCService() *FavoriteRPCService {
	return &FavoriteRPCService{}
}

func (s *FavoriteRPCService) ActionFavorite(ctx context.Context, request *pb.DouyinFavoriteActionRequest) (*pb.DouyinFavoriteActionResponse, error) {
	selfUserId := request.GetSelfUserId()
	videoId := request.GetVideoId()
	actionType := request.GetActionType()
	switch actionType {
	case 1:
		err := dao.CreateFav(ctx, selfUserId, videoId)
		if err != nil {
			return &pb.DouyinFavoriteActionResponse{
				StatusCode: 1,
				StatusMsg:  "点赞失败",
			}, err
		}
		return &pb.DouyinFavoriteActionResponse{
			StatusCode: 0,
			StatusMsg:  "点赞成功",
		}, nil
	case 2:
		err := dao.DeleteFav(ctx, selfUserId, videoId)
		if err != nil {
			return &pb.DouyinFavoriteActionResponse{
				StatusCode: 1,
				StatusMsg:  "取消点赞失败",
			}, err
		}
		return &pb.DouyinFavoriteActionResponse{
			StatusCode: 0,
			StatusMsg:  "取消点赞成功",
		}, nil
	default:
		return &pb.DouyinFavoriteActionResponse{
			StatusCode: 1,
			StatusMsg:  "参数有误",
		}, nil
	}
}

func (s *FavoriteRPCService) GetFavoriteCount(ctx context.Context, request *pb.DouyinFavoriteCountRequest) (*pb.DouyinFavoriteCountResponse, error) {
	// userId: 要查询赞数量和被赞数量的用户ID
	userId := request.GetUserId()
	favCount, getFavCount, err := dao.GetFavoriteCount(ctx, userId)
	if err != nil {
		return &pb.DouyinFavoriteCountResponse{
			StatusCode: 1,
			StatusMsg:  "查询失败",
		}, err
	}
	return &pb.DouyinFavoriteCountResponse{
		StatusCode:   0,
		StatusMsg:    "查询成功",
		FavCount:     favCount,
		GetFavCount_: getFavCount,
	}, nil
}

// TODO 点赞列表RPC接口
func (s *FavoriteRPCService) GetFavoriteList(ctx context.Context, request *pb.DouyinFavoriteListRequest) (*pb.DouyinFavoriteListResponse, error) {
	panic("implement me")
}
