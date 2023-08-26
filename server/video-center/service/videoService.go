package service

import (
	"context"
	"fmt"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"video-center/cache"
	"video-center/pkg/util"

	"video-center/dao"
)

type VideoService struct {
	ctx context.Context
}

func NewVideoService(context context.Context) *VideoService {
	return &VideoService{ctx: context}
}

// PublishList 查看自己发布的视频，已完整封装
func (s *VideoService) PublishList(authorId int64) ([]*pb.Video, error) {
	videos := make([]*pb.Video, 0)
	videoList, err := dao.QueryVideoListByAuthorId(s.ctx, authorId)
	if err != nil {
		return nil, err
	}
	for _, v := range videoList {
		//查询缓存，更新fav缓存 有变化则改动
		favKey := fmt.Sprintf("favorite:%d", v.Id)
		value, err := cache.RedisGetKey(s.ctx, favKey)
		if err != nil {
			util.LogrusObj.Errorf("Cache Error Key:%s ErrorCause:%s", favKey, err.Error())
		}
		toUint := util.StrToUint(value)
		if toUint != v.FavoriteCount {
			//更新Fav缓存
			err := cache.RedisSetKey(s.ctx, favKey, v.FavoriteCount)
			if err != nil {
				util.LogrusObj.Errorf("Cache Error Key:%s errorCause:%s", favKey, err.Error())
			}
		}
		videos = append(videos, &pb.Video{
			Id: v.Id,
			Author: &pb.User{
				Id: authorId,
			},
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    false,
			Title:         v.Title,
		})
	}
	return videos, nil
}
func (s *VideoService) PublishAction(authorId int64, playURL string, coverURL string, title string) error {
	err := dao.SaveVideo(s.ctx, []*dao.Video{
		{AuthorID: authorId,
			PlayUrl:  playURL,
			CoverUrl: coverURL,
			Title:    title},
	})
	if err != nil {
		return err
	}
	return nil
}
func (s *VideoService) FeedVideoList(lastTime int64, userId int64) ([]*pb.Video, error) {
	videoList, err := dao.QueryVideosByCurrentTime(s.ctx, lastTime, 0, 30)
	var isFav bool
	isFav = false
	videos := make([]*pb.Video, 0)
	if err != nil {
		return videos, err
	}
	//如果登录了传入userId，否则传0表示没有登录
	for _, v := range videoList {
		if userId != 0 { //！=0表示已登录，查询是否是作者的粉丝s
			favorite, err := dao.IsFavorite(context.Background(), v.Id, userId) //查看是否点赞
			if err != nil {
				isFav = false
			}
			isFav = favorite
		}
		//查询缓存，更新fav缓存 有变化则改动
		favKey := fmt.Sprintf("favorite:%d", v.Id)
		value, err := cache.RedisGetKey(s.ctx, favKey)
		if err != nil {
			util.LogrusObj.Errorf("Cache Error Key:%s ErrorCause:%s", favKey, err.Error())
			err := cache.RedisSetKey(s.ctx, favKey, v.FavoriteCount)
			if err != nil {
				util.LogrusObj.Errorf("Cache Error Key:%s errorCause:%s", favKey, err.Error())
			}
		}
		toUint := util.StrToUint(value)
		if toUint != v.FavoriteCount {
			//更新Fav缓存
			err := cache.RedisSetKey(s.ctx, favKey, v.FavoriteCount)
			if err != nil {
				util.LogrusObj.Errorf("Cache Error Key:%s errorCause:%s", favKey, err.Error())
			}
		}
		videos = append(videos, &pb.Video{
			Id: v.Id,
			Author: &pb.User{
				Id: userId,
			},
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    isFav,
			Title:         v.Title,
		})
	}
	return videos, nil
}
