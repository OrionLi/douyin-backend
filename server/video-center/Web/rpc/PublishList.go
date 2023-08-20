package rpc

import (
	"context"
	"video-center/pkg/errno"
	"video-center/pkg/pb"
)

func PublishList(ctx context.Context, req *pb.DouyinPublishListRequest) ([]*pb.Video, error) {
	videos := make([]*pb.Video, 0)
	r, err := VideoClient.PublishList(ctx, req)
	if err != nil {
		return videos, err
	}
	if r.StatusCode != 0 {
		return videos, errno.NewErrno(int64(r.StatusCode), *r.StatusMsg)
	}
	videos = r.VideoList
	return videos, nil
}
