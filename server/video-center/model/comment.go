package model

import (
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"time"
)

//type ApifoxModel struct {
//	CommentList []Comment `json:"comment_list"` // 评论列表
//	StatusCode  int64     `json:"status_code"`  // 状态码，0-成功，其他值-失败
//	StatusMsg   *string   `json:"status_msg"`   // 返回状态描述
//}

type Comment struct {
	Content    string    `json:"content"`     // 评论内容
	CreateDate time.Time `json:"create_date"` // 评论发布日期，格式 mm-dd
	ID         int64     `json:"id"`          // 评论ID
	UserId     int64     `json:"user"`        // 评论用户信息
	VideoId    int64     `json:"-"`           // 视频ID 忽略字段并指定列名
}

// TableName 对应数据库中的表名是comment
func (Comment) TableName() string {
	return "comment"
}

//type CommentAction struct {
//	ActionType  string  `json:"action_type"`            // 1-发布评论，2-删除评论
//	CommentID   *string `json:"comment_id,omitempty"`   // 要删除的评论id，在action_type=2的时候使用
//	CommentText *string `json:"comment_text,omitempty"` // 用户填写的评论内容，在action_type=1的时候使用
//	Token       string  `json:"token"`                  // 用户鉴权token
//	VideoID     string  `json:"video_id"`               // 视频id
//}

//// CommentApi 评论
//type CommentApi struct {
//	Content    string `json:"content"`     // 评论内容
//	CreateDate string `json:"create_date"` // 评论发布日期，格式 mm-dd
//	ID         int64  `json:"id"`          // 评论id
//	User       User   `json:"user"`        // 评论用户信息
//}

// User 评论用户信息
//type User struct {
//	Avatar          string `json:"avatar"`           // 用户头像
//	BackgroundImage string `json:"background_image"` // 用户个人页顶部大图
//	FavoriteCount   int64  `json:"favorite_count"`   // 喜欢数
//	FollowCount     int64  `json:"follow_count"`     // 关注总数
//	FollowerCount   int64  `json:"follower_count"`   // 粉丝总数
//	ID              int64  `json:"id"`               // 用户id
//	IsFollow        bool   `json:"is_follow"`        // true-已关注，false-未关注
//	Name            string `json:"name"`             // 用户名称
//	Signature       string `json:"signature"`        // 个人简介
//	TotalFavorited  string `json:"total_favorited"`  // 获赞数量
//	WorkCount       int64  `json:"work_count"`       // 作品数
//}

func ConvertToCommentApi(comment Comment, user *pb.User) pb.Comment {
	return pb.Comment{
		Id:         comment.ID,
		User:       user,
		Content:    comment.Content,
		CreateDate: comment.CreateDate.Format("01-02"),
	}
}
