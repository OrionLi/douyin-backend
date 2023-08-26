package controller

import (
	"context"
	"fmt"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"video-center/Web/rpc"
	"video-center/pkg/errno"
	"video-center/pkg/util"
)

func Feed(c *gin.Context) {
	fmt.Println("请求Feed")
	var params FeedParam
	isLogin := false
	isFollow := false
	if err := c.ShouldBindQuery(&params); err != nil {
		convertErr := errno.ConvertErr(err)
		//记录日志
		util.LogrusObj.Errorf("参数绑定错误 URL:%s form %v 错误原因:%s", c.Request.URL, c.Request.RequestURI, convertErr.ErrMsg)
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{StatusCode: int32(convertErr.ErrCode), StatusMsg: convertErr.ErrMsg},
		})
		return
	}
	token, err2 := util.ParseToken(params.Token)
	if err2 != nil {
		isLogin = false
	}
	if token.ID > 0 {
		isLogin = true
	}
	fmt.Printf("LatestTime: %d, Token: %s\n", params.LatestTime, params.Token)
	if params.LatestTime == 0 {
		params.LatestTime = time.Now().Unix()
	}
	videos, nextTime, err := rpc.Feed(context.Background(), &pb.DouyinFeedRequest{
		Token:      &params.Token,
		LatestTime: &params.LatestTime,
	})
	if err != nil {
		convertErr := errno.ConvertErr(err)
		util.LogrusObj.Errorf("rpc调用错误 URL:%s 错误原因:%s", c.Request.URL, convertErr.ErrMsg)
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{StatusCode: int32(convertErr.ErrCode), StatusMsg: convertErr.ErrMsg},
		})
		return
	}
	videoList := make([]*Video, 0)
	for _, video := range videos {
		info, err := rpc.GetUserInfo(context.Background(), &pb.DouyinUserRequest{UserId: video.Author.Id})
		if err != nil {
			util.LogrusObj.Errorf("获取User失败 UserId:%d UserToken:%d", video.Author.Id, &params.Token)
			continue
		}
		if isLogin {
			isFollow = rpc.IsFollow(context.Background(), &pb.IsFollowRequest{UserId: video.Id, FollowUserId: int64(token.ID)})
		}
		user := User{
			Id:            info.Id,
			Name:          info.Name,
			FollowerCount: info.FollowerCount,
			FollowCount:   info.FollowCount,
			IsFollow:      isFollow,
		}
		v := Video{
			id:            video.Id,
			user:          user,
			coverUrl:      video.CoverUrl,
			playUrl:       video.PlayUrl,
			favoriteCount: video.FavoriteCount,
			commentCount:  video.CommentCount,
			isFavorite:    video.IsFavorite,
			title:         video.Title,
		}
		videoList = append(videoList, &v)
	}
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: errno.SuccessCode, StatusMsg: errno.Success.ErrMsg},
		VideoList: videoList,
		NextTime:  nextTime,
	})
}
