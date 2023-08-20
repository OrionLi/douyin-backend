package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
	"video-center/Web/rpc"
	"video-center/pkg/errno"
	"video-center/pkg/pb"
)

func PublishAction(c *gin.Context) {
	var params PublishActionParam
	if err := c.ShouldBind(&params); err != nil {
		convertErr := errno.ConvertErr(err)
		c.JSON(http.StatusOK, PublishListResponse{
			Response: Response{StatusCode: int32(convertErr.ErrCode), StatusMsg: convertErr.ErrMsg},
		})
		return
	}
	//需要在这里判断Token的合法性，根据Token来获取UserId，然后通过查询此Id下的所有视频，找到当前的视频数据，查询其存储的路径
	if len(params.Token) == 0 {
		c.JSON(http.StatusOK, PublishListResponse{
			Response: Response{StatusCode: errno.ParamErrCode, StatusMsg: errno.ParamErr.ErrMsg},
		})
		return
	}
	//获取视频二进制数据
	file, err2 := c.FormFile("data")
	if err2 != nil {
		convertErr := errno.ConvertErr(err2)
		c.JSON(http.StatusOK, PublishListResponse{
			Response: Response{StatusCode: int32(convertErr.ErrCode), StatusMsg: convertErr.ErrMsg},
		})
		return
	}
	open, err2 := file.Open()
	if err2 != nil {
		convertErr := errno.ConvertErr(err2)
		c.JSON(http.StatusOK, PublishListResponse{
			Response: Response{StatusCode: int32(convertErr.ErrCode), StatusMsg: convertErr.ErrMsg},
		})
		return
	}
	fileData := make([]byte, file.Size)
	_, err2 = open.Read(fileData)
	defer func(open multipart.File) {
		err := open.Close()
		if err != nil {

		}
	}(open)
	if err2 != nil {
		convertErr := errno.ConvertErr(err2)
		c.JSON(http.StatusOK, PublishListResponse{
			Response: Response{StatusCode: int32(convertErr.ErrCode), StatusMsg: convertErr.ErrMsg},
		})
		return
	}
	err := rpc.PublishAction(context.Background(), &pb.DouyinPublishActionRequest{
		Token: params.Token,
		Data:  fileData,
		Title: params.Title,
	})
	if err != nil {
		convertErr := errno.ConvertErr(err)
		c.JSON(http.StatusOK, PublishListResponse{
			Response: Response{StatusCode: int32(convertErr.ErrCode), StatusMsg: convertErr.ErrMsg},
		})
		return
	}
	c.JSON(http.StatusOK, PublishListResponse{
		Response: Response{StatusCode: errno.SuccessCode, StatusMsg: errno.Success.ErrMsg},
	})
}
