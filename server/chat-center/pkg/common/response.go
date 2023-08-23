package common

import (
	"github.com/OrionLi/douyin-backend/pkg/pb"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type GetMessageResponse struct {
	Response
	MessageList []*pb.Message `json:"message_list"`
}

type SendMessageResponse struct {
	Response
}
