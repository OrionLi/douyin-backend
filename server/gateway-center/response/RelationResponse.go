package response

import (
	"github.com/OrionLi/douyin-backend/pkg/pb"
)

type RelationActionResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type GetFollowListResponse struct {
	StatusCode int32      `json:"status_code"`
	StatusMsg  string     `json:"status_msg"`
	UserList   []*pb.User `json:"user_list"`
}

type GetFollowerListResponse struct {
	StatusCode int32      `json:"status_code"`
	StatusMsg  string     `json:"status_msg"`
	UserList   []*pb.User `json:"user_list"`
}

type GetFriendListResponse struct {
	StatusCode int32            `json:"status_code"`
	StatusMsg  string           `json:"status_msg"`
	UserList   []*pb.FriendUser `json:"user_list"`
}

type FriendUser struct {
	User    *pb.User `json:"user"`
	Message string   `json:"message"`
	MsgType int32    `json:"msg_type"`
}
