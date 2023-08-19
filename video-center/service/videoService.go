package service

import (
	"context"
	"douyin-backend/video-center/dao"
	"douyin-backend/video-center/generated/video"
)

type VideoService struct {
	ctx context.Context
}

func NewVideoService(context context.Context) *VideoService {
	return &VideoService{ctx: context}
}

// PublishList 查看自己发布的视频，已完整封装
func (s *VideoService) PublishList(authorId int64) ([]*video.Video, error) {
	videos := make([]*video.Video, 0)
	videoList, err := dao.QueryVideoListByAuthorId(s.ctx, authorId)
	if err != nil {
		return nil, err
	}
	//根据userId封装Video，以及是否点赞
	user, err := dao.QueryUserByID(context.Background(), authorId)

	if err != nil {
		return nil, err
	}
	for _, v := range videoList {
		videos = append(videos, &video.Video{
			Id: v.Id,
			Author: &video.User{
				Id:            user.Id,
				Name:          user.Username,
				FollowCount:   &user.FollowCount,
				FollowerCount: &user.FanCount,
				IsFollow:      false,
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
func (s *VideoService) FeedVideoList(lastTime int64, userId int64) ([]*video.Video, error) {
	videoList, err := dao.QueryVideosByCurrentTime(s.ctx, lastTime, 0, 30)
	var isFan bool
	var isFav bool
	isFan = false
	isFav = false
	videos := make([]*video.Video, 0)
	if err != nil {
		return videos, err
	}
	//如果登录了传入userId，否则传0表示没有登录
	for _, v := range videoList {
		if userId != 0 { //！=0表示已登录，查询是否是作者的粉丝
			isFanof, err := dao.IsFanOf(v.AuthorID, uint(userId))
			if err != nil {
				isFan = false
			}
			isFan = isFanof
			favorite, err := dao.IsFavorite(context.Background(), v.Id, userId)
			if err != nil {
				isFav = false
			}
			isFav = favorite
		}
		user, err := dao.QueryUserByID(context.Background(), v.AuthorID)
		if err != nil {
			continue
		}

		videos = append(videos, &video.Video{
			Id: v.Id,
			Author: &video.User{
				Id:            user.Id,
				Name:          user.Username,
				FollowCount:   &user.FollowCount,
				FollowerCount: &user.FanCount,
				IsFollow:      isFan,
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
