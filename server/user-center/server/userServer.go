package server

import (
	"context"
	"fmt"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"user-center/cache"
	"user-center/grpc"
	"user-center/pkg/e"
	"user-center/pkg/util"
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
	//游客浏览
	if req.GetToken() == "" {
		userReq := service.GetUserByIdService{Id: uint(req.GetFollowId())}
		return userReq.GetUserById(ctx)
	}

	// 用户基本信息请求体
	videoCountKey := fmt.Sprintf("publishlist:%d", req.UserId)
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

	//用户作品数量缓存添加缓存
	key, err := cache.RedisGetKey(ctx, videoCountKey)
	var vidCount int64
	if err != nil {
		vids, err := grpc.GetPublishList(ctx, uint(req.GetFollowId()), req.GetToken())
		if err != nil {
			return nil, e.NewError(e.Error)
		}
		vidCount := len(vids.VideoList)
		err = cache.RedisSetKey(ctx, videoCountKey, vidCount)
	}
	vidCount = util.StrToInt64(key)
	// 获取用户发布视频列表
	userInfo := user.GetUser()
	return &pb.DouyinUserResponse{User: &pb.User{
		Id:              req.FollowId,
		Name:            userInfo.GetName(),
		FollowCount:     userInfo.FollowCount,
		FollowerCount:   userInfo.FollowerCount,
		IsFollow:        isFollow.GetIsFollow(),
		FavCount:        int64(favCount.GetFavCount()),
		WorkCount:       vidCount,
		TotalFavorited:  favCount.GetFavCount_,
		BackgroundImage: "https://ts2.cn.mm.bing.net/th?id=OIP-C.HfZqICAPqMQslH0cMrIDFQHaKe&w=210&h=297&c=8&rs=1&qlt=90&o=6&dpr=1.1&pid=3.1&rm=2",
		Signature:       "测试用户",
		Avatar:          "https://ts2.cn.mm.bing.net/th?id=OIP-C.druUEHdZrBEuZPn2w80Y1QHaNK&w=187&h=333&c=8&rs=1&qlt=90&o=6&dpr=1.1&pid=3.1&rm=2",
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
