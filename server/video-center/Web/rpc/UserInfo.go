package rpc

import (
	"context"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"video-center/pkg/util"
)

func GetUserInfo(ctx context.Context, req *pb.DouyinUserRequest) (*pb.User, error) {
	user, err := UserClient.GetUserById(ctx, req)
	if err != nil {
		util.LogrusObj.Errorf("RPC UserClient Error userId %d", req.UserId)
		return &pb.User{}, err
	}
	return user.User, nil
}
