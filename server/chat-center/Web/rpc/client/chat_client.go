package client

import (
	"context"
	"github.com/OrionLi/douyin-backend/pkg/pb"
)

func GetMessage(selfUserId int64, interActiveId int64, preMsgTime int64) (*pb.DouyinMessageChatResponse, error) {
	conn := GetMsgConn()

	client := pb.NewDouyinMessageServiceClient(conn)

	return client.GetMessage(context.Background(), &pb.DouyinMessageChatRequest{
		SelfUserId: selfUserId,
		ToUserId:   interActiveId,
		PreMsgTime: preMsgTime,
	})
}

func SendMessage(selfUserId int64, interActiveId int64, content string) (*pb.DouyinMessageActionResponse, error) {
	conn := GetMsgConn()

	client := pb.NewDouyinMessageServiceClient(conn)

	return client.SendMessage(context.Background(), &pb.DouyinMessageActionRequest{
		SelfUserId: selfUserId,
		ToUserId:   interActiveId,
		ActionType: 1,
		Content:    content,
	})
}
