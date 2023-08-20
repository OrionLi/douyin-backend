package rpc

import (
	"context"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"video-center/pkg/errno"
)

func PublishAction(ctx context.Context, req *pb.DouyinPublishActionRequest) error {
	r, err := VideoClient.PublishAction(ctx, req)
	if err != nil {
		return err
	}
	if r.StatusCode != 0 {
		return errno.NewErrno(int64(r.StatusCode), *r.StatusMsg)
	}
	return nil
}
