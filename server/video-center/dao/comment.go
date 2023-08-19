package dao

import (
	"douyin-backend/server/video-center/model"
	"gorm.io/gorm"
)

// SaveComment 发布评论
func SaveComment(comment model.Comment) model.Comment {
	DB.Create(&comment)
	return comment
}

// DeleteComment 删除评论
func DeleteComment(comment model.Comment) {
	//todo 该评论是否为用户发表？
	DB.Delete(&comment)
}

// CommentList 根据视频ID查看所有评论
func CommentList(videoId int64) []model.Comment {
	var comments []model.Comment
	DB.Where("video_id = ?", videoId).Find(&comments)
	return comments
}

// IsUserComment 该评论是否为用户发布的 是返回true
func IsUserComment(userId int64, commentId int64, videoId int64) (bool, error) {
	var comment model.Comment
	// 查询数据库，找到指定的评论
	if err := DB.Where("id = ? AND user_id = ? AND video_id = ?", commentId, userId, videoId).First(&comment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil // 没有找到匹配的评论
		}
		return false, err // 查询出错
	}
	return true, nil // 找到匹配的评论并验证通过
}
