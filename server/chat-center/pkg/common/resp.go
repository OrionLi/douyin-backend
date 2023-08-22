package common

import "chat-center/model"

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type GetMessageResponse struct {
	Response
	MessageList []model.Message `json:"message_list"`
}

type SendMessageResponse struct {
	Response
}
