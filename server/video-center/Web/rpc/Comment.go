package rpc

import (
	"context"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"sync"
	"time"
	"video-center/dao"
	"video-center/model"
	"video-center/pkg/errno"
)

type CommentRPCService struct {
	pb.UnimplementedDouyinVideoInteractionServiceServer
}

func NewCommentRPCService() *CommentRPCService {
	return &CommentRPCService{}
}

func (s *CommentRPCService) ActionComment(ctx context.Context, request *pb.DouyinCommentActionRequest) (*pb.DouyinCommentActionResponse, error) {
	actionType := request.GetActionType()
	userId := request.GetSelfUserId()
	videoId := request.GetVideoId()
	if actionType == 0 {
		// 发布评论
		commentText := request.GetCommentText()
		comment := model.Comment{
			Content:    commentText,
			CreateDate: time.Now(),
			UserId:     userId,
			VideoId:    videoId,
		}
		b, err := dao.SaveComment(ctx, comment)
		if err != nil || !b {
			return &pb.DouyinCommentActionResponse{
				StatusCode: errno.CommentPostingCode,
				StatusMsg:  errno.CommentPostingErr.ErrMsg,
				Comment:    &pb.Comment{},
			}, err
		}
		//todo 通过用户ID查询用户然信息后封装
		user := pb.User{Id: userId}
		commentApi := model.ConvertToCommentApi(comment, &user)

		return &pb.DouyinCommentActionResponse{
			StatusCode: errno.SuccessCode,
			StatusMsg:  errno.CommentPostingErr.ErrMsg,
			Comment:    &commentApi,
		}, nil
	}
	if actionType == 1 {
		// 删除评论
		commentId := request.GetCommentId()
		isUserComment, err := dao.IsUserComment(ctx, userId, commentId, videoId)
		if err != nil {
			return &pb.DouyinCommentActionResponse{
				StatusCode: errno.DeleteCommentCode,
				StatusMsg:  errno.DeleteCommentErr.ErrMsg,
				Comment:    &pb.Comment{},
			}, err
		}
		if isUserComment {
			comment := model.Comment{
				VideoId: videoId,
				ID:      commentId,
			}
			b, err := dao.DeleteComment(ctx, comment)
			if err != nil || !b {
				return &pb.DouyinCommentActionResponse{
					StatusCode: errno.DeleteCommentCode,
					StatusMsg:  errno.DeleteCommentErr.ErrMsg,
					Comment:    &pb.Comment{},
				}, err
			}
			//该评论存在而且正确删除
			//todo 通过ID查询用户信息
			user := pb.User{Id: userId}

			commentApi := model.ConvertToCommentApi(comment, &user)
			return &pb.DouyinCommentActionResponse{
				StatusCode: errno.SuccessCode,
				StatusMsg:  errno.Success.ErrMsg,
				Comment:    &commentApi,
			}, nil
		} else {
			//该评论不是你的
			return &pb.DouyinCommentActionResponse{
				StatusCode: errno.NoMyCommentCode,
				StatusMsg:  errno.NoMyCommentErr.ErrMsg,
				Comment:    &pb.Comment{},
			}, nil
		}
	}
	return &pb.DouyinCommentActionResponse{
		StatusCode: errno.ParamErrCode,
		StatusMsg:  errno.ParamErr.ErrMsg,
		Comment:    &pb.Comment{},
	}, nil
}

func (s *CommentRPCService) ListComment(ctx context.Context, request *pb.DouyinCommentListRequest) (*pb.DouyinCommentListResponse, error) {
	videoId := request.GetVideoId()
	comments, err := dao.CommentList(ctx, videoId)
	if err != nil {
		return &pb.DouyinCommentListResponse{
			StatusCode:  errno.NoCommentExistsCode,
			StatusMsg:   errno.NoCommentExistsErr.ErrMsg,
			CommentList: nil,
		}, err
	}
	if len(comments) == 0 {
		// 没有评论，可以返回一个空切片或者合适的提示
		return &pb.DouyinCommentListResponse{
			StatusCode:  errno.NoCommentExistsCode,
			StatusMsg:   errno.NoCommentExistsErr.ErrMsg,
			CommentList: nil,
		}, nil
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
	return &pb.DouyinCommentListResponse{
		StatusCode:  errno.SuccessCode,
		StatusMsg:   errno.Success.ErrMsg,
		CommentList: commentApis,
	}, nil
}
