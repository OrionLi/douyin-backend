package dao

import (
	"context"
	"gorm.io/gorm"
	"video-center/model"
)

// SaveComment 发布评论
func SaveComment(ctx context.Context, comment model.Comment) (bool, error) {
	result := DB.WithContext(ctx).Create(&comment)
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected > 0 {
		// 更新视频评论计数
		err := UpdateVideoCommentCount(ctx, comment.VideoId, 1)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

// DeleteComment 删除评论
func DeleteComment(ctx context.Context, comment model.Comment) (bool, error) {
	result := DB.WithContext(ctx).Delete(&comment)
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected > 0 {
		// 更新视频评论计数
		err := UpdateVideoCommentCount(ctx, comment.VideoId, -1)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

// UpdateVideoCommentCount 更新视频的评论计数
func UpdateVideoCommentCount(ctx context.Context, videoID int64, increment int) error {
	result := DB.WithContext(ctx).Model(Video{}).Where("id = ?", videoID).UpdateColumn("comment_count", gorm.Expr("comment_count + ?", increment))
	return result.Error
}

// CommentList 根据视频ID查看所有评论
func CommentList(ctx context.Context, videoId int64) ([]model.Comment, error) {
	var comments []model.Comment
	err := DB.Where("video_id = ?", videoId).Find(&comments).Error
	if err == nil {
		return comments, nil
	} else {
		return nil, err
	}
}

// IsUserComment 该评论是否为用户发布的 是返回true
// 是否找到评论 评论是否为该用户的
func IsUserComment(ctx context.Context, userId int64, commentId int64, videoId int64) (bool, bool, model.Comment, error) {
	var comment model.Comment
	// 查询数据库，找到指定的评论
	if err := DB.WithContext(ctx).Where("id = ?  AND video_id = ?", commentId, videoId).First(&comment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, false, comment, nil // 没有找到匹配的评论
		}
		return true, false, comment, err // 查询出错
	}
	if comment.UserId == userId {
		return true, true, comment, nil // 找到匹配的评论并验证通过
	}
	return true, false, comment, nil
}
