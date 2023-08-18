package rpc

import (
	"chat-center/dao"
	"chat-center/generated/message"
	"chat-center/model"
	"chat-center/pkg/common"
	"context"
	"time"
)

type ChatRPCService struct {
	message.UnimplementedDouyinMessageServiceServer
}

func NewChatRPCService() *ChatRPCService {
	return &ChatRPCService{}
}

func (s *ChatRPCService) GetMessage(ctx context.Context, request *message.DouyinMessageChatRequest) (*message.DouyinMessageChatResponse, error) {
	// HACK 如果不鉴权： userId := ctx.Value("userId")
	// HACK 如果鉴权： token := req.GetToken()，解析token，获取userId
	// HACK userId暂时定为固定值1
	fromUserId := 1
	toUserId := request.GetToUserId()
	preMsgTime := request.GetPreMsgTime()
	if preMsgTime == 0 {
		messageListTemp, err := dao.GetAllMessagesByToUserId(int(toUserId), fromUserId)
		messageList := messageListToPbMessageList(messageListTemp)
		if err != nil {
			return &message.DouyinMessageChatResponse{}, err
		}
		return &message.DouyinMessageChatResponse{
			MessageList: messageList,
		}, nil
	} else {
		// 将时间戳转换为 time.Time 类型
		timeObj := time.Unix(preMsgTime, 0)
		messageListTemp, err := dao.GetMessageByToUserId(timeObj, int(toUserId), fromUserId)
		messageList := messageListToPbMessageList(messageListTemp)
		if err != nil {
			return &message.DouyinMessageChatResponse{}, err
		}
		return &message.DouyinMessageChatResponse{
			MessageList: messageList,
		}, nil
	}
}

func (s *ChatRPCService) SendMessage(ctx context.Context, request *message.DouyinMessageActionRequest) (*message.DouyinMessageActionResponse, error) {
	// HACK 如果不鉴权： userId := ctx.Value("userId")
	// HACK 如果鉴权： token := req.GetToken()，解析token，获取userId
	// HACK userId暂时定为固定值1
	fromUserId := 1
	toUserId := request.GetToUserId()
	content := request.GetContent()
	msg := model.Message{
		ToUserId:   int64(int(toUserId)),
		FromUserId: int64(fromUserId),
		Content:    content,
	}
	err := dao.SendMessage(msg)
	if err != nil {
		return &message.DouyinMessageActionResponse{}, err
	}
	return &message.DouyinMessageActionResponse{
		StatusCode: common.SuccessCode,
		StatusMsg:  common.SuccessMsg,
	}, nil
}

func messageListToPbMessageList(msg []model.Message) []*message.Message {
	var messageList []*message.Message
	for _, v := range msg {
		messageList = append(messageList, &message.Message{
			Id:         v.Id,
			FromUserId: v.FromUserId,
			ToUserId:   v.ToUserId,
			Content:    v.Content,
			CreateTime: v.CreateTime,
		})
	}
	return messageList
}
