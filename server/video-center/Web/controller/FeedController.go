package controller

import (
	"context"
	"fmt"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"video-center/Web/pkg/baseResponse"
	"video-center/Web/rpc"
	"video-center/cache"
	"video-center/pkg/errno"
	"video-center/pkg/util"
)

func Feed(c *gin.Context) {
	fmt.Println("请求Feed")
	var params baseResponse.FeedParam
	var FeedKey string
	var VideoList baseResponse.VideoArray
	var videos []*pb.Video
	var err error
	var nextTime int64
	if err := c.ShouldBindQuery(&params); err != nil {
		convertErr := errno.ConvertErr(err)
		//记录日志
		util.LogrusObj.Errorf("参数绑定错误 URL:%s form %v 错误原因:%s", c.Request.URL, c.Request.RequestURI, convertErr.ErrMsg)
		c.JSON(http.StatusOK, baseResponse.FeedResponse{
			Response: baseResponse.Response{StatusCode: int32(convertErr.ErrCode), StatusMsg: convertErr.ErrMsg},
		})
		return
	}
	isLogin := false
	isFollow := false
	token, err2 := util.ParseToken(params.Token)
	if err2 != nil {
		isLogin = false
	} else {
		isLogin = true
	}
	fmt.Printf("LatestTime: %d, Token: %s\n", params.LatestTime, params.Token)
	if params.LatestTime == 0 {
		params.LatestTime = time.Now().Unix()
	}

	if isLogin {
		FeedKey = fmt.Sprintf("HttpFeed:time:%d userId:%d", params.LatestTime, token.ID)
	} else {
		FeedKey = fmt.Sprintf("HttpFeed:time:%d userId:0", params.LatestTime)
	}

	VideoList, err2 = cache.RedisGetHttpVideoList(context.Background(), FeedKey)
	if err2 != nil {
		//设置从rpc获取
		videos, nextTime, err = rpc.Feed(context.Background(), &pb.DouyinFeedRequest{
			Token:      &params.Token,
			LatestTime: &params.LatestTime,
		})
		if err != nil {
			convertErr := errno.ConvertErr(err)
			util.LogrusObj.Errorf("rpc调用错误 URL:%s 错误原因:%s", c.Request.URL, convertErr.ErrMsg)
			c.JSON(http.StatusOK, baseResponse.FeedResponse{
				Response: baseResponse.Response{StatusCode: int32(convertErr.ErrCode), StatusMsg: convertErr.ErrMsg},
			})
			return
		}
		VideoList = baseResponse.VideoArray{}
		for _, video := range videos {
			///todo 封装User信息
			//info, err := rpc.GetUserInfo(context.Background(), &pb.DouyinUserRequest{UserId: video.Author.Id})
			//if err != nil {
			//	util.LogrusObj.Errorf("获取User失败 UserId:%d UserToken:%d", video.Author.Id, &params.Token)
			//	continue
			//}
			//if isLogin {
			//	isFollow = rpc.IsFollow(context.Background(), &pb.IsFollowRequest{UserId: video.Id, FollowUserId: int64(token.ID)})
			//}
			//user := User{
			//	Id:            info.Id,
			//	Name:          info.Name,
			//	FollowerCount: info.FollowerCount,
			//	FollowCount:   info.FollowCount,
			//	IsFollow:      isFollow,
			//}
			v := &baseResponse.Video{
				Id:            video.Id,
				User:          baseResponse.User{},
				CoverUrl:      video.CoverUrl,
				PlayUrl:       video.PlayUrl,
				FavoriteCount: video.FavoriteCount,
				CommentCount:  video.CommentCount,
				IsFavorite:    video.IsFavorite,
				Title:         video.Title,
			}
			VideoList = append(VideoList, *v)
		}
		err := cache.RedisSetHttpVideoList(context.Background(), FeedKey, VideoList)
		if err != nil {
			util.LogrusObj.Errorf("Cache Error ERRMSG:%s", err.Error())
		}
	}
	fmt.Println(isLogin)
	fmt.Println(isFollow)

	c.JSON(http.StatusOK, baseResponse.FeedResponse{
		Response:  baseResponse.Response{StatusCode: errno.SuccessCode, StatusMsg: errno.Success.ErrMsg},
		VideoList: VideoList,
		NextTime:  nextTime,
	})
}
