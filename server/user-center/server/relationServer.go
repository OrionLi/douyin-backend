package server

import (
	"context"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	userId, err1 := getUserIdByToken(req.Token)
	if err1 != nil {
		return nil, err1
	}

	var err error
	if req.ActionType == 1 {
		err = s.relationService.Follow(ctx, userId, req.ToUserId)
	} else if req.ActionType == 2 {
		err = s.relationService.Unfollow(ctx, userId, req.ToUserId)
	}

	if err != nil {
		return &pb.RelationActionResponse{StatusCode: -1, StatusMsg: err.Error()}, nil
	}

	return &pb.RelationActionResponse{StatusCode: 0}, nil
}

func (s *RelationServer) GetFollowList(ctx context.Context, req *pb.GetFollowListRequest) (*pb.GetFollowListResponse, error) {
	userId, err := getUserIdByToken(req.Token)
	if err != nil {
		return nil, err
	}

	users, err := s.relationService.GetFollowList(ctx, userId)
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
		})
	}

	return &pb.GetFollowListResponse{
		StatusCode: 0,
		UserList:   respUsers,
	}, nil
}

func (s *RelationServer) GetFollowerList(ctx context.Context, req *pb.GetFollowerListRequest) (*pb.GetFollowerListResponse, error) {
	userId, err := getUserIdByToken(req.Token)
	if err != nil {
		return nil, err
	}

	users, err := s.relationService.GetFollowerList(ctx, userId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "获取粉丝列表失败: %v", err)
	}

	respUsers := make([]*pb.User, 0, len(users))
	for _, user := range users {
		respUsers = append(respUsers, &pb.User{
			Id:            int64(user.ID),
			Name:          user.Username,
			FollowCount:   user.FollowCount,
			FollowerCount: user.FanCount,
		})
	}

	return &pb.GetFollowerListResponse{
		StatusCode: 0,
		UserList:   respUsers,
	}, nil
}

func (s *RelationServer) GetFriendList(ctx context.Context, req *pb.GetFriendListRequest) (*pb.GetFriendListResponse, error) {
	userId, err := getUserIdByToken(req.Token)
	if err != nil {
		return nil, err
	}

	users, err := s.relationService.GetFriendList(ctx, userId)
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
			},
		})
	}

	return &pb.GetFriendListResponse{
		StatusCode: 0,
		UserList:   respUsers,
	}, nil
}

// 用户鉴权逻辑
func getUserIdByToken(token string) (int64, error) {
	// TODO: 实现具体逻辑
	return 1, nil
}
