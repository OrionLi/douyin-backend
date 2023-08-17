package service

import "douyin-backend/chat-center/generated/message"

type ChatService interface {
	GetAllHistoryMessage(toUserId string) (messageList []message.Message, err error)
	GetMessage(toUserId string, preMsgTime string) (messageList []message.Message, err error)
	SendMessage(toUserId string, fromUserId string, content string) (err error)
}

type ChatServiceImpl struct{}

// GetAllHistoryMessage 根据toUserId查询数据库中所有聊天记录
// 当preMsgTime为0时，查询所有聊天记录
func (c ChatServiceImpl) GetAllHistoryMessage(toUserId string) (messageList []message.Message, err error) {
	//TODO implement me
	panic("implement me")
}

// GetMessage 根据toUserId查询数据库中所有聊天记录
func (c ChatServiceImpl) GetMessage(toUserId string, preMsgTime string) (messageList []message.Message, err error) {
	//TODO implement me
	// 1. 根据toUserId查询数据库中所有聊天记录
	// 2. 将聊天记录转换为message.Message
	// 3. 返回messageList
	panic("implement me")
}

// SendMessage 将消息存入数据库
func (c ChatServiceImpl) SendMessage(toUserId string, fromUserId string, content string) (err error) {
	//TODO implement me
	panic("implement me")
}

func NewChatService() ChatService {
	return &ChatServiceImpl{}
}
