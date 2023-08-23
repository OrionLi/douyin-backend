package handler

import (
	"context"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"sync"
	"time"
	"video-center/dao"
	"video-center/model"
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
	// TODO service调用dao
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
	//favs := dao.ListFav(ctx, userId)
	//if len(favs) == 0 {
	//	return &pb.DouyinFavoriteListResponse{
	//		StatusCode: errno.FavListEmptyCode,
	//		StatusMsg:  errno.FavListEmptyErr.ErrMsg,
	//		VideoList:  nil,
	//	}, errno.NewErrno(errno.FavListEmptyCode, errno.FavListEmptyErr.ErrMsg)
	//}
	////pb
	//var favVideoList []*pb.Video
	//for _, v := range favs {
	//	//todo 得到用户ID，然后调用rpc查询用户信息
	//	//var user := xxx(v.AuthorID) //然后修改Author
	//	video := &pb.Video{
	//		Id:            v.Id,
	//		Author:        nil,
	//		PlayUrl:       v.PlayUrl,
	//		CoverUrl:      v.CoverUrl,
	//		FavoriteCount: v.FavoriteCount,
	//		CommentCount:  v.CommentCount,
	//		IsFavorite:    true,
	//		Title:         v.Title,
	//	}
	//	favVideoList = append(favVideoList, video)
	//}
	//return &pb.DouyinFavoriteListResponse{
	//	StatusCode: 0,
	//	StatusMsg:  errno.Success.ErrMsg,
	//	VideoList:  favVideoList,
	//}, nil
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
		service.NewCommentService(ctx).SaveComment(userId, videoId, commentText)
		/////////////////////////////////////////////////////////////////////
		comment := model.Comment{
			Content:    commentText,
			CreateDate: time.Now(),
			UserId:     userId,
			VideoId:    videoId,
		}
		// TODO service调用
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

func (v VideoInteractionServer) ListComment(ctx context.Context, request *pb.DouyinCommentListRequest) (*pb.DouyinCommentListResponse, error) {
	videoId := request.GetVideoId()
	// TODO service调用
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
