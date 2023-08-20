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
