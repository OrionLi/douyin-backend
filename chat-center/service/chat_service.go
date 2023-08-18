package service

import (
	"douyin-backend/chat-center/dao"
	"douyin-backend/chat-center/generated/message"
	"log"
	"strconv"
)

type ChatService interface {
	GetAllHistoryMessage(toUserId string) ([]message.Message, error)
	GetMessage(toUserId string, preMsgTime int64) ([]message.Message, error)
	SendMessage(toUserId string, fromUserId string, content string) error
}

type ChatServiceImpl struct{}

// GetAllHistoryMessage 根据toUserId查询数据库中所有聊天记录
// 当preMsgTime为0时，查询所有聊天记录
func (c ChatServiceImpl) GetAllHistoryMessage(toUserId string) ([]message.Message, error) {
	// 查询相关记录
	parseInt, err := strconv.ParseInt(toUserId, 10, 64)
	if err != nil {
		log.Fatalf("Parse error: %v", err)
		return nil, err
	}
	messageList, err := dao.GetAllMessagesByToUserId(int(parseInt))
	if err != nil {
		return nil, err
	}
	return messageList, nil
}

// GetMessage 根据toUserId查询数据库中所有聊天记录
func (c ChatServiceImpl) GetMessage(toUserId string, preMsgTime int64) ([]message.Message, error) {
	//TODO implement me
	panic("implement me")
}

// SendMessage 将消息存入数据库
func (c ChatServiceImpl) SendMessage(toUserId string, fromUserId string, content string) error {
	//TODO implement me
	panic("implement me")
}

func NewChatService() ChatService {
	return &ChatServiceImpl{}
}
