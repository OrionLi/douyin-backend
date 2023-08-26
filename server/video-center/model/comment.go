package model

import (
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"time"
)

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

func ConvertToCommentApi(comment Comment, user *pb.User) pb.Comment {
	return pb.Comment{
		Id:         comment.ID,
		User:       user,
		Content:    comment.Content,
		CreateDate: comment.CreateDate.Format("01-02"),
	}
}
