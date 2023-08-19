package handler

import (
	"chat-center/pkg/common"
	"chat-center/pkg/util"
	"chat-center/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ChatHandler struct {
	ChatService service.ChatService
}

func NewDiaryHandler(service service.ChatService) *ChatHandler {
	return &ChatHandler{
		ChatService: service,
	}
}

func (h *ChatHandler) GetMessage(c *gin.Context) {
	// TODO 解析token
	// HACK userId暂时定为固定值1
	var currentId int64 = 1

	// 获取参数
	interActiveId := util.StringToInt64(c.Query("to_user_id"))
	preMsgTime := util.StringToInt64(c.Query("pre_msg_time"))
	// 判断是否无法转为int64
	if interActiveId == -1 || preMsgTime == -1 {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "data": nil, "msg": common.ParamErrorMsg})
		return
	}

	if preMsgTime == 0 {
		messageList, err := h.ChatService.GetAllHistoryMessage(currentId, interActiveId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "data": nil, "msg": err.Error()})
		}
		c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": messageList, "msg": common.SuccessMsg})
	} else {
		messageList, err := h.ChatService.GetMessageByPreMsgTime(currentId, interActiveId, preMsgTime)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "data": nil, "msg": err.Error()})
		}
		c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": messageList, "msg": common.SuccessMsg})
	}
}

func (h *ChatHandler) SendMessage(c *gin.Context) {
	// TODO 解析token
	// HACK userId暂时定为固定值1
	var currentId int64 = 1

	// 判断action_type是否为1，不为1返回不支持的action_type
	action := c.Query("action_type")
	if action != "1" {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "data": nil, "msg": common.ActionTypeErrorMsg})
		return
	}

	interActiveId := util.StringToInt64(c.Query("to_user_id"))
	content := c.Query("content")
	// 判断是否无法转为int64
	if interActiveId == -1 || content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "data": nil, "msg": common.ParamErrorMsg})
		return
	}

	err := h.ChatService.SendMessage(currentId, interActiveId, content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "data": nil, "msg": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": nil, "msg": common.SuccessMsg})
}
