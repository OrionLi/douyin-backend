package service

import (
	"chat-center/dao"
	"chat-center/model"
	"log"
	"strconv"
	"time"
)

type ChatService interface {
	GetAllHistoryMessage(toUserId string) ([]model.Message, error)
	GetMessage(toUserId string, preMsgTime int64) ([]model.Message, error)
	SendMessage(fromUserId string, content string) error
}

type ChatServiceImpl struct{}

// TODO: 该源文件下所有函数实现鉴权

// GetAllHistoryMessage 根据toUserId查询数据库中所有聊天记录
// 当preMsgTime为0时，查询所有聊天记录
func (c ChatServiceImpl) GetAllHistoryMessage(toUserId string) ([]model.Message, error) {
	// 查询相关记录
	parseInt, err := strconv.ParseInt(toUserId, 10, 64)
	if err != nil {
		log.Fatalf("Parse error: %v", err)
		return nil, err
	}
	// FIXME 获取登录用户ID
	fromUserId := 1
	messageList, err := dao.GetAllMessagesByToUserId(int(parseInt), fromUserId)
	if err != nil {
		return nil, err
	}
	return messageList, nil
}

// GetMessage 根据toUserId查询数据库中所有聊天记录
func (c ChatServiceImpl) GetMessage(toUserId string, preMsgTime int64) ([]model.Message, error) {
	timeObj := time.Unix(preMsgTime, 0)
	// 查询相关记录
	parseInt, err := strconv.ParseInt(toUserId, 10, 64)
	if err != nil {
		log.Fatalf("Parse error: %v", err)
		return nil, err
	}
	// FIXME 获取登录用户ID
	fromUserId := 1
	messageList, err := dao.GetMessageByToUserId(timeObj, int(parseInt), fromUserId)
	if err != nil {
		return nil, err
	}
	return messageList, nil
}

// SendMessage 将消息存入数据库
func (c ChatServiceImpl) SendMessage(fromUserId string, content string) error {
	// FIXME 获取登录用户ID
	toUserId := 1
	fromUserIdInt, err := strconv.ParseInt(fromUserId, 10, 64)
	if err != nil {
		log.Fatalf("Parse error: %v", err)
		return err
	}
	message := model.Message{
		ToUserId:   int64(toUserId),
		FromUserId: fromUserIdInt,
		Content:    content,
	}
	err = dao.SendMessage(message)
	if err != nil {
		return err
	}
	return nil
}

func NewChatService() ChatService {
	return &ChatServiceImpl{}
}
