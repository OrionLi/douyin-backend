package rpc

import (
	"context"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"video-center/pkg/util"
)

func IsFollow(ctx context.Context, req *pb.IsFollowRequest) bool {
	follow, err := UserClient.IsFollow(ctx, req)
	if err != nil {
		util.LogrusObj.Errorf("RPC UserClient Error userId %d followerId %d", req.UserId, req.FollowUserId)
		return false
	}
	return follow.IsFollow
}
