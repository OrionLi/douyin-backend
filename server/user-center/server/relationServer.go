package server

import (
	"context"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"user-center/pkg/util"
	"user-center/service"
)

type RelationServer struct {
	relationService *service.RelationService
	pb.UnimplementedRelationServiceServer
}

func NewRelationServer(relationService *service.RelationService) *RelationServer {
	return &RelationServer{relationService: relationService}
}

func (s *RelationServer) RelationAction(ctx context.Context, req *pb.RelationActionRequest) (*pb.RelationActionResponse, error) {
	parseToken, err1 := util.ParseToken(req.Token)
	if err1 != nil {
		return nil, err1
	}
	var err error
	if req.ActionType == 1 {
		err = s.relationService.Follow(ctx, int64(parseToken.ID), req.ToUserId)
	} else if req.ActionType == 2 {
		err = s.relationService.Unfollow(ctx, int64(parseToken.ID), req.ToUserId)
	}

	if err != nil {
		return &pb.RelationActionResponse{StatusCode: -1, StatusMsg: err.Error()}, nil
	}

	return &pb.RelationActionResponse{StatusCode: 0}, nil
}

func (s *RelationServer) GetFollowList(ctx context.Context, req *pb.GetFollowListRequest) (*pb.GetFollowListResponse, error) {
	parseToken, err := util.ParseToken(req.Token)
	if err != nil {
		return nil, err
	}

	users, err := s.relationService.GetFollowList(ctx, int64(parseToken.ID))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "获取关注列表失败: %v", err)
	}

	respUsers := make([]*pb.User, 0, len(users))
	for _, user := range users {
		respUsers = append(respUsers, &pb.User{
			Id:            int64(user.ID),
			Name:          user.Username,
			FollowCount:   user.FollowCount,
			FollowerCount: user.FanCount,
			IsFollow:      true,
		})
	}

	return &pb.GetFollowListResponse{
		StatusCode: 0,
		UserList:   respUsers,
	}, nil
}

func (s *RelationServer) GetFollowerList(ctx context.Context, req *pb.GetFollowerListRequest) (*pb.GetFollowerListResponse, error) {
	parseToken, err := util.ParseToken(req.Token)
	if err != nil {
		return nil, err
	}

	users, err := s.relationService.GetFollowerList(ctx, int64(parseToken.ID))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "获取粉丝列表失败: %v", err)
	}

	isFollowService := service.IsFollowService{
		UserId:       parseToken.ID,
		FollowUserId: 0,
	}

	respUsers := make([]*pb.User, 0, len(users))
	for _, user := range users {
		isFollowService.FollowUserId = user.ID
		isFollow, _ := isFollowService.IsFollow(ctx)
		respUsers = append(respUsers, &pb.User{
			Id:            int64(user.ID),
			Name:          user.Username,
			FollowCount:   user.FollowCount,
			FollowerCount: user.FanCount,
			IsFollow:      isFollow.GetIsFollow(),
		})
	}

	return &pb.GetFollowerListResponse{
		StatusCode: 0,
		UserList:   respUsers,
	}, nil
}

func (s *RelationServer) GetFriendList(ctx context.Context, req *pb.GetFriendListRequest) (*pb.GetFriendListResponse, error) {
	parseToken, err := util.ParseToken(req.Token)
	if err != nil {
		return nil, err
	}

	users, err := s.relationService.GetFriendList(ctx, int64(parseToken.ID))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "获取好友列表失败: %v", err)
	}

	respUsers := make([]*pb.FriendUser, 0, len(users))
	for _, user := range users {
		respUsers = append(respUsers, &pb.FriendUser{
			User: &pb.User{
				Id:            int64(user.ID),
				Name:          user.Username,
				FollowCount:   user.FollowCount,
				FollowerCount: user.FanCount,
				IsFollow:      true,
			},
		})
	}

	return &pb.GetFriendListResponse{
		StatusCode: 0,
		UserList:   respUsers,
	}, nil
}
