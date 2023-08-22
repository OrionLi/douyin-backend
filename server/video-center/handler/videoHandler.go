package handler

import (
	"context"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"io"
	"time"
	"video-center/oss"
	"video-center/pkg/errno"
	"video-center/pkg/util"
	"video-center/service"

	"fmt"
)

type VideoServer struct {
	pb.UnsafeVideoCenterServer
}

func (s *VideoServer) PublishAction(server pb.VideoCenter_PublishActionServer) error {
	for {
		request, err := server.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		user, err := util.ParseToken(request.Token)
		if err != nil {
			fmt.Println(err)
			return err
		}
		playUrl, coverUrl, err := oss.UploadVideo(context.Background(), int64(user.ID), request.Data, request.Title)
		if err != nil {
			fmt.Println(err)
			return err
		}
		err = service.NewVideoService(context.Background()).PublishAction(int64(user.ID), playUrl, coverUrl, request.Title)
		if err != nil {
			fmt.Println(err)
			return err
		}
		response := &pb.DouyinPublishActionResponse{
			StatusCode: errno.SuccessCode,
			StatusMsg:  &errno.Success.ErrMsg,
		}
		if err := server.SendAndClose(response); err != nil {
			return err
		}
	}
	return nil
}

func (s *VideoServer) PublishList(ctx context.Context, req *pb.DouyinPublishListRequest) (*pb.DouyinPublishListResponse, error) {
	resp := new(pb.DouyinPublishListResponse)
	//判断Token
	token := req.Token
	if len(token) == 0 {
		resp.StatusCode = errno.ParamErrCode
		resp.StatusMsg = &errno.ParamErr.ErrMsg
		return resp, nil
	}
	result, err2 := util.ParseToken(token)
	if err2 != nil {
		convertErr := errno.ConvertErr(err2)
		resp.StatusCode = int32(convertErr.ErrCode)
		resp.StatusMsg = &convertErr.ErrMsg
		return resp, nil
	}
	if result.ID != uint(req.UserId) {
		resp.StatusCode = errno.TokenErrCode
		resp.StatusMsg = &errno.TokenErr.ErrMsg
		return resp, nil
	}
	//Token验证成功之后，根据UserId返回
	list, err := service.NewVideoService(ctx).PublishList(req.UserId)
	if err != nil {
		convertErr := errno.ConvertErr(err)
		resp.StatusCode = int32(convertErr.ErrCode)
		resp.StatusMsg = &convertErr.ErrMsg
		return resp, nil
	}
	//正常返回
	resp.StatusCode = errno.SuccessCode
	resp.StatusMsg = &errno.Success.ErrMsg
	resp.VideoList = list
	return resp, nil
}

// Feed 根据分页查询，随机时间种子来实现，每页最多30个
func (s *VideoServer) Feed(ctx context.Context, req *pb.DouyinFeedRequest) (*pb.DouyinFeedResponse, error) {
	resp := new(pb.DouyinFeedResponse)
	var userId int64
	userId = 0
	//判断token,并获取userId
	if len(req.GetToken()) != 0 {
		token, err := util.ParseToken(req.GetToken())
		if err != nil {
			convertErr := errno.ConvertErr(err)
			resp.StatusCode = int32(convertErr.ErrCode)
			resp.StatusMsg = &convertErr.ErrMsg
			return resp, nil
		}
		userId = int64(token.ID)
	}
	lastTime := req.LatestTime
	fmt.Println(userId)
	list, err := service.NewVideoService(ctx).FeedVideoList(*lastTime, userId)
	if err != nil {
		convertErr := errno.ConvertErr(err)
		resp.StatusCode = int32(convertErr.ErrCode)
		resp.StatusMsg = &convertErr.ErrMsg
		return resp, nil
	}
	//正常返回
	resp.StatusCode = errno.SuccessCode
	resp.StatusMsg = &errno.Success.ErrMsg
	resp.VideoList = list
	//获取当前时间
	unix := time.Now().Unix()
	resp.NextTime = &unix
	return resp, nil
}

//无分页
// FeedVideo feed流基础实现
//func FeedVideo(ctx context.Context, req *video.DouyinFeedRequest) (*video.DouyinFeedResponse, error) {
//	resp := new(video.DouyinFeedResponse)
//	//判断Token
//	lastTime := req.LatestTime
//	list, err := service.NewVideoService(ctx).FeedVideo(*lastTime)
//	if err != nil {
//		convertErr := errno.ConvertErr(err)
//		resp.StatusCode = int32(convertErr.ErrCode)
//		resp.StatusMsg = &convertErr.ErrMsg
//		return resp, nil
//	}
//	//正常返回
//	resp.StatusCode = errno.SuccessCode
//	resp.StatusMsg = &errno.Success.ErrMsg
//	//将返回的video封装为目标video
//	videoList := make([]*video.Video, 0)
//	for _, v := range list {
//		// 根据video查询user,将User信息放入video中
//
//		videoList = append(videoList, &video.Video{
//			Id: v.Id,
//			//Author: 即为查询到的user
//			PlayUrl:       v.PlayUrl,
//			CoverUrl:      v.CoverUrl,
//			FavoriteCount: v.FavoriteCount,
//			CommentCount:  v.CommentCount,
//			IsFavorite:    v.IsFavorite,
//			Title:         v.Title,
//		})
//	}
//	resp.VideoList = videoList
//	return resp, nil
//}
