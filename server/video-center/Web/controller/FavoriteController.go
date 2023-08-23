package controller

import (
	"context"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
	"video-center/Web/rpc"
	"video-center/pkg/errno"
	"video-center/pkg/util"
	"video-center/service"
)

type FavoriteController struct {
	ChatService service.FavoriteService
}

type FavoriteParam struct {
	Token      string `json:"token"`
	VideoId    string `json:"video_id"`
	ActionType string `json:"action_type"`
}

func NewFavoriteController(service service.FavoriteService) *FavoriteController {
	return &FavoriteController{
		ChatService: service,
	}
}

func (h *FavoriteController) ActionFav(c *gin.Context) {
	var requestBody FavoriteParam
	err := c.ShouldBindJSON(&requestBody)
	if err != nil {
		c.JSON(http.StatusOK, FavActionResponse{Response{StatusCode: errno.SuccessCode, StatusMsg: errno.ParamErr.ErrMsg}})
		return
	}
	userId := validateToken(requestBody.Token)
	if userId == -1 {
		c.JSON(http.StatusOK, FavActionResponse{Response{StatusCode: errno.TokenErrCode, StatusMsg: errno.TokenErr.ErrMsg}})
		return
	}
	videoId := util.StringToInt64(requestBody.VideoId)
	if videoId == -1 {
		c.JSON(http.StatusOK, FavActionResponse{Response{StatusCode: errno.ParamErrCode, StatusMsg: errno.ParamErr.ErrMsg}})
		return
	}
	actionType := util.StringToInt64(requestBody.ActionType)
	resp, err := rpc.ActionFavorite(context.Background(), userId, videoId, int32(actionType))
	if err != nil || resp.StatusCode != errno.SuccessCode {
		c.JSON(http.StatusOK, FavActionResponse{Response{StatusCode: errno.FavActionErrCode, StatusMsg: errno.FavActionErr.ErrMsg}})
		return
	}
	c.JSON(http.StatusOK, FavActionResponse{Response{StatusCode: errno.SuccessCode, StatusMsg: errno.Success.ErrMsg}})
}

// ListFav 获取喜欢列表
func (h *FavoriteController) ListFav(context *gin.Context) {
	userId := context.Query("user_id")
	token := context.Query("token")
	if userId == "" || token == "" {
		context.JSON(http.StatusOK, FavListResponse{
			Response: Response{StatusCode: errno.ParamErrCode, StatusMsg: errno.ParamErr.ErrMsg},
		})
		return
	}
	tokenUserId := validateToken(token)
	UserIdParseInt, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		context.JSON(http.StatusOK, FavListResponse{
			Response: Response{StatusCode: errno.ParamErrCode, StatusMsg: errno.ParamErr.ErrMsg},
		})
		return
	}
	if tokenUserId != UserIdParseInt {
		context.JSON(http.StatusOK, FavListResponse{
			Response: Response{StatusCode: errno.ParamErrCode, StatusMsg: errno.ParamErr.ErrMsg},
		})
		return
	}
	b, favs := h.ChatService.ListFav(UserIdParseInt)
	if !b {
		context.JSON(http.StatusOK, FavListResponse{
			Response: Response{StatusCode: errno.FavListEmptyCode, StatusMsg: errno.FavListEmptyErr.ErrMsg},
			FavList:  []*pb.Video{},
		})
		return
	}
	if b {
		context.JSON(http.StatusOK, FavListResponse{
			Response: Response{StatusCode: errno.SuccessCode, StatusMsg: errno.Success.ErrMsg},
			FavList:  favs,
		})
		return
	}
}

// validateToken 验证token
func validateToken(token string) int64 {
	parseToken, err := util.ParseToken(token)
	if err != nil {
		return -1
	}
	// 判断 token 是否过期
	if parseToken.ExpiresAt < time.Now().Unix() {
		return -1
	}
	return int64(parseToken.ID)
}
