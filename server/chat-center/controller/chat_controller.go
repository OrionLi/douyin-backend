package controller

import (
	"chat-center/model"
	"chat-center/pkg/common"
	"chat-center/pkg/util"
	"chat-center/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type ChatController struct {
	ChatService service.ChatService
}

func NewChatHandler(service service.ChatService) *ChatController {
	return &ChatController{
		ChatService: service,
	}
}

func (h *ChatController) GetMessage(c *gin.Context) {
	currentId := validateToken(c.Query("token"))
	if currentId == -1 {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusForbidden, "data": nil, "msg": common.ForbiddenMsg})
		return
	}

	// 获取参数
	interActiveId := util.StringToInt64(c.Query("to_user_id"))
	preMsgTime := util.StringToInt64(c.Query("pre_msg_time"))
	// 判断是否无法转为int64
	if interActiveId == -1 || preMsgTime == -1 {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusBadRequest, "data": nil, "msg": common.ParamErrorMsg})
		return
	}

	if preMsgTime == 0 {
		messageList, err := h.ChatService.GetAllHistoryMessage(currentId, interActiveId)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError, "data": nil, "msg": err.Error()})
		}
		c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": messageList, "msg": common.SuccessMsg})
	} else {
		messageList, err := h.ChatService.GetMessageByPreMsgTime(currentId, interActiveId, preMsgTime)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError, "data": nil, "msg": err.Error()})
		}
		c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": messageList, "msg": common.SuccessMsg})
	}
}

func (h *ChatController) SendMessage(c *gin.Context) {
	var requestBody model.ActionRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusBadRequest, "data": nil, "msg": common.ParamErrorMsg})
		return
	}

	currentId := validateToken(requestBody.Token)
	if currentId == -1 {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusForbidden, "data": nil, "msg": common.ForbiddenMsg})
		return
	}

	// 判断action_type是否为1，不为1返回不支持的action_type
	if requestBody.ActionType != "1" {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusBadRequest, "data": nil, "msg": common.ActionTypeErrorMsg})
		return
	}

	interActiveId := util.StringToInt64(requestBody.ToUserID)
	// 判断是否无法转为int64
	if interActiveId == -1 {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusBadRequest, "data": nil, "msg": common.ParamErrorMsg})
		return
	}
	if requestBody.Content == "" {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusBadRequest, "data": nil, "msg": common.ContentNullErrorMsg})
		return
	}

	err := h.ChatService.SendMessage(currentId, interActiveId, requestBody.Content)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError, "data": nil, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": nil, "msg": common.SuccessMsg})
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
