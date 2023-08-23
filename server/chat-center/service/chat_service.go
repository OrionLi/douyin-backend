package service

import (
	"chat-center/dao"
	"chat-center/pkg/utils"
	"fmt"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"github.com/sony/sonyflake"
	"time"
)

type ChatService interface {
	GetAllHistoryMessage(currentId int64, interActiveId int64) ([]*pb.Message, error)
	GetMessageByPreMsgTime(currentId int64, interActiveId int64, preMsgTime int64) ([]*pb.Message, error)
	SendMessage(currentId int64, interActiveId int64, content string) error
}

type ChatServiceImpl struct{}

// GetAllHistoryMessage 根据toUserId查询数据库中所有聊天记录
// 当preMsgTime为0时，查询所有聊天记录
func (c ChatServiceImpl) GetAllHistoryMessage(currentId int64, interActiveId int64) ([]*pb.Message, error) {
	// 查询相关记录
	messageList, err := dao.GetAllMessagesByToUserId(currentId, interActiveId)
	if err != nil {
		utils.LogrusObj.Error("<GetAllHistoryMessage> Get message error: ", err)
		return nil, err
	}
	return messageList, nil
}

// GetMessageByPreMsgTime 根据toUserId查询数据库中所有聊天记录
func (c ChatServiceImpl) GetMessageByPreMsgTime(currentId int64, interActiveId int64, preMsgTime int64) ([]*pb.Message, error) {
	timeObj := time.Unix(preMsgTime, 0)
	// 查询相关记录, 从preMsgTime开始
	// Dao层中toUserId和fromUserId的顺序是反的，因为前端传参中toUserId为对方的ID，fromUserId为自己的ID
	messageList, err := dao.GetMessageByToUserId(timeObj, currentId, interActiveId)
	if err != nil {
		utils.LogrusObj.Error("<GetMessageByPreMsgTime> Get message error: ", err)
		return nil, err
	}
	return messageList, nil
}

// SendMessage 将消息存入数据库
func (c ChatServiceImpl) SendMessage(currentId int64, interActiveId int64, content string) error {
	// 创建一个 SonyFlake 实例
	sf := sonyflake.NewSonyflake(sonyflake.Settings{})

	// 生成唯一 ID
	id, err := sf.NextID()
	if err != nil {
		fmt.Printf("Error generating ID: %s\n", err)
		return err
	}
	now := time.Now()
	format := now.Format("2006-01-02 15:04:05")
	message := &pb.Message{
		Id:         int64(id),
		ToUserId:   interActiveId,
		FromUserId: currentId,
		Content:    content,
		CreateTime: format,
	}
	err = dao.SendMessage(message)
	if err != nil {
		utils.LogrusObj.Error("<SendMessage> Send message error: ", err, " [be from]:", message)
		return err
	}
	return nil
}

func NewChatService() ChatService {
	return &ChatServiceImpl{}
}
