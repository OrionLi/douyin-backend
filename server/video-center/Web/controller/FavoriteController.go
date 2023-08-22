package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
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

// TODO msg定义为常量
func (h *FavoriteController) ActionFav(c *gin.Context) {
	var requestBody FavoriteParam
	err := c.ShouldBindJSON(&requestBody)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusBadRequest, "data": nil, "msg": "参数错误"})
		return
	}
	userId := validateToken(requestBody.Token)
	if userId == -1 {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusForbidden, "data": nil, "msg": "token错误"})
		return
	}
	videoId := util.StringToInt64(requestBody.VideoId)
	if videoId == -1 {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusBadRequest, "data": nil, "msg": "参数错误"})
		return
	}
	if requestBody.ActionType == "1" {
		err = h.ChatService.CreateFav(videoId, userId)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError, "data": nil, "msg": err.Error()})
			return
		}
	} else if requestBody.ActionType == "2" {
		err := h.ChatService.DeleteFav(videoId, userId)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError, "data": nil, "msg": err.Error()})
			return
		}
	} else {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusBadRequest, "data": nil, "msg": "参数错误"})
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
