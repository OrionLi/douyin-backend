package response

import (
	"encoding/json"
	"fmt"
	"github.com/OrionLi/douyin-backend/pkg/pb"
)

type Video struct {
	Id            int64  `json:"id"`
	User          Vuser  `json:"user"`
	PlayUrl       string `json:"playUrl"`
	CoverUrl      string `json:"coverUrl"`
	FavoriteCount int64  `json:"favoriteCount"`
	CommentCount  int64  `json:"commentCount"`
	IsFavorite    bool   `json:"isFavorite"`
	Title         string `json:"title"`
}

type Vuser struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	FollowCount   int64  `json:"followCount"`
	FollowerCount int64  `json:"followerCount"`
	IsFollow      bool   `json:"isFollow"`
}
type VideoArray []Video

func (v *VideoArray) MarshalBinary() (data []byte, err error) {
	fmt.Println("MarshalBinary")
	return json.Marshal(v)
}
func (v *VideoArray) UnmarshalBinary(data []byte) error {
	fmt.Println("UnmarshalBinary")
	return json.Unmarshal(data, v)

}

type VBResponse struct {
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
	LatestTime int64  `json:"latest_time" form:"latest_time"`
	Token      string `json:"token" form:"token"`
}
type FeedResponse struct {
	VBResponse
	VideoList VideoArray `json:"video_list,omitempty"`
	NextTime  int64      `json:"next_time,omitempty"`
}
type PublishListResponse struct {
	VBResponse
	VideoList VideoArray `json:"video_list,omitempty"`
}
type PublishActionResponse struct {
	VBResponse
}

type CommentActionParam struct {
	ActionType  string `form:"action_type" json:"action_type"`   // 1-发布评论，2-删除评论
	CommentID   string `form:"comment_id" json:"comment_id"`     // 要删除的评论id，在action_type=2的时候使用
	CommentText string `form:"comment_text" json:"comment_text"` // 用户填写的评论内容，在action_type=1的时候使用
	Token       string `form:"token" json:"token"`               // 用户鉴权token
	VideoID     string `form:"video_id" json:"video_id"`         // 视频id
}
type CommentActionResponse struct {
	VBResponse
	Comment *pb.Comment `json:"comment"`
}

type CommentListParam struct {
	Token   string `form:"token" json:"token"`       // 用户鉴权token
	VideoID string `form:"video_id" json:"video_id"` // 视频id
}
type CommentListResponse struct {
	VBResponse
	Comment []*pb.Comment `json:"comment_list"`
}
type FavListResponse struct {
	VBResponse
	FavList VideoArray `form:"video_list" json:"video_list"`
}

type FavActionResponse struct {
	VBResponse
}
