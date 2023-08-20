package controller

import (
	"video-center/pkg/pb"
)

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
	ActionType  string `json:"action_type"`            // 1-发布评论，2-删除评论
	CommentID   string `json:"comment_id,omitempty"`   // 要删除的评论id，在action_type=2的时候使用
	CommentText string `json:"comment_text,omitempty"` // 用户填写的评论内容，在action_type=1的时候使用
	Token       string `json:"token"`                  // 用户鉴权token
	VideoID     string `json:"video_id"`               // 视频id
}
type CommentActionResponse struct {
	Response
	Comment *pb.Comment `json:"comment,omitempty"`
}

type CommentListParam struct {
	Token   string `json:"token"`    // 用户鉴权token
	VideoID string `json:"video_id"` // 视频id
}
type CommentListResponse struct {
	Response
	Comment []*pb.Comment `json:"comment_list,omitempty"`
}
