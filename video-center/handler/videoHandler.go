package handler

import (
	"context"
	"douyin-backend/video-center/generated/video"
	"douyin-backend/video-center/pkg/errno"
	"douyin-backend/video-center/service"
)

func PublishList(ctx context.Context, req *video.DouyinPublishListRequest) (*video.DouyinPublishListResponse, error) {
	resp := new(video.DouyinPublishListResponse)
	//判断Token
	authorId := req.UserId
	list, err := service.NewVideoService(ctx).PublishList(authorId)
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
	//todo 根据video查询user,将User信息放入video中

	for _, v := range list {
		videoList = append(videoList, &video.Video{
			Id: v.Id,
			//todo Author: 即为查询到的user
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

// FeedVideo feed流基础实现
func FeedVideo(ctx context.Context, req *video.DouyinFeedRequest) (*video.DouyinFeedResponse, error) {
	resp := new(video.DouyinFeedResponse)
	//判断Token
	lastTime := req.LatestTime
	list, err := service.NewVideoService(ctx).FeedVideo(*lastTime)
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
		//todo 根据video查询user,将User信息放入video中

		videoList = append(videoList, &video.Video{
			Id: v.Id,
			//todo Author: 即为查询到的user
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

// FeedVideoList 根据分页查询，随机时间种子来实现，每页最多30个
func FeedVideoList(ctx context.Context, req *video.DouyinFeedRequest) (*video.DouyinFeedResponse, error) {
	resp := new(video.DouyinFeedResponse)
	//判断Token
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
		//todo 根据video查询user,将User信息放入video中

		videoList = append(videoList, &video.Video{
			Id: v.Id,
			//todo Author: 即为查询到的user
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

func PublishAction(ctx context.Context, req *video.DouyinPublishActionRequest) (*video.DouyinPublishActionResponse, error) {
	resp := new(video.DouyinPublishActionResponse)
	token := req.Token
	authorId := 1
	if token != "" {
		//todo 对token判断,通过token获取authorId
	}
	//todo 向七牛云存放视频资源
	playUrl := ""
	coverUrl := ""
	err := service.NewVideoService(ctx).PublishAction(int64(authorId), playUrl, coverUrl, req.Title)
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
