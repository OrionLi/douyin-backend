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

func NewChatController(service service.ChatService) *ChatController {
	return &ChatController{
		ChatService: service,
	}
}

func (h *ChatController) GetMessage(c *gin.Context) {
	//currentId := validateToken(c.Query("token"))
	var currentId int64 = 1
	if currentId == -1 {
		c.JSON(http.StatusOK, common.GetMessageResponse{Response: common.Response{StatusCode: http.StatusForbidden, StatusMsg: common.ForbiddenMsg}})
		return
	}

	// 获取参数
	interActiveId := util.StringToInt64(c.Query("to_user_id"))
	preMsgTime := util.StringToInt64(c.Query("pre_msg_time"))
	// 判断是否无法转为int64
	if interActiveId == -1 || preMsgTime == -1 {
		c.JSON(http.StatusOK, common.GetMessageResponse{Response: common.Response{StatusCode: http.StatusBadRequest, StatusMsg: common.ParamErrorMsg}})
		return
	}

	if preMsgTime == 0 {
		messageList, err := h.ChatService.GetAllHistoryMessage(currentId, interActiveId)
		if err != nil {
			c.JSON(http.StatusOK, common.GetMessageResponse{Response: common.Response{StatusCode: http.StatusInternalServerError, StatusMsg: err.Error()}})
			return
		}
		c.JSON(http.StatusOK, common.GetMessageResponse{Response: common.Response{StatusCode: common.SuccessCode, StatusMsg: common.SuccessMsg}, MessageList: messageList})
		return
	} else {
		messageList, err := h.ChatService.GetMessageByPreMsgTime(currentId, interActiveId, preMsgTime)
		if err != nil {
			c.JSON(http.StatusOK, common.GetMessageResponse{Response: common.Response{StatusCode: http.StatusInternalServerError, StatusMsg: err.Error()}})
			return
		}
		c.JSON(http.StatusOK, common.GetMessageResponse{Response: common.Response{StatusCode: common.SuccessCode, StatusMsg: common.SuccessMsg}, MessageList: messageList})
	}
}

func (h *ChatController) SendMessage(c *gin.Context) {
	var requestBody model.ActionRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusOK, common.SendMessageResponse{Response: common.Response{StatusCode: http.StatusBadRequest, StatusMsg: common.ParamErrorMsg}})
		return
	}

	//currentId := validateToken(requestBody.Token)
	var currentId int64 = 1
	if currentId == -1 {
		c.JSON(http.StatusOK, common.SendMessageResponse{Response: common.Response{StatusCode: http.StatusForbidden, StatusMsg: common.ForbiddenMsg}})
		return
	}

	// 判断action_type是否为1，不为1返回不支持的action_type
	if requestBody.ActionType != "1" {
		c.JSON(http.StatusOK, common.SendMessageResponse{Response: common.Response{StatusCode: http.StatusBadRequest, StatusMsg: common.ParamErrorMsg}})
		return
	}

	interActiveId := util.StringToInt64(requestBody.ToUserID)
	// 判断是否无法转为int64
	if interActiveId == -1 {
		c.JSON(http.StatusOK, common.SendMessageResponse{Response: common.Response{StatusCode: http.StatusBadRequest, StatusMsg: common.ParamErrorMsg}})
		return
	}
	if requestBody.Content == "" {
		c.JSON(http.StatusOK, common.SendMessageResponse{Response: common.Response{StatusCode: http.StatusBadRequest, StatusMsg: common.ParamErrorMsg}})
		return
	}

	err := h.ChatService.SendMessage(currentId, interActiveId, requestBody.Content)
	if err != nil {
		c.JSON(http.StatusOK, common.SendMessageResponse{Response: common.Response{StatusCode: http.StatusInternalServerError, StatusMsg: err.Error()}})
		return
	}

	c.JSON(http.StatusOK, common.SendMessageResponse{Response: common.Response{StatusCode: common.SuccessCode, StatusMsg: common.SuccessMsg}})
}

// validateToken 验证token
func validateToken(token string) int64 {
	if len(token) == 0 {
		return -1
	}
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

func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		request := c.Request
		util.LogrusObj.Infof("URL:%s host:%s method:%s remoteIp:%s", request.URL, request.Host, request.Method, request.RemoteAddr)
		c.Next()
	}
}
