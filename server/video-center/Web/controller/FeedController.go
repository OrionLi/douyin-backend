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
)

func Feed(c *gin.Context) {
	fmt.Println("请求Feed")
	var params FeedParam
	if err := c.ShouldBindQuery(&params); err != nil {
		convertErr := errno.ConvertErr(err)
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{StatusCode: int32(convertErr.ErrCode), StatusMsg: convertErr.ErrMsg},
		})
		return
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
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{StatusCode: int32(convertErr.ErrCode), StatusMsg: convertErr.ErrMsg},
		})
		return
	}
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: errno.SuccessCode, StatusMsg: errno.Success.ErrMsg},
		VideoList: videos,
		NextTime:  nextTime,
	})
}
