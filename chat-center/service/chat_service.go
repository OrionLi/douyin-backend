package service

import (
	"chat-center/dao"
	"chat-center/model"
	"time"
)

type ChatService interface {
	GetAllHistoryMessage(currentId int64, interActiveId int64) ([]model.Message, error)
	GetMessageByPreMsgTime(currentId int64, interActiveId int64, preMsgTime int64) ([]model.Message, error)
	SendMessage(currentId int64, interActiveId int64, content string) error
}

type ChatServiceImpl struct{}

// GetAllHistoryMessage 根据toUserId查询数据库中所有聊天记录
// 当preMsgTime为0时，查询所有聊天记录
func (c ChatServiceImpl) GetAllHistoryMessage(currentId int64, interActiveId int64) ([]model.Message, error) {
	// 查询相关记录
	messageList, err := dao.GetAllMessagesByToUserId(currentId, interActiveId)
	if err != nil {
		return nil, err
	}
	return messageList, nil
}

// GetMessageByPreMsgTime 根据toUserId查询数据库中所有聊天记录
func (c ChatServiceImpl) GetMessageByPreMsgTime(currentId int64, interActiveId int64, preMsgTime int64) ([]model.Message, error) {
	timeObj := time.Unix(preMsgTime, 0)
	// 查询相关记录, 从preMsgTime开始
	// Dao层中toUserId和fromUserId的顺序是反的，因为前端传参中toUserId为对方的ID，fromUserId为自己的ID
	messageList, err := dao.GetMessageByToUserId(timeObj, currentId, interActiveId)
	if err != nil {
		return nil, err
	}
	return messageList, nil
}

// SendMessage 将消息存入数据库
func (c ChatServiceImpl) SendMessage(currentId int64, interActiveId int64, content string) error {
	message := model.Message{
		ToUserId:   interActiveId,
		FromUserId: currentId,
		Content:    content,
	}
	err := dao.SendMessage(message)
	if err != nil {
		return err
	}
	return nil
}

func NewChatService() ChatService {
	return &ChatServiceImpl{}
}
