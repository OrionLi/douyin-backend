package controller

import (
	"context"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"github.com/gin-gonic/gin"
	"net/http"
	"video-center/Web/rpc"
	"video-center/pkg/errno"
	"video-center/pkg/util"
)

func PublishList(c *gin.Context) {
	var params PublishListParam
	if err := c.ShouldBindJSON(&params); err != nil {
		convertErr := errno.ConvertErr(err)
		util.LogrusObj.Errorf("参数绑定错误 URL:%s form %v 错误原因:%s", c.Request.URL, c.Request.PostForm, convertErr.ErrMsg)
		c.JSON(http.StatusOK, PublishListResponse{
			Response: Response{StatusCode: int32(convertErr.ErrCode), StatusMsg: convertErr.ErrMsg},
		})
		return
	}
	if params.UserId <= 0 || len(params.Token) == 0 {
		util.LogrusObj.Errorf("Token格式错误 URL:%s Token:%s UserId:%d", c.Request.RequestURI, params.Token, params.UserId)
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
		util.LogrusObj.Errorf("RPC Error ErrorMSG:%s", convertErr.ErrMsg)
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
