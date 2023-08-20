package controller

import "github.com/OrionLi/douyin-backend/pkg/pb"

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}
type PublishActionParam struct {
	Token string `form:"token"`
	Title string `form:"title"`
}
type PublishListParam struct {
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}
type FeedParam struct {
	LatestTime int64  `json:"latest_time"`
	Token      string `json:"token"`
}
type FeedResponse struct {
	Response
	VideoList []*pb.Video `json:"video_list,omitempty"`
	NextTime  int64       `json:"next_time,omitempty"`
}
type PublishListResponse struct {
	Response
	VideoList []*pb.Video `json:"video_list,omitempty"`
}
type PublishActionResponse struct {
	Response
}

type CommentActionParam struct {
	ActionType  string `form:"action_type" json:"action_type"`   // 1-发布评论，2-删除评论
	CommentID   string `form:"comment_id" json:"comment_id"`     // 要删除的评论id，在action_type=2的时候使用
	CommentText string `form:"comment_text" json:"comment_text"` // 用户填写的评论内容，在action_type=1的时候使用
	Token       string `form:"token" json:"token"`               // 用户鉴权token
	VideoID     string `form:"video_id" json:"video_id"`         // 视频id
}
type CommentActionResponse struct {
	Response
	Comment *pb.Comment `json:"comment"`
}

type CommentListParam struct {
	Token   string `form:"token" json:"token"`       // 用户鉴权token
	VideoID string `form:"video_id" json:"video_id"` // 视频id
}
type CommentListResponse struct {
	Response
	Comment []*pb.Comment `json:"comment_list"`
}
