package service

import (
	"context"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"sync"
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

// SaveComment 发表评论
func (s *CommentService) SaveComment(LoginUserId int64, videoId int64, content string) (bool, *pb.Comment, error) {
	// 发表评论
	comment := model.Comment{
		Content:    content,
		CreateDate: time.Now(),
		UserId:     LoginUserId,
		VideoId:    videoId,
	}
	b, err := dao.SaveComment(s.ctx, comment)
	if err != nil {
		return false, &pb.Comment{}, err
	}
	if !b {
		return false, &pb.Comment{}, nil
	}
	user := pb.User{Id: LoginUserId}
	commentApi := model.ConvertToCommentApi(comment, &user)
	return true, &commentApi, nil
}

// DeleteComment 删除评论
// bool评论是否是自己发布的
func (s *CommentService) DeleteComment(LoginUserId int64, videoId int64, commentId int64) (bool, *pb.Comment, error) {
	ExistComment, isUserComment, comment, err := dao.IsUserComment(s.ctx, LoginUserId, commentId, videoId)
	if err != nil {
		return false, &pb.Comment{}, err
	}
	if ExistComment {
		//评论存在
		if isUserComment {
			//是自己的评论
			b, err := dao.DeleteComment(s.ctx, comment)
			if err != nil || !b {
				//删除出错
				return false, &pb.Comment{}, err
			}
			//该评论存在而且正确删除
			user := pb.User{Id: LoginUserId}

			commentApi := model.ConvertToCommentApi(comment, &user)
			return true, &commentApi, nil
		} else {
			//不是自己的评论
			return false, &pb.Comment{
				Content: comment.Content,
			}, nil
		}
	} else {
		//评论不存在
		return false, &pb.Comment{}, nil
	}
}

// ListComment 查看所有评论
func (s *CommentService) ListComment(videoId int64) ([]*pb.Comment, error) {
	comments, err := dao.CommentList(s.ctx, videoId)
	if err != nil {
		return []*pb.Comment{}, err
	}
	if len(comments) == 0 {
		// 没有评论，可以返回一个空切片或者合适的提示
		return []*pb.Comment{}, nil
	}
	var mu sync.Mutex
	var commentApis []*pb.Comment
	for _, comment := range comments {
		//每个评论里面的userID comment.UserId
		user := pb.User{Id: comment.UserId}

		CommentApi := model.ConvertToCommentApi(comment, &user)
		mu.Lock()
		commentApis = append(commentApis, &CommentApi)
		mu.Unlock()
	}
	return commentApis, nil
}
