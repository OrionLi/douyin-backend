package handler

import (
	"context"
	"douyin-backend/pkg/pb"
	"douyin-backend/server/video-center/cache"
	"douyin-backend/server/video-center/oss"
	"douyin-backend/server/video-center/pkg/errno"
	"douyin-backend/server/video-center/service"
	"fmt"
	"strconv"
)

type VideoServer struct {
	pb.UnsafeVideoCenterServer
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
	key, err2 := cache.RedisGetKey(ctx, token)
	if err2 != nil {
		convertErr := errno.ConvertErr(err2)
		resp.StatusCode = int32(convertErr.ErrCode)
		resp.StatusMsg = &convertErr.ErrMsg
		return resp, nil
	}
	formatInt := strconv.FormatInt(req.UserId, 10)
	if key != formatInt {
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
		key, err := cache.RedisGetKey(ctx, req.GetToken())
		if err != nil {
			convertErr := errno.ConvertErr(err)
			resp.StatusCode = int32(convertErr.ErrCode)
			resp.StatusMsg = &convertErr.ErrMsg
			return resp, nil
		}
		userId, err = strconv.ParseInt(key, 10, 64)
		if err != nil {
			convertErr := errno.ConvertErr(err)
			resp.StatusCode = int32(convertErr.ErrCode)
			resp.StatusMsg = &convertErr.ErrMsg
			return resp, nil
		}
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
	return resp, nil
}

func (s *VideoServer) PublishAction(ctx context.Context, req *pb.DouyinPublishActionRequest) (*pb.DouyinPublishActionResponse, error) {
	resp := new(pb.DouyinPublishActionResponse)
	token := req.Token
	if len(token) == 0 {
		resp.StatusCode = errno.ParamErrCode
		resp.StatusMsg = &errno.ParamErr.ErrMsg
		return resp, nil
	}
	key, err2 := cache.RedisGetKey(ctx, token)
	if err2 != nil {
		convertErr := errno.ConvertErr(err2)
		resp.StatusCode = int32(convertErr.ErrCode)
		resp.StatusMsg = &convertErr.ErrMsg
		return resp, nil
	}
	AuthorId, err2 := strconv.ParseInt(key, 10, 64)
	if err2 != nil {
		convertErr := errno.ConvertErr(err2)
		resp.StatusCode = int32(convertErr.ErrCode)
		resp.StatusMsg = &convertErr.ErrMsg
		return resp, nil
	}
	////判断AuthorId是否在user表中
	//isExist := dao.ExistID(ctx, AuthorId)
	//if !isExist {
	//	resp.StatusCode = errno.ParamErrCode
	//	resp.StatusMsg = &errno.ParamErr.ErrMsg
	//	return resp, nil
	//}
	//向七牛云存放视频资源
	playUrl, coverUrl, err2 := oss.UploadVideo(ctx, AuthorId, req.Data, req.Title)
	if err2 != nil {
		convertErr := errno.ConvertErr(err2)
		resp.StatusCode = int32(convertErr.ErrCode)
		resp.StatusMsg = &convertErr.ErrMsg
		return resp, nil
	}
	err := service.NewVideoService(ctx).PublishAction(AuthorId, playUrl, coverUrl, req.Title)
	if err != nil {
		convertErr := errno.ConvertErr(err)
		resp.StatusCode = int32(convertErr.ErrCode)
		resp.StatusMsg = &convertErr.ErrMsg
		return resp, nil
	}
	resp.StatusCode = errno.SuccessCode
	resp.StatusMsg = &errno.Success.ErrMsg
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
