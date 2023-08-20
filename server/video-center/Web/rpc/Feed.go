package rpc

import (
	"context"
	"video-center/pkg/errno"
	"video-center/pkg/pb"
)

func Feed(ctx context.Context, req *pb.DouyinFeedRequest) ([]*pb.Video, int64, error) {
	videos := make([]*pb.Video, 0)
	r, err := VideoClient.Feed(ctx, req)
	if err != nil {
		return videos, 0, err
	}
	if r.StatusCode != 0 {
		return videos, 0, errno.NewErrno(int64(r.StatusCode), *r.StatusMsg)
	}
	videos = r.VideoList
	return videos, *r.NextTime, nil
}
