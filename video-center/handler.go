package main

import (
	"context"
	videoCenter "douyin-backend/video-center/kitex_gen/videoCenter"
)

// VideoCenterImpl implements the last service interface defined in the IDL.
type VideoCenterImpl struct{}

// Feed implements the VideoCenterImpl interface.
func (s *VideoCenterImpl) Feed(ctx context.Context, req *videoCenter.DouyinFeedRequest) (resp *videoCenter.DouyinFeedResponse, err error) {
	return
}

// PublishList implements the VideoCenterImpl interface.
func (s *VideoCenterImpl) PublishList(ctx context.Context, req *videoCenter.DouyinPublishListRequest) (resp *videoCenter.DouyinPublishListResponse, err error) {
	return
}

// PublishAction implements the VideoCenterImpl interface.
func (s *VideoCenterImpl) PublishAction(ctx context.Context, req *videoCenter.DouyinPublishActionRequest) (resp *videoCenter.DouyinPublishActionResponse, err error) {
	return
}
