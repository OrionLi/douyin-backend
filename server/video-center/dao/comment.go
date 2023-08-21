package dao

import (
	"gorm.io/gorm"
	"video-center/model"
)

// SaveComment 发布评论
func SaveComment(comment model.Comment) (bool, error) {
	result := DB.Create(&comment)
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected > 0 {
		// 更新视频评论计数
		err := UpdateVideoCommentCount(comment.VideoId, 1)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

// DeleteComment 删除评论
func DeleteComment(comment model.Comment) (bool, error) {
	result := DB.Delete(&comment)
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected > 0 {
		// 更新视频评论计数
		err := UpdateVideoCommentCount(comment.VideoId, -1)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

// UpdateVideoCommentCount 更新视频的评论计数
func UpdateVideoCommentCount(videoID int64, increment int) error {
	result := DB.Model(Video{}).Where("id = ?", videoID).UpdateColumn("comment_count", gorm.Expr("comment_count + ?", increment))
	return result.Error
}

// CommentList 根据视频ID查看所有评论
func CommentList(videoId int64) ([]model.Comment, error) {
	var comments []model.Comment
	err := DB.Where("video_id = ?", videoId).Find(&comments).Error
	if err == nil {
		return comments, nil
	} else {
		return nil, err
	}
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
