package response

import (
	"github.com/OrionLi/douyin-backend/pkg/pb"
)

type DouyinChatGetMessageResponse struct {
	StatusCode  int32         `json:"status_code"`
	StatusMsg   string        `json:"status_msg,omitempty"`
	MessageList []*pb.Message `json:"message_list,omitempty"`
}

type DouyinChatSendMessageResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}
