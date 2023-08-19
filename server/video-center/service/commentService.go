package service

import (
	"context"
	"time"
	"video-center/dao"
	"video-center/model"
)

type CommentService struct {
	ctx context.Context
}

func NewCommentService(context context.Context) *CommentService {
	return &CommentService{ctx: context}
}

//todo 返回的东西还没有写好

// SaveComment 发表评论
func (s *CommentService) SaveComment(LoginUserId int64, videoId int64, content string) model.CommentApi {
	// 发表评论
	comment := dao.SaveComment(model.Comment{
		Content:    content,
		CreateDate: time.Now(),
		UserId:     LoginUserId,
		VideoId:    videoId,
	})
	//todo 通过用户ID查询用户然信息后封装
	user := model.User{ID: comment.UserId}

	CommentApi := model.ConvertToCommentApi(comment, user)
	return CommentApi
}

// DeleteComment 删除评论
func (s *CommentService) DeleteComment(LoginUserId int64, videoId int64, commentId int64) (bool, error) {
	isComment, err := dao.IsUserComment(LoginUserId, commentId, videoId)
	if err != nil {
		return false, err
	}
	if isComment {
		dao.DeleteComment(model.Comment{
			VideoId: videoId,
			ID:      commentId,
		})
		return true, nil
	} else {
		return false, nil
	}
}

// ListComment 查看所有评论
func (s *CommentService) ListComment(videoId int64) []model.CommentApi {
	comments := dao.CommentList(videoId)
	var commentApis []model.CommentApi
	for _, comment := range comments {
		//每个评论里面的userID comment.UserId
		//todo 利用grpc查询user信息
		user := model.User{ID: comment.UserId}

		CommentApi := model.ConvertToCommentApi(comment, user)
		commentApis = append(commentApis, CommentApi)
	}
	return commentApis
}
