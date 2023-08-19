package error

import (
	"chat-center/pkg/common"
	"fmt"
)

type ChatError struct {
	ErrCode int64
	ErrMsg  string
}

func (e ChatError) Error() string {
	return fmt.Sprintf("err_code=%d, err_msg=%s", e.ErrCode, e.ErrMsg)
}

func NewChatError(code int64, msg string) ChatError {
	return ChatError{code, msg}
}

func GetError(err error) ChatError {
	return NewChatError(common.ErrorGetCode, common.ErrorGetMsg+":"+err.Error())
}

func SendError(err error) ChatError {
	return NewChatError(common.ErrorSendCode, common.ErrorSendMsg+":"+err.Error())
}
