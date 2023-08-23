package server

import (
	"chat-center/pkg/common"
	"chat-center/pkg/utils"
	"chat-center/service"
	"context"
	"github.com/OrionLi/douyin-backend/pkg/pb"
)

type ChatServer struct {
	pb.UnimplementedDouyinMessageServiceServer
}

func NewChatServer() *ChatServer {
	return &ChatServer{}
}

func (s *ChatServer) GetMessage(_ context.Context, request *pb.DouyinMessageChatRequest) (*pb.DouyinMessageChatResponse, error) {
	fromUserId := request.GetSelfUserId()
	toUserId := request.GetToUserId()
	preMsgTime := request.GetPreMsgTime()
	if preMsgTime == 0 {
		messageList, err := service.NewChatService().GetAllHistoryMessage(fromUserId, toUserId)
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
		messageList, err := service.NewChatService().GetMessageByPreMsgTime(fromUserId, toUserId, preMsgTime)
		if err != nil {
			utils.LogrusObj.Error("<GetMessage-rpc> Get message error: ", err, " [be from req]:", request)
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

func (s *ChatServer) SendMessage(_ context.Context, request *pb.DouyinMessageActionRequest) (*pb.DouyinMessageActionResponse, error) {
	fromUserId := request.GetSelfUserId()
	toUserId := request.GetToUserId()
	content := request.GetContent()

	err := service.NewChatService().SendMessage(fromUserId, toUserId, content)

	if err != nil {
		utils.LogrusObj.Error("<SendMessage-rpc> Send message error: ", err, " [be from req]:", request)
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
