package grpcClient

import (
	"context"
	"github.com/OrionLi/douyin-backend/pkg/pb"
)

func UserRegister(ctx context.Context, username, password string) (*pb.DouyinUserRegisterResponse, error) {

	conn := GetConn()

	clientUser := pb.NewUserServiceClient(conn)
	pb.NewRelationServiceClient(conn)

	resp, err := clientUser.Register(ctx, &pb.DouyinUserRegisterRequest{
		Username: username,
		Password: password,
	})
	return resp, err
}

func UserLogin(ctx context.Context, username, password string) (*pb.DouyinUserLoginResponse, error) {

	conn := GetConn()

	clientUser := pb.NewUserServiceClient(conn)
	resp, err := clientUser.Login(ctx, &pb.DouyinUserLoginRequest{
		Username: username,
		Password: password,
	})
	return resp, err
}

func GetUserById(ctx context.Context, uId uint) (*pb.DouyinUserResponse, error) {

	conn := GetConn()
	clientUser := pb.NewUserServiceClient(conn)
	resp, err := clientUser.GetUserById(ctx, &pb.DouyinUserRequest{UserId: int64(uId)})
	return resp, err
}

func IsFollow(ctx context.Context, uId, followId uint) (*pb.IsFollowResponse, error) {

	conn := GetConn()

	clientUser := pb.NewUserServiceClient(conn)
	resp, err := clientUser.IsFollow(ctx, &pb.IsFollowRequest{
		UserId:       int64(uId),
		FollowUserId: int64(followId),
	})
	return resp, err
}
