package controller

import (
	"context"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"github.com/gin-gonic/gin"
	"net/http"
	"video-center/Web/rpc"
	"video-center/pkg/errno"
)

func PublishList(c *gin.Context) {
	var params PublishListParam
	if err := c.ShouldBindJSON(&params); err != nil {
		convertErr := errno.ConvertErr(err)
		c.JSON(http.StatusOK, PublishListResponse{
			Response: Response{StatusCode: int32(convertErr.ErrCode), StatusMsg: convertErr.ErrMsg},
		})
		return
	}
	if params.UserId <= 0 || len(params.Token) == 0 {
		c.JSON(http.StatusOK, PublishListResponse{
			Response: Response{StatusCode: errno.ParamErrCode, StatusMsg: errno.ParamErr.ErrMsg},
		})
		return
	}
	videos, err := rpc.PublishList(context.Background(), &pb.DouyinPublishListRequest{
		UserId: params.UserId,
		Token:  params.Token,
	})
	if err != nil {
		convertErr := errno.ConvertErr(err)
		c.JSON(http.StatusOK, PublishListResponse{
			Response: Response{StatusCode: int32(convertErr.ErrCode), StatusMsg: convertErr.ErrMsg},
		})
		return
	}
	c.JSON(http.StatusOK, PublishListResponse{
		Response:  Response{StatusCode: errno.SuccessCode, StatusMsg: errno.Success.ErrMsg},
		VideoList: videos,
	})
}
