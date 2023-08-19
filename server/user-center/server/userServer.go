package server

import (
	"context"
	"user-center/pb"
	"user-center/service"
)

type UserRPCServer struct {
	pb.UnimplementedUserServiceServer
}

func NewUserRPCServer() *UserRPCServer {
	return &UserRPCServer{}
}

func (s *UserRPCServer) GetUserById(ctx context.Context, req *pb.DouyinUserRequest) (*pb.DouyinUserResponse, error) {
	userReq := service.GetUserByIdService{Id: uint(req.GetUserId())}
	return userReq.GetUserById(ctx)
}
func (s *UserRPCServer) Register(ctx context.Context, req *pb.DouyinUserRegisterRequest) (*pb.DouyinUserRegisterResponse, error) {
	regReq := service.CreateUserService{
		UserName: req.GetUsername(),
		Password: req.GetPassword(),
	}
	return regReq.Register(ctx)
}
func (s *UserRPCServer) Login(ctx context.Context, req *pb.DouyinUserLoginRequest) (*pb.DouyinUserLoginResponse, error) {
	loginReq := service.LoginUserService{
		UserName: req.GetUsername(),
		Password: req.GetPassword(),
	}
	return loginReq.Login(ctx)
}
func (s *UserRPCServer) IsFollow(ctx context.Context, req *pb.IsFollowRequest) (*pb.IsFollowResponse, error) {
	isFollowReq := service.IsFollowService{
		UserId:       uint(req.GetUserId()),
		FollowUserId: uint(req.GetFollowUserId()),
	}
	return isFollowReq.IsFollow(ctx)
}
