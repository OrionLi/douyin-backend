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
func (s *CommentService) SaveComment(LoginUserId int64, videoId int64, content string) (bool, model.CommentApi, error) {
	// 发表评论
	comment := model.Comment{
		Content:    content,
		CreateDate: time.Now(),
		UserId:     LoginUserId,
		VideoId:    videoId,
	}
	b, err := dao.SaveComment(comment)
	if err != nil {
		return false, model.CommentApi{}, err
	}
	if !b {
		//todo
		return false, model.CommentApi{}, nil
	}
	//todo 通过用户ID查询用户然信息后封装
	user := model.User{ID: comment.UserId}

	CommentApi := model.ConvertToCommentApi(comment, user)
	return true, CommentApi, nil
}

// DeleteComment 删除评论
func (s *CommentService) DeleteComment(LoginUserId int64, videoId int64, commentId int64) (bool, model.CommentApi, error) {
	isComment, err := dao.IsUserComment(LoginUserId, commentId, videoId)
	if err != nil {
		return false, model.CommentApi{}, err
	}
	if isComment {
		comment := model.Comment{
			VideoId: videoId,
			ID:      commentId,
		}
		b, err := dao.DeleteComment(comment)
		if err != nil || !b {
			return false, model.CommentApi{}, err
		}
		//该评论存在而且正确删除
		//todo 通过ID查询用户信息
		user := model.User{ID: LoginUserId}
		return true, model.ConvertToCommentApi(comment, user), nil
	} else {
		//该评论不是你的
		return false, model.CommentApi{}, nil
	}
}

// ListComment 查看所有评论
func (s *CommentService) ListComment(videoId int64) ([]model.CommentApi, error) {
	comments, err := dao.CommentList(videoId)
	if err != nil {
		return []model.CommentApi{}, err
	}
	var commentApis []model.CommentApi
	for _, comment := range comments {
		//每个评论里面的userID comment.UserId
		//todo 利用grpc查询user信息
		user := model.User{ID: comment.UserId}

		CommentApi := model.ConvertToCommentApi(comment, user)
		commentApis = append(commentApis, CommentApi)
	}
	return commentApis, nil
}
