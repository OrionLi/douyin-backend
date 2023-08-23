package controller

import (
	"chat-center/Web/rpc/client"
	"chat-center/pkg/common"
	"chat-center/pkg/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type ActionRequest struct {
	Token      string `json:"token"`
	ToUserID   string `json:"to_user_id"`
	ActionType string `json:"action_type"`
	Content    string `json:"content"`
}

func GetMessage(c *gin.Context) {
	currentId := validateToken(c.Query("token"))
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

	resp, err := client.GetMessage(currentId, interActiveId, preMsgTime)
	if err != nil || resp.StatusCode != common.SuccessCode {
		c.JSON(http.StatusOK, common.GetMessageResponse{Response: common.Response{StatusCode: http.StatusInternalServerError, StatusMsg: err.Error()}})
		return
	}
	c.JSON(http.StatusOK, common.GetMessageResponse{Response: common.Response{StatusCode: common.SuccessCode, StatusMsg: common.SuccessMsg}, MessageList: resp.GetMessageList()})
}

func SendMessage(c *gin.Context) {
	var requestBody ActionRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusOK, common.SendMessageResponse{Response: common.Response{StatusCode: http.StatusBadRequest, StatusMsg: common.ParamErrorMsg}})
		return
	}

	currentId := validateToken(requestBody.Token)
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

	resp, err := client.SendMessage(currentId, interActiveId, requestBody.Content)
	if err != nil || resp.StatusCode != common.SuccessCode {
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
