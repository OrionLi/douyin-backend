package server

import (
	"context"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"user-center/grpc"
	"user-center/pkg/e"
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
	// 用户基本信息请求体
	userReq := service.GetUserByIdService{Id: uint(req.GetFollowId())}
	user, err := userReq.GetUserById(ctx)
	if err != nil {
		return nil, err
	}
	// 是否关注请求体
	isFollowreq := service.IsFollowService{
		UserId:       uint(req.GetUserId()),
		FollowUserId: uint(req.GetFollowId()),
	}
	isFollow, err := isFollowreq.IsFollow(ctx)
	if err != nil {
		return nil, err
	}
	// 获取用户关注数与被关注数
	favCount, err := grpc.GetFavCount(ctx, uint(req.GetFollowId()))
	if err != nil {
		return nil, err
	}
	if favCount.StatusCode != 0 {
		return nil, e.NewError(e.Error)
	}
	// 获取用户发布视频列表
	vids, err := grpc.GetPublishList(ctx, uint(req.GetFollowId()), req.GetToken())
	vidCount := len(vids.VideoList)
	userInfo := user.GetUser()
	return &pb.DouyinUserResponse{User: &pb.User{
		Id:            req.FollowId,
		Name:          userInfo.GetName(),
		FollowCount:   userInfo.FollowCount,
		FollowerCount: userInfo.FollowerCount,
		IsFollow:      isFollow.GetIsFollow(),
		FavCount:      int64(favCount.GetFavCount()),
		WorkCount:     int64(vidCount),
		GetFavCount_:  favCount.GetFavCount_,
	}}, nil
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
