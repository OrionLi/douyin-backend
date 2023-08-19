package service

import (
	"context"
	"douyin-backend/server/video-center/dao"
)

type VideoService struct {
	ctx context.Context
}

func NewVideoService(context context.Context) *VideoService {
	return &VideoService{ctx: context}
}

func (s *VideoService) PublishList(authorId int64) ([]*dao.Video, error) {
	videoList, err := dao.QueryVideoListByAuthorId(s.ctx, authorId)
	if err != nil {
		return nil, err
	}
	return videoList, nil
}
func (s *VideoService) FeedVideo(lastTime int64) ([]*dao.Video, error) {
	videoList, err := dao.FeedByTime(s.ctx, lastTime)
	if err != nil {
		return nil, err
	}
	return videoList, nil
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
func (s *VideoService) FeedVideoList(lastTime int64) ([]*dao.Video, error) {
	videoList, err := dao.QueryVideosByCurrentTime(s.ctx, lastTime, 0, 30)
	if err != nil {
		return nil, err
	}
	return videoList, nil
}
