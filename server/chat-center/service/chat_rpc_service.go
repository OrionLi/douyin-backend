package service

import (
	"chat-center/dao"
	"chat-center/model"
	"chat-center/pkg/common"
	"context"
	"fmt"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"github.com/sony/sonyflake"
	"time"
)

type ChatRPCService struct {
	pb.UnimplementedDouyinMessageServiceServer
}

func NewChatRPCService() *ChatRPCService {
	return &ChatRPCService{}
}

func (s *ChatRPCService) GetMessage(ctx context.Context, request *pb.DouyinMessageChatRequest) (*pb.DouyinMessageChatResponse, error) {
	fromUserId := request.GetSelfUserId()
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
	fromUserId := request.GetSelfUserId()
	toUserId := request.GetToUserId()
	content := request.GetContent()

	// 创建一个 SonyFlake 实例
	sf := sonyflake.NewSonyflake(sonyflake.Settings{})

	// 生成唯一 ID
	id, err := sf.NextID()
	if err != nil {
		fmt.Printf("Error generating ID: %s\n", err)
		return nil, err
	}
	now := time.Now()
	format := now.Format("2006-01-02 15:04:05")
	message := model.Message{
		Id:         int64(id),
		ToUserId:   toUserId,
		FromUserId: fromUserId,
		Content:    content,
		CreateTime: format,
	}

	err = dao.SendMessage(message)
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
