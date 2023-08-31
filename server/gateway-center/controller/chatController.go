package controller

import (
	"gateway-center/grpcClient"
	"gateway-center/pkg/e"
	"gateway-center/response"
	"gateway-center/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetMessage(c *gin.Context) {
	currentIdAny, _ := c.Get("UserId")
	currentId := int64(currentIdAny.(uint))
	interActiveId := util.StringToInt64(c.Query("to_user_id"))
	preMsgTime := util.StringToInt64(c.Query("pre_msg_time"))
	// 判断是否无法转为int64
	if interActiveId == -1 || preMsgTime == -1 {
		c.JSON(http.StatusOK, response.DouyinChatGetMessageResponse{StatusCode: e.InvalidParams, StatusMsg: e.GetMsg(e.InvalidParams)})
		return
	}

	resp, err := grpcClient.GetMessage(currentId, interActiveId, preMsgTime)
	if err != nil || resp.StatusCode != e.Success {
		c.JSON(http.StatusOK, response.DouyinChatGetMessageResponse{StatusCode: e.Error, StatusMsg: e.GetMsg(e.Error)})
		return
	}
	messageList := util.PbMessageListToMessageList(resp.MessageList)
	c.JSON(http.StatusOK, response.DouyinChatGetMessageResponse{
		StatusCode:  e.Success,
		StatusMsg:   e.GetMsg(e.Success),
		MessageList: messageList,
	})
}

func SendMessage(c *gin.Context) {
	currentIdAny, _ := c.Get("UserId")
	currentId := int64(currentIdAny.(uint))
	interActiveIdStr := c.Query("to_user_id")
	actionType := c.Query("action_type")
	content := c.Query("content")
	// 判断action_type是否为1，不为1返回不支持的action_type
	if actionType != "1" {
		c.JSON(http.StatusOK, response.DouyinChatSendMessageResponse{StatusCode: e.InvalidParams, StatusMsg: e.GetMsg(e.InvalidParams)})
		return
	}

	interActiveId := util.StringToInt64(interActiveIdStr)
	// 判断是否无法转为int64
	if interActiveId == -1 {
		c.JSON(http.StatusOK, response.DouyinChatSendMessageResponse{StatusCode: e.InvalidParams, StatusMsg: e.GetMsg(e.InvalidParams)})
		return
	}
	if content == "" {
		c.JSON(http.StatusOK, response.DouyinChatSendMessageResponse{StatusCode: e.InvalidParams, StatusMsg: e.GetMsg(e.InvalidParams)})
		return
	}

	resp, err := grpcClient.SendMessage(currentId, interActiveId, content)
	if err != nil || resp.StatusCode != e.Success {
		c.JSON(http.StatusOK, response.DouyinChatSendMessageResponse{StatusCode: e.Error, StatusMsg: e.GetMsg(e.Error)})
		return
	}
	c.JSON(http.StatusOK, response.DouyinChatSendMessageResponse{StatusCode: e.Success, StatusMsg: e.GetMsg(e.Success)})
}
