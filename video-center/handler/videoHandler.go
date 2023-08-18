package handler

import (
	"context"
	"douyin-backend/video-center/cache"
	"douyin-backend/video-center/dao"
	"douyin-backend/video-center/generated/video"
	"douyin-backend/video-center/oss"
	"douyin-backend/video-center/pkg/errno"
	"douyin-backend/video-center/service"
	"strconv"
)

type VideoServer struct {
	video.UnsafeVideoCenterServer
}

func (s *VideoServer) PublishList(ctx context.Context, req *video.DouyinPublishListRequest) (*video.DouyinPublishListResponse, error) {
	resp := new(video.DouyinPublishListResponse)
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
	//将返回的video封装为目标video
	videoList := make([]*video.Video, 0)
	user, err2 := dao.QueryUserByID(ctx, req.UserId)
	if err2 != nil {
		convertErr := errno.ConvertErr(err)
		resp.StatusCode = int32(convertErr.ErrCode)
		resp.StatusMsg = &convertErr.ErrMsg
		return resp, nil
	}
	for _, v := range list {
		videoList = append(videoList, &video.Video{
			Id: v.Id,
			//Author: 即为查询到的user
			Author: &video.User{
				Id:            user.Id,
				Name:          user.Username,
				FollowCount:   &user.FollowCount,
				FollowerCount: &user.FanCount,
				IsFollow:      false,
			},
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    v.IsFavorite,
			Title:         v.Title,
		})
	}
	resp.VideoList = videoList
	return resp, nil
}

// Feed 根据分页查询，随机时间种子来实现，每页最多30个
func (s *VideoServer) Feed(ctx context.Context, req *video.DouyinFeedRequest) (*video.DouyinFeedResponse, error) {
	resp := new(video.DouyinFeedResponse)
	hasToken := false
	isfan := false
	var userId int64
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
		hasToken = true
	}
	lastTime := req.LatestTime
	list, err := service.NewVideoService(ctx).FeedVideoList(*lastTime)
	if err != nil {
		convertErr := errno.ConvertErr(err)
		resp.StatusCode = int32(convertErr.ErrCode)
		resp.StatusMsg = &convertErr.ErrMsg
		return resp, nil
	}
	//正常返回
	resp.StatusCode = errno.SuccessCode
	resp.StatusMsg = &errno.Success.ErrMsg
	//将返回的video封装为目标video
	videoList := make([]*video.Video, 0)
	for _, v := range list {
		//根据video查询user,将User信息放入video中
		user, err := dao.QueryUserByID(ctx, v.AuthorID)
		if err != nil {
			continue
		}
		if hasToken { //如果有token，则判断和视频主人是否为粉丝
			of, err := dao.IsFanOf(v.AuthorID, uint(userId))
			if err != nil {
				isfan = false
			}
			isfan = of
		}
		videoList = append(videoList, &video.Video{
			Id: v.Id,
			//Author: 即为查询到的user
			Author: &video.User{
				Id:            user.Id,
				Name:          user.Username,
				FollowCount:   &user.FollowCount,
				FollowerCount: &user.FanCount,
				IsFollow:      isfan,
			},
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    v.IsFavorite,
			Title:         v.Title,
		})
	}
	resp.VideoList = videoList
	return resp, nil
}

func (s *VideoServer) PublishAction(ctx context.Context, req *video.DouyinPublishActionRequest) (*video.DouyinPublishActionResponse, error) {
	resp := new(video.DouyinPublishActionResponse)
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
	//判断AuthorId是否在user表中
	isExist := dao.ExistID(ctx, AuthorId)
	if !isExist {
		resp.StatusCode = errno.ParamErrCode
		resp.StatusMsg = &errno.ParamErr.ErrMsg
		return resp, nil
	}
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
