package util

import (
	"gateway-center/response"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"strconv"
)

func StringToInt64(str string) int64 {
	parseInt, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return -1
	}
	return parseInt
}

func PbMessageListToMessageList(pbMessageList []*pb.Message) []response.Message {
	var messageList []response.Message
	for _, v := range pbMessageList {
		var message response.Message
		message.Id = v.Id
		message.ToUserId = v.ToUserId
		message.FromUserId = v.FromUserId
		message.Content = v.Content
		time, _ := strconv.ParseInt(v.CreateTime, 10, 64)
		message.CreateTime = time
		messageList = append(messageList, message)
	}
	return messageList
}
