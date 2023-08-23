package controller

import (
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"video-center/pkg/errno"
	"video-center/pkg/util"
	"video-center/service"
)

func CommentAction(c *gin.Context) {
	var param CommentActionParam
	if err := c.ShouldBind(&param); err != err {
		convertErr := errno.ConvertErr(err)
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{StatusCode: int32(convertErr.ErrCode), StatusMsg: convertErr.ErrMsg},
			Comment:  &pb.Comment{},
		})
		return
	}
	if len(param.Token) == 0 || param.ActionType == "" || param.VideoID == "" {
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{StatusCode: errno.ParamErrCode, StatusMsg: errno.ParamErr.ErrMsg},
			Comment:  &pb.Comment{},
		})
		return
	}
	//todo 用户token验证
	user, err := util.ParseToken(param.Token)
	if err != nil {
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{StatusCode: errno.ParamErrCode, StatusMsg: errno.ParamErr.ErrMsg},
			Comment:  &pb.Comment{},
		})
		return
	}
	//判断是哪种操作？
	if param.ActionType == "1" {
		videoId, err := strconv.ParseInt(param.VideoID, 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{StatusCode: errno.ParamErrCode, StatusMsg: errno.ParamErr.ErrMsg},
				Comment:  &pb.Comment{},
			})
			return
		}
		b, comment, err := service.NewCommentService(c).SaveComment(int64(user.ID), videoId, param.CommentText)
		if err != nil {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{StatusCode: errno.CommentPostingCode, StatusMsg: errno.CommentPostingErr.ErrMsg},
				Comment:  &pb.Comment{},
			})
			return
		}
		if b {
			// 评论成功
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{errno.SuccessCode, errno.Success.ErrMsg},
				Comment:  comment,
			})
			return
		} else {
			// 评论失败
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{StatusCode: errno.CommentPostingCode, StatusMsg: errno.CommentPostingErr.ErrMsg},
				Comment:  &pb.Comment{},
			})
			return
		}
	}
	if param.ActionType == "2" {
		videoId, err := strconv.ParseInt(param.VideoID, 10, 64)
		commentID, err := strconv.ParseInt(param.CommentID, 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{StatusCode: errno.ParamErrCode, StatusMsg: errno.ParamErr.ErrMsg},
				Comment:  &pb.Comment{},
			})
			return
		}
		b, comment, err := service.NewCommentService(c).DeleteComment(int64(user.ID), videoId, commentID)
		if err != nil {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{StatusCode: errno.DeleteCommentCode, StatusMsg: errno.DeleteCommentErr.ErrMsg},
				Comment:  &pb.Comment{},
			})
			return
		}
		if !b {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{StatusCode: errno.NoMyCommentCode, StatusMsg: errno.NoMyCommentErr.ErrMsg},
				Comment:  &pb.Comment{},
			})
			return
		}
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{StatusCode: errno.SuccessCode, StatusMsg: errno.Success.ErrMsg},
			Comment:  comment,
		})
		return
	}
	c.JSON(http.StatusOK, CommentActionResponse{
		Response: Response{StatusCode: errno.ParamErrCode, StatusMsg: errno.ParamErr.ErrMsg},
		Comment:  &pb.Comment{},
	})
}

func CommentList(c *gin.Context) {
	token := c.Query("token")
	videoId := c.Query("video_id")
	if len(token) == 0 || videoId == "" {
		c.JSON(http.StatusOK, CommentListResponse{
			Response: Response{StatusCode: errno.ParamErrCode, StatusMsg: errno.ParamErr.ErrMsg},
			Comment:  []*pb.Comment{},
		})
		return
	}
	videoID, err1 := strconv.ParseInt(videoId, 10, 64)
	if err1 != nil {
		c.JSON(http.StatusOK, CommentListResponse{
			Response: Response{StatusCode: errno.ParamErrCode, StatusMsg: errno.ParamErr.ErrMsg},
			Comment:  []*pb.Comment{},
		})
		return
	}
	comments, err2 := service.NewCommentService(c).ListComment(videoID)
	if err2 != nil || len(comments) == 0 {
		c.JSON(http.StatusOK, CommentListResponse{
			Response: Response{StatusCode: errno.NoCommentExistsCode, StatusMsg: errno.NoCommentExistsErr.ErrMsg},
			Comment:  []*pb.Comment{},
		})
		return
	}
	c.JSON(http.StatusOK, CommentListResponse{
		Response: Response{StatusCode: errno.SuccessCode, StatusMsg: errno.Success.ErrMsg},
		Comment:  comments,
	})
}
