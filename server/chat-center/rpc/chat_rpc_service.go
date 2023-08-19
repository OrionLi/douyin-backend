package rpc

import (
	"chat-center/dao"
	"chat-center/model"
	"chat-center/pkg/common"
	"context"
	"github.com/OrionLi/douyin-backend/pkg/pb"

	"time"
)

type ChatRPCService struct {
	pb.UnimplementedDouyinMessageServiceServer
}

func NewChatRPCService() *ChatRPCService {
	return &ChatRPCService{}
}

func (s *ChatRPCService) GetMessage(ctx context.Context, request *pb.DouyinMessageChatRequest) (*pb.DouyinMessageChatResponse, error) {
	// HACK 如果不鉴权： userId := ctx.Value("userId")
	// HACK 如果鉴权： token := req.GetToken()，解析token，获取userId
	// HACK userId暂时定为固定值1
	// toUserId为对方id fromUserId为自己id
	var fromUserId int64 = 1
	toUserId := request.GetToUserId()
	preMsgTime := request.GetPreMsgTime()
	if preMsgTime == 0 {
		messageListTemp, err := dao.GetAllMessagesByToUserId(toUserId, fromUserId)
		messageList := messageListToPbMessageList(messageListTemp)
		if err != nil {
			return &pb.DouyinMessageChatResponse{
				StatusCode:  common.ErrorGetCode,
				StatusMsg:   common.ErrorGetMsg,
				MessageList: nil,
			}, err
		}
		return &pb.DouyinMessageChatResponse{
			StatusCode:  common.SuccessCode,
			StatusMsg:   common.SuccessMsg,
			MessageList: messageList,
		}, nil
	} else {
		// 将时间戳转换为 time.Time 类型
		timeObj := time.Unix(preMsgTime, 0)
		messageListTemp, err := dao.GetMessageByToUserId(timeObj, toUserId, fromUserId)
		messageList := messageListToPbMessageList(messageListTemp)
		if err != nil {
			return &pb.DouyinMessageChatResponse{
				StatusCode:  common.ErrorGetCode,
				StatusMsg:   common.ErrorGetMsg,
				MessageList: nil,
			}, err
		}
		return &pb.DouyinMessageChatResponse{
			StatusCode:  common.SuccessCode,
			StatusMsg:   common.SuccessMsg,
			MessageList: messageList,
		}, nil
	}
}

func (s *ChatRPCService) SendMessage(ctx context.Context, request *pb.DouyinMessageActionRequest) (*pb.DouyinMessageActionResponse, error) {
	// HACK 如果不鉴权： userId := ctx.Value("userId")
	// HACK 如果鉴权： token := req.GetToken()，解析token，获取userId
	// HACK userId暂时定为固定值1
	// toUserId为对方id，发送消息操作，即为消息接收者
	var fromUserId int64 = 1
	toUserId := request.GetToUserId()
	content := request.GetContent()
	msg := model.Message{
		ToUserId:   toUserId,
		FromUserId: fromUserId,
		Content:    content,
	}
	err := dao.SendMessage(msg)
	if err != nil {
		return &pb.DouyinMessageActionResponse{
			StatusCode: common.ErrorSendCode,
			StatusMsg:  common.ErrorSendMsg,
		}, err
	}
	return &pb.DouyinMessageActionResponse{
		StatusCode: common.SuccessCode,
		StatusMsg:  common.SuccessMsg,
	}, nil
}

func messageListToPbMessageList(msg []model.Message) []*pb.Message {
	var messageList []*pb.Message
	for _, v := range msg {
		messageList = append(messageList, &pb.Message{
			Id:         v.Id,
			FromUserId: v.FromUserId,
			ToUserId:   v.ToUserId,
			Content:    v.Content,
			CreateTime: v.CreateTime,
		})
	}
	return messageList
}

func validateToken(token string) int64 {
	//TODO 解析token
	//HACK parseToken(token)
	panic("implement me")
}
