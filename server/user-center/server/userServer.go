package server

import (
	"context"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"user-center/service"
)

type UserRPCServer struct {
	pb.UnimplementedUserServiceServer
}

func NewUserRPCServer() *UserRPCServer {
	return &UserRPCServer{}
}

// GetUserById 通过id获取用户基本信息
func (s *UserRPCServer) GetUserById(ctx context.Context, req *pb.DouyinUserRequest) (*pb.DouyinUserResponse, error) {
	userReq := service.GetUserByIdService{Id: uint(req.GetUserId())}
	isFollow := service.IsFollowService{
		UserId:       uint(req.GetUserId()),
		FollowUserId: uint(req.GetFollowId()),
	}
	count:=&
	return
}

// Register 用户注册
func (s *UserRPCServer) Register(ctx context.Context, req *pb.DouyinUserRegisterRequest) (*pb.DouyinUserRegisterResponse, error) {
	regReq := service.CreateUserService{
		UserName: req.GetUsername(),
		Password: req.GetPassword(),
	}
	return regReq.Register(ctx)
}

// Login 用户登录
func (s *UserRPCServer) Login(ctx context.Context, req *pb.DouyinUserLoginRequest) (*pb.DouyinUserLoginResponse, error) {
	loginReq := service.LoginUserService{
		UserName: req.GetUsername(),
		Password: req.GetPassword(),
	}

	return loginReq.Login(ctx)
}

// IsFollow 判断是否关注
func (s *UserRPCServer) IsFollow(ctx context.Context, req *pb.IsFollowRequest) (*pb.IsFollowResponse, error) {
	isFollowReq := service.IsFollowService{
		UserId:       uint(req.GetUserId()),
		FollowUserId: uint(req.GetFollowUserId()),
	}
	return isFollowReq.IsFollow(ctx)
}
