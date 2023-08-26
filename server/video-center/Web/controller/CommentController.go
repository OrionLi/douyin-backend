package controller

import (
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"video-center/Web/pkg/baseResponse"
	"video-center/Web/rpc"
	"video-center/pkg/errno"
	"video-center/pkg/util"
)

func CommentAction(c *gin.Context) {
	var param baseResponse.CommentActionParam
	if err := c.ShouldBind(&param); err != err {
		convertErr := errno.ConvertErr(err)
		c.JSON(http.StatusOK, baseResponse.CommentActionResponse{
			Response: baseResponse.Response{StatusCode: int32(convertErr.ErrCode), StatusMsg: convertErr.ErrMsg},
			Comment:  &pb.Comment{},
		})
		return
	}
	if len(param.Token) == 0 || param.ActionType == "" || param.VideoID == "" {
		c.JSON(http.StatusOK, baseResponse.CommentActionResponse{
			Response: baseResponse.Response{StatusCode: errno.ParamErrCode, StatusMsg: errno.ParamErr.ErrMsg},
			Comment:  &pb.Comment{},
		})
		return
	}
	user, err := util.ParseToken(param.Token)
	if err != nil {
		c.JSON(http.StatusOK, baseResponse.CommentActionResponse{
			Response: baseResponse.Response{StatusCode: errno.ParamErrCode, StatusMsg: errno.ParamErr.ErrMsg},
			Comment:  &pb.Comment{},
		})
		return
	}
	//判断是哪种操作？
	if param.ActionType == "1" {
		videoId, err := strconv.ParseInt(param.VideoID, 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, baseResponse.CommentActionResponse{
				Response: baseResponse.Response{StatusCode: errno.ParamErrCode, StatusMsg: errno.ParamErr.ErrMsg},
				Comment:  &pb.Comment{},
			})
			return
		}
		request := pb.DouyinCommentActionRequest{
			SelfUserId:  int64(user.ID),
			VideoId:     videoId,
			ActionType:  0, //保存
			CommentText: param.CommentText,
		}
		response, err := rpc.ActionComment(c, &request)
		if err != nil {
			c.JSON(http.StatusOK, baseResponse.CommentActionResponse{
				Response: baseResponse.Response{StatusCode: errno.FailedToCallRpcCode, StatusMsg: errno.FailedToCallRpcErr.ErrMsg},
				Comment:  &pb.Comment{},
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
				Response: baseResponse.Response{StatusCode: errno.ParamErrCode, StatusMsg: errno.ParamErr.ErrMsg},
				Comment:  &pb.Comment{},
			})
			return
		}
		request := pb.DouyinCommentActionRequest{
			SelfUserId: int64(user.ID),
			VideoId:    videoId,
			ActionType: 1, //删除
			CommentId:  commentID,
		}
		response, err2 := rpc.ActionComment(c, &request)
		if err2 != nil {
			c.JSON(http.StatusOK, baseResponse.CommentActionResponse{
				Response: baseResponse.Response{StatusCode: errno.FailedToCallRpcCode, StatusMsg: errno.FailedToCallRpcErr.ErrMsg},
				Comment:  &pb.Comment{},
			})
			return
		}
		c.JSON(http.StatusOK, response)
		return
	}
	c.JSON(http.StatusOK, baseResponse.CommentActionResponse{
		Response: baseResponse.Response{StatusCode: errno.ParamErrCode, StatusMsg: errno.ParamErr.ErrMsg},
		Comment:  &pb.Comment{},
	})
}

func CommentList(c *gin.Context) {
	token := c.Query("token")
	videoId := c.Query("video_id")
	if len(token) == 0 || videoId == "" {
		c.JSON(http.StatusOK, baseResponse.CommentListResponse{
			Response: baseResponse.Response{StatusCode: errno.ParamErrCode, StatusMsg: errno.ParamErr.ErrMsg},
			Comment:  []*pb.Comment{},
		})
		return
	}
	videoID, err1 := strconv.ParseInt(videoId, 10, 64)
	if err1 != nil {
		c.JSON(http.StatusOK, baseResponse.CommentListResponse{
			Response: baseResponse.Response{StatusCode: errno.ParamErrCode, StatusMsg: errno.ParamErr.ErrMsg},
			Comment:  []*pb.Comment{},
		})
		return
	}
	parseToken, err1 := util.ParseToken(token)
	if err1 != nil {
		c.JSON(http.StatusOK, baseResponse.CommentListResponse{
			Response: baseResponse.Response{StatusCode: errno.ParamErrCode, StatusMsg: errno.ParamErr.ErrMsg},
			Comment:  []*pb.Comment{},
		})
		return
	}

	request := pb.DouyinCommentListRequest{
		SelfUserId: int64(parseToken.ID),
		VideoId:    videoID,
	}
	response, _ := rpc.ListComment(c, &request)
	c.JSON(http.StatusOK, response)
}
