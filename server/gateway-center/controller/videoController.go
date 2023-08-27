package controller

import (
	"context"
	"gateway-center/grpcClient"
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
	var param baseResponse.CommentActionParam
	if err := c.ShouldBind(&param); err != err {
		convertErr := errno.ConvertErr(err)
		c.JSON(http.StatusOK, baseResponse.CommentActionResponse{
			VBResponse: baseResponse.VBResponse{StatusCode: int32(convertErr.ErrCode), StatusMsg: convertErr.ErrMsg},
			Comment:    &pb.Comment{},
		})
		return
	}
	if len(param.Token) == 0 || param.ActionType == "" || param.VideoID == "" {
		c.JSON(http.StatusOK, baseResponse.CommentActionResponse{
			VBResponse: baseResponse.VBResponse{StatusCode: errno.ParamErrCode, StatusMsg: errno.ParamErr.ErrMsg},
			Comment:    &pb.Comment{},
		})
		return
	}
	user, err := util.ParseToken(param.Token)
	if err != nil {
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
			SelfUserId:  int64(user.ID),
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
			SelfUserId: int64(user.ID),
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
	token := c.Query("token")
	videoId := c.Query("video_id")
	if len(token) == 0 || videoId == "" {
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
	parseToken, err1 := util.ParseToken(token)
	if err1 != nil {
		c.JSON(http.StatusOK, baseResponse.CommentListResponse{
			VBResponse: baseResponse.VBResponse{StatusCode: errno.ParamErrCode, StatusMsg: errno.ParamErr.ErrMsg},
			Comment:    []*pb.Comment{},
		})
		return
	}

	request := pb.DouyinCommentListRequest{
		SelfUserId: int64(parseToken.ID),
		VideoId:    videoID,
	}
	response, _ := grpcClient.ListComment(c, &request)
	c.JSON(http.StatusOK, response)
}

// ActionFav 点赞操作
func ActionFav(c *gin.Context) {
	var requestBody FavoriteParam
	err := c.ShouldBindJSON(&requestBody)
	if err != nil {
		c.JSON(http.StatusOK, baseResponse.FavActionResponse{baseResponse.VBResponse{StatusCode: errno.ParamErrCode, StatusMsg: errno.ParamErr.ErrMsg}})
		return
	}
	userId := validateToken(requestBody.Token)
	if userId == -1 {
		c.JSON(http.StatusOK, baseResponse.FavActionResponse{baseResponse.VBResponse{StatusCode: errno.TokenErrCode, StatusMsg: errno.TokenErr.ErrMsg}})
		return
	}
	videoId := util.StringToInt64(requestBody.VideoId)
	if videoId == -1 {
		c.JSON(http.StatusOK, baseResponse.FavActionResponse{baseResponse.VBResponse{StatusCode: errno.ParamErrCode, StatusMsg: errno.ParamErr.ErrMsg}})
		return
	}
	actionType := util.StringToInt64(requestBody.ActionType)
	resp, err := grpcClient.ActionFavorite(context.Background(), userId, videoId, int32(actionType))
	if err != nil || resp.StatusCode != errno.SuccessCode {
		c.JSON(http.StatusOK, baseResponse.FavActionResponse{baseResponse.VBResponse{StatusCode: errno.FavActionErrCode, StatusMsg: errno.FavActionErr.ErrMsg}})
		return
	}
	c.JSON(http.StatusOK, baseResponse.FavActionResponse{baseResponse.VBResponse{StatusCode: errno.SuccessCode, StatusMsg: errno.Success.ErrMsg}})
}

// ListFav 获取喜欢列表
func ListFav(c *gin.Context) {
	userId := c.Query("user_id")
	token := c.Query("token")
	if userId == "" || token == "" {
		c.JSON(http.StatusOK, baseResponse.FavListResponse{
			VBResponse: baseResponse.VBResponse{StatusCode: errno.ParamErrCode, StatusMsg: errno.ParamErr.ErrMsg},
			FavList:    []*pb.Video{},
		})
		return
	}
	tokenUserId := validateToken(token)
	UserIdParseInt, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, baseResponse.FavListResponse{
			VBResponse: baseResponse.VBResponse{StatusCode: errno.ParamErrCode, StatusMsg: errno.ParamErr.ErrMsg},
			FavList:    []*pb.Video{},
		})
		return
	}
	if tokenUserId != UserIdParseInt {
		c.JSON(http.StatusOK, baseResponse.FavListResponse{
			VBResponse: baseResponse.VBResponse{StatusCode: errno.ParamErrCode, StatusMsg: errno.ParamErr.ErrMsg},
			FavList:    []*pb.Video{},
		})
		return
	}
	request := pb.DouyinFavoriteListRequest{UserId: tokenUserId}
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
