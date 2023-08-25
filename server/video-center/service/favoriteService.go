package service

import (
	"context"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"github.com/go-redis/redis/v8"
	"video-center/cache"
	"video-center/dao"
)

type FavoriteService interface {
	CreateFav(videoId int64, userId int64) error
	DeleteFav(videoId int64, userId int64) error
	IsFav(videoId int64, userId int64) (bool, error)
	ListFav(userId int64) (bool, []*pb.Video)
	CountFav(userId int64) (int32, int32, error)
}

type FavoriteServiceImpl struct {
	ctx context.Context
}

func (f FavoriteServiceImpl) CreateFav(videoId int64, userId int64) error {
	// TODO 验证是否已经点赞 此项应在压测通过后实现
	// HACK IsFav() 验证 重复点赞问题
	err := dao.CreateFav(f.ctx, videoId, userId)
	if err != nil {
		return err
	}
	return cache.ActionFavoriteCache(videoId, 1)
}

func (f FavoriteServiceImpl) DeleteFav(videoId int64, userId int64) error {
	// TODO 验证是否未点赞 此项应在压测通过后实现
	// HACK IsFav() 验证 未点赞试图取消点赞问题
	err := dao.DeleteFav(f.ctx, videoId, userId)
	if err != nil {
		return err
	}
	return cache.ActionFavoriteCache(videoId, 2)
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

// TODO 是否考虑缓存穿透问题
func (f FavoriteServiceImpl) CountFav(userId int64) (int32, int32, error) {
	favoriteCount, err := dao.GetFavoriteCount(f.ctx, userId)
	if err != nil {
		return 0, 0, err
	}
	receivedFavoriteCount, err := cache.GetFavoriteCountCache(userId)
	if err != nil {
		if err != redis.Nil {
			return 0, 0, err
		}
		count, err := dao.GetReceivedFavoriteCount(f.ctx, userId)
		if err != nil {
			return 0, 0, err
		}
		return favoriteCount, count, nil
	}
	return favoriteCount, int32(receivedFavoriteCount), nil
}

func NewFavoriteService(context context.Context) FavoriteService {
	return &FavoriteServiceImpl{
		ctx: context,
	}
}
