package handler

import (
	"context"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"video-center/dao"
	"video-center/pkg/errno"
	"video-center/service"
)

type VideoInteractionServer struct {
	pb.UnsafeDouyinVideoInteractionServiceServer
}

func (v VideoInteractionServer) ActionFavorite(ctx context.Context, request *pb.DouyinFavoriteActionRequest) (*pb.DouyinFavoriteActionResponse, error) {
	selfUserId := request.GetSelfUserId()
	videoId := request.GetVideoId()
	actionType := request.GetActionType()
	switch actionType {
	case 1:
		err := service.NewFavoriteService(ctx).CreateFav(videoId, selfUserId)
		if err != nil {
			return &pb.DouyinFavoriteActionResponse{
				StatusCode: errno.FavActionErrCode,
				StatusMsg:  errno.FavActionErr.ErrMsg,
			}, err
		}
		return &pb.DouyinFavoriteActionResponse{
			StatusCode: errno.SuccessCode,
			StatusMsg:  errno.Success.ErrMsg,
		}, nil
	case 2:
		err := dao.DeleteFav(ctx, selfUserId, videoId)
		if err != nil {
			return &pb.DouyinFavoriteActionResponse{
				StatusCode: errno.FavActionErrCode,
				StatusMsg:  errno.FavActionErr.ErrMsg,
			}, err
		}
		return &pb.DouyinFavoriteActionResponse{
			StatusCode: errno.SuccessCode,
			StatusMsg:  errno.Success.ErrMsg,
		}, nil
	default:
		return &pb.DouyinFavoriteActionResponse{
			StatusCode: errno.ParamErrCode,
			StatusMsg:  errno.ParamErr.ErrMsg,
		}, nil
	}
}

func (v VideoInteractionServer) ListFavorite(ctx context.Context, request *pb.DouyinFavoriteListRequest) (*pb.DouyinFavoriteListResponse, error) {
	userId := request.GetUserId()
	b, favs := service.NewFavoriteService(ctx).ListFav(userId)
	if !b || len(favs) == 0 {
		return &pb.DouyinFavoriteListResponse{
			StatusCode: errno.FavListEmptyCode,
			StatusMsg:  errno.FavListEmptyErr.ErrMsg,
			VideoList:  nil,
		}, errno.NewErrno(errno.FavListEmptyCode, errno.FavListEmptyErr.ErrMsg)
	}
	return &pb.DouyinFavoriteListResponse{
		StatusCode: 0,
		StatusMsg:  errno.Success.ErrMsg,
		VideoList:  favs,
	}, nil
}

func (v VideoInteractionServer) CountFavorite(ctx context.Context, request *pb.DouyinFavoriteCountRequest) (*pb.DouyinFavoriteCountResponse, error) {
	// userId: 要查询赞数量和被赞数量的用户ID
	userId := request.GetUserId()
	favCount, getFavCount, err := service.NewFavoriteService(context.Background()).CountFav(userId)
	if err != nil {
		return &pb.DouyinFavoriteCountResponse{
			StatusCode: errno.FavCountErrCode,
			StatusMsg:  errno.FavCountErr.ErrMsg,
		}, err
	}
	return &pb.DouyinFavoriteCountResponse{
		StatusCode:   errno.SuccessCode,
		StatusMsg:    errno.Success.ErrMsg,
		FavCount:     favCount,
		GetFavCount_: getFavCount,
	}, nil
}

func (v VideoInteractionServer) ActionComment(ctx context.Context, request *pb.DouyinCommentActionRequest) (*pb.DouyinCommentActionResponse, error) {
	actionType := request.GetActionType()
	userId := request.GetSelfUserId()
	videoId := request.GetVideoId()
	if actionType == 0 {
		// 发布评论
		commentText := request.GetCommentText()
		b, comment, err := service.NewCommentService(ctx).SaveComment(userId, videoId, commentText)
		if err != nil {
			return nil, err
		}
		if !b || err != nil {
			return &pb.DouyinCommentActionResponse{
				StatusCode: errno.CommentPostingCode,
				StatusMsg:  errno.CommentPostingErr.ErrMsg,
				Comment:    &pb.Comment{},
			}, err
		}
		return &pb.DouyinCommentActionResponse{
			StatusCode: errno.SuccessCode,
			StatusMsg:  errno.Success.ErrMsg,
			Comment:    comment,
		}, nil
	}
	if actionType == 1 {
		// 删除评论
		commentId := request.GetCommentId()
		b, comment, err := service.NewCommentService(ctx).DeleteComment(userId, videoId, commentId)
		if !b && err != nil {
			//执行出错
			return &pb.DouyinCommentActionResponse{
				StatusCode: errno.DeleteCommentCode,
				StatusMsg:  errno.DeleteCommentErr.ErrMsg,
				Comment:    &pb.Comment{},
			}, err
		}
		if !b && comment.Content == "" && err == nil {
			//评论不存在
			return &pb.DouyinCommentActionResponse{
				StatusCode: errno.NoCommentExistsCode,
				StatusMsg:  errno.NoCommentExistsErr.ErrMsg,
				Comment:    &pb.Comment{},
			}, nil
		}
		if !b && comment.Content != "" && err == nil {
			//不是自己的评论
			//该评论不是你的
			return &pb.DouyinCommentActionResponse{
				StatusCode: errno.NoMyCommentCode,
				StatusMsg:  errno.NoMyCommentErr.ErrMsg,
				Comment:    &pb.Comment{},
			}, nil
		}
		//删除成功
		return &pb.DouyinCommentActionResponse{
			StatusCode: errno.SuccessCode,
			StatusMsg:  errno.Success.ErrMsg,
			Comment:    comment,
		}, nil
	}
	return &pb.DouyinCommentActionResponse{
		StatusCode: errno.ParamErrCode,
		StatusMsg:  errno.ParamErr.ErrMsg,
		Comment:    &pb.Comment{},
	}, nil
}

func (v VideoInteractionServer) ListComment(ctx context.Context, request *pb.DouyinCommentListRequest) (*pb.DouyinCommentListResponse, error) {
	videoId := request.GetVideoId()
	comments, err := service.NewCommentService(ctx).ListComment(videoId)
	if err != nil {
		return &pb.DouyinCommentListResponse{
			StatusCode:  errno.ListCommentCode,
			StatusMsg:   errno.ListCommentErr.ErrMsg,
			CommentList: nil,
		}, err
	}
	if len(comments) == 0 {
		return &pb.DouyinCommentListResponse{
			StatusCode:  errno.NoCommentExistsCode,
			StatusMsg:   errno.NoCommentExistsErr.ErrMsg,
			CommentList: nil,
		}, nil
	}
	return &pb.DouyinCommentListResponse{
		StatusCode:  errno.SuccessCode,
		StatusMsg:   errno.Success.ErrMsg,
		CommentList: comments,
	}, nil
}
