package controller

import (
	"gateway-center/grpcClient"
	"gateway-center/pkg/e"
	"gateway-center/response"
	"gateway-center/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// TODO 重新定义部分返回Code和Msg

type ActionRequest struct {
	Token      string `json:"token"`
	ToUserID   string `json:"to_user_id"`
	ActionType string `json:"action_type"`
	Content    string `json:"content"`
}

func GetMessage(c *gin.Context) {
	currentId := validateToken(c.Query("token"))
	if currentId == -1 {
		c.JSON(http.StatusOK, response.DouyinChatGetMessageResponse{StatusCode: e.ErrorAuthToken, StatusMsg: e.GetMsg(e.ErrorAuthToken)})
		return
	}

	// 获取参数
	interActiveId := util.StringToInt64(c.Query("to_user_id"))
	preMsgTime := util.StringToInt64(c.Query("pre_msg_time"))
	// 判断是否无法转为int64
	if interActiveId == -1 || preMsgTime == -1 {
		c.JSON(http.StatusOK, response.DouyinChatGetMessageResponse{StatusCode: e.InvalidParams, StatusMsg: e.GetMsg(e.InvalidParams)})
		return
	}

	resp, err := grpcClient.GetMessage(currentId, interActiveId, preMsgTime)
	if err != nil || resp.StatusCode != e.Success {
		c.JSON(http.StatusOK, response.DouyinChatGetMessageResponse{StatusCode: e.Error, StatusMsg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.DouyinChatGetMessageResponse{StatusCode: e.Success, StatusMsg: e.GetMsg(e.Success), MessageList: resp.GetMessageList()})
}

func SendMessage(c *gin.Context) {
	var requestBody ActionRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusOK, response.DouyinChatSendMessageResponse{StatusCode: e.InvalidParams, StatusMsg: e.GetMsg(e.InvalidParams)})
		return
	}

	currentId := validateToken(requestBody.Token)
	if currentId == -1 {
		c.JSON(http.StatusOK, response.DouyinChatSendMessageResponse{StatusCode: e.ErrorAuthToken, StatusMsg: e.GetMsg(e.ErrorAuthToken)})
		return
	}

	// 判断action_type是否为1，不为1返回不支持的action_type
	if requestBody.ActionType != "1" {
		c.JSON(http.StatusOK, response.DouyinChatSendMessageResponse{StatusCode: e.InvalidParams, StatusMsg: e.GetMsg(e.InvalidParams)})
		return
	}

	interActiveId := util.StringToInt64(requestBody.ToUserID)
	// 判断是否无法转为int64
	if interActiveId == -1 {
		c.JSON(http.StatusOK, response.DouyinChatSendMessageResponse{StatusCode: e.InvalidParams, StatusMsg: e.GetMsg(e.InvalidParams)})
		return
	}
	if requestBody.Content == "" {
		c.JSON(http.StatusOK, response.DouyinChatSendMessageResponse{StatusCode: e.InvalidParams, StatusMsg: e.GetMsg(e.InvalidParams)})
		return
	}

	resp, err := grpcClient.SendMessage(currentId, interActiveId, requestBody.Content)
	if err != nil || resp.StatusCode != e.Success {
		c.JSON(http.StatusOK, response.DouyinChatSendMessageResponse{StatusCode: e.Error, StatusMsg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.DouyinChatSendMessageResponse{StatusCode: e.Success, StatusMsg: e.GetMsg(e.Success)})
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
