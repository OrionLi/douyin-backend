package controller

import (
	"context"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"github.com/gin-gonic/gin"
	"net/http"
	"video-center/Web/pkg/baseResponse"
	"video-center/Web/rpc"
	"video-center/pkg/errno"
	"video-center/pkg/util"
)

func PublishList(c *gin.Context) {
	var params baseResponse.PublishListParam
	if err := c.ShouldBindJSON(&params); err != nil {
		convertErr := errno.ConvertErr(err)
		util.LogrusObj.Errorf("参数绑定错误 URL:%s form %v 错误原因:%s", c.Request.URL, c.Request.PostForm, convertErr.ErrMsg)
		c.JSON(http.StatusOK, baseResponse.PublishListResponse{
			Response: baseResponse.Response{StatusCode: int32(convertErr.ErrCode), StatusMsg: convertErr.ErrMsg},
		})
		return
	}
	if params.UserId <= 0 || len(params.Token) == 0 {
		util.LogrusObj.Errorf("Token格式错误 URL:%s Token:%s UserId:%d", c.Request.RequestURI, params.Token, params.UserId)
		c.JSON(http.StatusOK, baseResponse.PublishListResponse{
			Response: baseResponse.Response{StatusCode: errno.ParamErrCode, StatusMsg: errno.ParamErr.ErrMsg},
		})
		return
	}
	_, err2 := util.ParseToken(params.Token)
	if err2 != nil {
		util.LogrusObj.Errorf("Token验证失败 URL:%s Token:%s UserId:%d", c.Request.RequestURI, params.Token, params.UserId)
		c.JSON(http.StatusOK, baseResponse.PublishListResponse{
			Response: baseResponse.Response{StatusCode: errno.TokenErrCode, StatusMsg: errno.TokenErr.ErrMsg},
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
		c.JSON(http.StatusOK, baseResponse.PublishListResponse{
			Response: baseResponse.Response{StatusCode: int32(convertErr.ErrCode), StatusMsg: convertErr.ErrMsg},
		})
		return
	}
	videoList := make([]*baseResponse.Video, 0)
	for _, video := range videos {
		//todo 封装User信息
		//info, err := rpc.GetUserInfo(context.Background(), &pb.DouyinUserRequest{UserId: video.Author.Id})
		//if err != nil {
		//	util.LogrusObj.Errorf("获取User失败 UserId:%d UserToken:%d", video.Author.Id, &params.Token)
		//	continue
		//}
		//user := User{
		//	Id:            info.Id,
		//	Name:          info.Name,
		//	FollowerCount: info.FollowerCount,
		//	FollowCount:   info.FollowCount,
		//	IsFollow:      false,
		//}
		v := baseResponse.Video{
			Id:            video.Id,
			User:          baseResponse.User{},
			CoverUrl:      video.CoverUrl,
			PlayUrl:       video.PlayUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    video.IsFavorite,
			Title:         video.Title,
		}
		videoList = append(videoList, &v)
	}
	c.JSON(http.StatusOK, baseResponse.PublishListResponse{
		Response:  baseResponse.Response{StatusCode: errno.SuccessCode, StatusMsg: errno.Success.ErrMsg},
		VideoList: videoList,
	})
}
