package rpc

import (
	"context"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"video-center/pkg/errno"
)

func PublishAction(ctx context.Context, req *pb.DouyinPublishActionRequest) error {
	err := VideoStreamClient.Send(req)
	if err != nil {
		return err
	}
	recv, err := VideoStreamClient.CloseAndRecv()
	ResetVideoStreamClient()
	if err != nil {
		return err
	}
	if recv.StatusCode == 0 {
		return nil
	} else {
		return errno.NewErrno(int64(recv.StatusCode), *recv.StatusMsg)
	}
}
