package service

import (
	"context"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"video-center/dao"
)

type FavoriteService interface {
	CreateFav(videoId int64, userId int64) error
	DeleteFav(videoId int64, userId int64) error
	IsFav(videoId int64, userId int64) (bool, error)
	ListFav(userId int64) (bool, []*pb.Video)
}

type FavoriteServiceImpl struct {
	ctx context.Context
}

func (f FavoriteServiceImpl) CreateFav(videoId int64, userId int64) error {
	return dao.CreateFav(f.ctx, videoId, userId)
}

func (f FavoriteServiceImpl) DeleteFav(videoId int64, userId int64) error {
	return dao.DeleteFav(f.ctx, videoId, userId)
}

func (f FavoriteServiceImpl) IsFav(videoId int64, userId int64) (bool, error) {
	return dao.IsFavorite(f.ctx, videoId, userId)
}

func (f FavoriteServiceImpl) ListFav(userId int64) (bool, []*pb.Video) {
	//得到喜欢的视频集合
	favs := dao.ListFav(f.ctx, userId)
	if len(favs) == 0 {
		return false, []*pb.Video{}
	}
	//pb
	var favVideoList []*pb.Video
	for _, v := range favs {
		//todo 得到用户ID，然后调用rpc查询用户信息
		//var user := xxx(v.AuthorID) //然后修改Author
		video := &pb.Video{
			Id:            v.Id,
			Author:        nil,
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    true,
			Title:         v.Title,
		}
		favVideoList = append(favVideoList, video)
	}
	return true, favVideoList
}

func NewFavoriteService(context context.Context) FavoriteService {
	return &FavoriteServiceImpl{
		ctx: context,
	}
}
