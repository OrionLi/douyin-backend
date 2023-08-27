package controller

import (
	"context"
	"gateway-center/grpcClient"
	"gateway-center/pkg/e"
	"gateway-center/pkg/errno"
	baseResponse "gateway-center/response"
	"gateway-center/util"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// FavoriteParam 点赞请求参数
type FavoriteParam struct {
	Token      string `json:"token"`
	VideoId    string `json:"video_id"`
	ActionType string `json:"action_type"`
}

// CommentAction 评论操作
func CommentAction(c *gin.Context) {
	userIdAny, _ := c.Get("UserId")
	userId := userIdAny.(int64)
	var param baseResponse.CommentActionParam
	if err := c.ShouldBind(&param); err != err {
		convertErr := errno.ConvertErr(err)
		c.JSON(http.StatusOK, baseResponse.CommentActionResponse{
			VBResponse: baseResponse.VBResponse{StatusCode: int32(convertErr.ErrCode), StatusMsg: convertErr.ErrMsg},
			Comment:    &pb.Comment{},
		})
		return
	}
	if param.ActionType == "" || param.VideoID == "" {
		c.JSON(http.StatusOK, baseResponse.CommentActionResponse{
			VBResponse: baseResponse.VBResponse{StatusCode: errno.ParamErrCode, StatusMsg: errno.ParamErr.ErrMsg},
			Comment:    &pb.Comment{},
		})
		return
	}
	//判断是哪种操作？
	if param.ActionType == "1" {
		videoId, err := strconv.ParseInt(param.VideoID, 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, baseResponse.CommentActionResponse{
				VBResponse: baseResponse.VBResponse{StatusCode: errno.ParamErrCode, StatusMsg: errno.ParamErr.ErrMsg},
				Comment:    &pb.Comment{},
			})
			return
		}
		request := pb.DouyinCommentActionRequest{
			SelfUserId:  userId,
			VideoId:     videoId,
			ActionType:  0, //保存
			CommentText: param.CommentText,
		}
		response, err := grpcClient.ActionComment(c, &request)
		if err != nil {
			c.JSON(http.StatusOK, baseResponse.CommentActionResponse{
				VBResponse: baseResponse.VBResponse{StatusCode: errno.FailedToCallRpcCode, StatusMsg: errno.FailedToCallRpcErr.ErrMsg},
				Comment:    &pb.Comment{},
			})
		}
		c.JSON(http.StatusOK, response)
		return
	}
	if param.ActionType == "2" {
		videoId, err := strconv.ParseInt(param.VideoID, 10, 64)
		commentID, err := strconv.ParseInt(param.CommentID, 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, baseResponse.CommentActionResponse{
				VBResponse: baseResponse.VBResponse{StatusCode: errno.ParamErrCode, StatusMsg: errno.ParamErr.ErrMsg},
				Comment:    &pb.Comment{},
			})
			return
		}
		request := pb.DouyinCommentActionRequest{
			SelfUserId: userId,
			VideoId:    videoId,
			ActionType: 1, //删除
			CommentId:  commentID,
		}
		response, err2 := grpcClient.ActionComment(c, &request)
		if err2 != nil {
			c.JSON(http.StatusOK, baseResponse.CommentActionResponse{
				VBResponse: baseResponse.VBResponse{StatusCode: errno.FailedToCallRpcCode, StatusMsg: errno.FailedToCallRpcErr.ErrMsg},
				Comment:    &pb.Comment{},
			})
			return
		}
		c.JSON(http.StatusOK, response)
		return
	}
	c.JSON(http.StatusOK, baseResponse.CommentActionResponse{
		VBResponse: baseResponse.VBResponse{StatusCode: errno.ParamErrCode, StatusMsg: errno.ParamErr.ErrMsg},
		Comment:    &pb.Comment{},
	})
}

// CommentList 评论列表
func CommentList(c *gin.Context) {
	userIdAny, _ := c.Get("UserId")
	userId := userIdAny.(int64)
	videoId := c.Query("video_id")
	if videoId == "" {
		c.JSON(http.StatusOK, baseResponse.CommentListResponse{
			VBResponse: baseResponse.VBResponse{StatusCode: errno.ParamErrCode, StatusMsg: errno.ParamErr.ErrMsg},
			Comment:    []*pb.Comment{},
		})
		return
	}
	videoID, err1 := strconv.ParseInt(videoId, 10, 64)
	if err1 != nil {
		c.JSON(http.StatusOK, baseResponse.CommentListResponse{
			VBResponse: baseResponse.VBResponse{StatusCode: errno.ParamErrCode, StatusMsg: errno.ParamErr.ErrMsg},
			Comment:    []*pb.Comment{},
		})
		return
	}

	request := pb.DouyinCommentListRequest{
		SelfUserId: userId,
		VideoId:    videoID,
	}
	response, _ := grpcClient.ListComment(c, &request)
	c.JSON(http.StatusOK, response)
}

// ActionFav 点赞操作
func ActionFav(c *gin.Context) {
	userIdAny, _ := c.Get("UserId")
	userId := userIdAny.(int64)
	videoId := util.StringToInt64(c.Query("video_id"))
	if videoId == -1 {
		c.JSON(http.StatusOK, baseResponse.DouyinFavoriteActionResponse{StatusCode: e.InvalidParams, StatusMsg: e.GetMsg(e.InvalidParams)})
		return
	}
	actionType := util.StringToInt64(c.Query("action_type"))
	resp, err := grpcClient.ActionFavorite(context.Background(), userId, videoId, int32(actionType))
	if err != nil || resp.StatusCode != e.Success {
		c.JSON(http.StatusOK, baseResponse.DouyinFavoriteActionResponse{StatusCode: e.Error, StatusMsg: e.GetMsg(e.Error)})
		return
	}
	c.JSON(http.StatusOK, baseResponse.DouyinFavoriteActionResponse{StatusCode: e.Success, StatusMsg: e.GetMsg(e.Success)})
}

// ListFav 获取喜欢列表
func ListFav(c *gin.Context) {
	userId := c.Query("user_id")
	userIdAny, _ := c.Get("UserId")
	userIdToken := userIdAny.(int64)
	if userId == "" {
		c.JSON(http.StatusOK, baseResponse.FavListResponse{
			VBResponse: baseResponse.VBResponse{StatusCode: errno.ParamErrCode, StatusMsg: errno.ParamErr.ErrMsg},
			FavList:    []*pb.Video{},
		})
		return
	}

	UserIdParseInt, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, baseResponse.FavListResponse{
			VBResponse: baseResponse.VBResponse{StatusCode: errno.ParamErrCode, StatusMsg: errno.ParamErr.ErrMsg},
			FavList:    []*pb.Video{},
		})
		return
	}
	if userIdToken != UserIdParseInt {
		c.JSON(http.StatusOK, baseResponse.FavListResponse{
			VBResponse: baseResponse.VBResponse{StatusCode: errno.ParamErrCode, StatusMsg: errno.ParamErr.ErrMsg},
			FavList:    []*pb.Video{},
		})
		return
	}
	request := pb.DouyinFavoriteListRequest{UserId: userIdToken}
	response, _ := grpcClient.GetFavoriteList(c, &request)
	if response == nil {
		c.JSON(http.StatusOK, &pb.DouyinFavoriteListResponse{
			StatusCode: errno.FavListEmptyCode,
			StatusMsg:  errno.FavListEmptyErr.ErrMsg,
			VideoList:  []*pb.Video{},
		})
		return
	}
	println(response)
	c.JSON(http.StatusOK, response)
}
