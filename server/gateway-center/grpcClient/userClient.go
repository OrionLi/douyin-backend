package grpcClient

import (
	"context"
	"github.com/OrionLi/douyin-backend/pkg/pb"
)

func UserRegister(ctx context.Context, username, password string) (*pb.DouyinUserRegisterResponse, error) {

	resp, err := UserClient.Register(ctx, &pb.DouyinUserRegisterRequest{
		Username: username,
		Password: password,
	})
	return resp, err
}

func UserLogin(ctx context.Context, username, password string) (*pb.DouyinUserLoginResponse, error) {

	resp, err := UserClient.Login(ctx, &pb.DouyinUserLoginRequest{
		Username: username,
		Password: password,
	})
	return resp, err
}

func GetUserById(ctx context.Context, myId, uId uint, token string) (*pb.DouyinUserResponse, error) {

	resp, err := UserClient.GetUserById(ctx, &pb.DouyinUserRequest{
		UserId:   int64(myId),
		FollowId: int64(uId),
		Token:    token,
	})
	return resp, err
}
