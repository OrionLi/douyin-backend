package service

import (
	"context"
	"video-center/dao"
)

type FavoriteService interface {
	CreateFav(videoId int64, userId int64) error
	DeleteFav(videoId int64, userId int64) error
	IsFav(videoId int64, userId int64) (bool, error)
}

type FavoriteServiceImpl struct {
	ctx context.Context
}

func (f FavoriteServiceImpl) CreateFav(videoId int64, userId int64) error {
	return dao.CreateFav(f.ctx, videoId, userId)
}

func (f FavoriteServiceImpl) DeleteFav(videoId int64, userId int64) error {
	return dao.CreateFav(f.ctx, videoId, userId)
}

func (f FavoriteServiceImpl) IsFav(videoId int64, userId int64) (bool, error) {
	return dao.IsFavorite(f.ctx, videoId, userId)
}

func NewFavoriteService(context context.Context) FavoriteService {
	return &FavoriteServiceImpl{
		ctx: context,
	}
}
