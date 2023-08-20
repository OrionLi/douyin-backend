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

//todo 返回的东西还没有写好

// SaveComment 发表评论
func (s *CommentService) SaveComment(LoginUserId int64, videoId int64, content string) (bool, *pb.Comment, error) {
	// 发表评论
	comment := model.Comment{
		Content:    content,
		CreateDate: time.Now(),
		UserId:     LoginUserId,
		VideoId:    videoId,
	}
	b, err := dao.SaveComment(comment)
	if err != nil {
		return false, &pb.Comment{}, err
	}
	if !b {
		//todo
		return false, &pb.Comment{}, nil
	}
	//todo 通过用户ID查询用户然信息后封装
	user := pb.User{Id: 1}

	commentApi := model.ConvertToCommentApi(comment, &user)
	return true, &commentApi, nil
}

// DeleteComment 删除评论
func (s *CommentService) DeleteComment(LoginUserId int64, videoId int64, commentId int64) (bool, *pb.Comment, error) {
	isComment, err := dao.IsUserComment(LoginUserId, commentId, videoId)
	if err != nil {
		return false, &pb.Comment{}, err
	}
	if isComment {
		comment := model.Comment{
			VideoId: videoId,
			ID:      commentId,
		}
		b, err := dao.DeleteComment(comment)
		if err != nil || !b {
			return false, &pb.Comment{}, err
		}
		//该评论存在而且正确删除
		//todo 通过ID查询用户信息
		user := pb.User{Id: LoginUserId}
		commentApi := model.ConvertToCommentApi(comment, &user)
		return true, &commentApi, nil
	} else {
		//该评论不是你的
		return false, &pb.Comment{}, nil
	}
}

// ListComment 查看所有评论
func (s *CommentService) ListComment(videoId int64) ([]*pb.Comment, error) {
	comments, err := dao.CommentList(videoId)
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
		//todo 利用grpc查询user信息
		user := pb.User{Id: comment.UserId}

		CommentApi := model.ConvertToCommentApi(comment, &user)
		mu.Lock()
		commentApis = append(commentApis, &CommentApi)
		mu.Unlock()
	}
	return commentApis, nil
}
