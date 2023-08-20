package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"video-center/pkg/errno"
	"video-center/service"
)

func CommentAction(c *gin.Context) {
	var cap CommentActionParam
	if err := c.ShouldBind(&cap); err != err {
		convertErr := errno.ConvertErr(err)
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{StatusCode: int32(convertErr.ErrCode), StatusMsg: convertErr.ErrMsg},
		})
		return
	}
	if len(cap.Token) == 0 || cap.ActionType == "" || cap.VideoID == "" {
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{StatusCode: errno.ParamErrCode, StatusMsg: errno.ParamErr.ErrMsg},
		})
		return
	}
	//todo 用户token验证
	var userId int64

	//判断是哪种操作？
	if cap.ActionType == "1" {
		videoId, err := strconv.ParseInt(cap.VideoID, 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{StatusCode: errno.ParamErrCode, StatusMsg: errno.ParamErr.ErrMsg},
			})
			return
		}
		b, comment, err := service.NewCommentService(c).SaveComment(userId, videoId, cap.CommentText)
		if err != nil {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{StatusCode: errno.CommentPostingCode, StatusMsg: errno.CommentPostingErr.ErrMsg},
			})
			return
		}
		if b {
			// 评论成功
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{errno.SuccessCode, "发布评论成功！"},
				Comment:  comment,
			})
			return
		} else {
			// 评论失败
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{StatusCode: errno.ParamErrCode, StatusMsg: errno.ParamErr.ErrMsg},
			})
			return
		}
	}
	if cap.ActionType == "2" {
		videoId, err := strconv.ParseInt(cap.VideoID, 10, 64)
		commentID, err := strconv.ParseInt(cap.CommentID, 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{StatusCode: errno.ParamErrCode, StatusMsg: errno.ParamErr.ErrMsg},
			})
			return
		}
		b, comment, err := service.NewCommentService(c).DeleteComment(userId, videoId, commentID)
		if err != nil {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{StatusCode: errno.DeleteCommentCode, StatusMsg: errno.DeleteCommentErr.ErrMsg},
			})
			return
		}
		if !b {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{StatusCode: errno.NoMyCommentCode, StatusMsg: errno.NoMyCommentErr.ErrMsg},
			})
			return
		}
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{StatusCode: errno.SuccessCode, StatusMsg: "删除评论成功！"},
			Comment:  comment,
		})
		return
	}
	c.JSON(http.StatusOK, CommentActionResponse{
		Response: Response{StatusCode: errno.ParamErrCode, StatusMsg: errno.ParamErr.ErrMsg},
	})
}

func CommentList(c *gin.Context) {
	token := c.Query("token")
	videoId := c.Query("video_id")
	if len(token) == 0 || videoId == "" {
		c.JSON(http.StatusOK, CommentListResponse{
			Response: Response{StatusCode: errno.ParamErrCode, StatusMsg: errno.ParamErr.ErrMsg},
		})
		return
	}
	//todo 用户token验证
	//var userId int64
	//
	videoID, err1 := strconv.ParseInt(videoId, 10, 64)
	if err1 != nil {
		c.JSON(http.StatusOK, CommentListResponse{
			Response: Response{StatusCode: errno.ParamErrCode, StatusMsg: errno.ParamErr.ErrMsg},
		})
		return
	}
	comments, err2 := service.NewCommentService(c).ListComment(videoID)
	if err2 != nil || len(comments) == 0 {
		c.JSON(http.StatusOK, CommentListResponse{
			Response: Response{StatusCode: errno.NoCommentExistsCode, StatusMsg: errno.NoCommentExistsErr.ErrMsg},
		})
		return
	}
	c.JSON(http.StatusOK, CommentListResponse{
		Response: Response{StatusCode: errno.SuccessCode, StatusMsg: "获取评论信息成功"},
		Comment:  comments,
	})

}
