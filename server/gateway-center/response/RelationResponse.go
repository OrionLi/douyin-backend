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
	StatusCode int32        `json:"status_code"`
	StatusMsg  string       `json:"status_msg"`
	UserList   []FriendUser `json:"user_list"`
}

type FriendUser struct {
	Avatar          string `json:"avatar"`           // 用户头像
	BackgroundImage string `json:"background_image"` // 用户个人页顶部大图
	FavoriteCount   int64  `json:"favorite_count"`   // 喜欢数
	FollowCount     int64  `json:"follow_count"`     // 关注总数
	FollowerCount   int64  `json:"follower_count"`   // 粉丝总数
	ID              int64  `json:"id"`               // 用户id
	IsFollow        bool   `json:"is_follow"`        // true-已关注，false-未关注
	Name            string `json:"name"`             // 用户名称
	Signature       string `json:"signature"`        // 个人简介
	TotalFavorited  int64  `json:"total_favorited"`  // 获赞数量
	WorkCount       int64  `json:"work_count"`       // 作品数
	Message         string `json:"message"`
	MsgType         int32  `json:"msg_type"`
}
