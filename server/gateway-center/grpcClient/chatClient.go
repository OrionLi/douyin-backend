package grpcClient

import (
	"context"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"google.golang.org/grpc"
)

var Conn *grpc.ClientConn

func GetMessage(selfUserId int64, interActiveId int64, preMsgTime int64) (*pb.DouyinMessageChatResponse, error) {
	return ChatClient.GetMessage(context.Background(), &pb.DouyinMessageChatRequest{
		SelfUserId: selfUserId,
		ToUserId:   interActiveId,
		PreMsgTime: preMsgTime,
	})
}

func SendMessage(selfUserId int64, interActiveId int64, content string) (*pb.DouyinMessageActionResponse, error) {
	return ChatClient.SendMessage(context.Background(), &pb.DouyinMessageActionRequest{
		SelfUserId: selfUserId,
		ToUserId:   interActiveId,
		ActionType: 1,
		Content:    content,
	})
}
