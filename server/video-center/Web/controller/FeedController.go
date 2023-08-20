package controller

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"video-center/Web/rpc"
	"video-center/pkg/errno"
	"video-center/pkg/pb"
)

func Feed(c *gin.Context) {
	var params FeedParam
	if err := c.ShouldBindJSON(&params); err != nil {
		convertErr := errno.ConvertErr(err)
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{StatusCode: int32(convertErr.ErrCode), StatusMsg: convertErr.ErrMsg},
		})
		return
	}
	fmt.Printf("LatestTime: %d, Token: %s\n", params.LatestTime, params.Token)
	if params.LatestTime <= 0 || len(params.Token) == 0 {
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{StatusCode: errno.ParamErrCode, StatusMsg: errno.ParamErr.ErrMsg},
		})
		return
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
