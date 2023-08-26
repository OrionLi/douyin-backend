package e

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CustomError 自定义的错误详情结构体
type CustomError struct {
	Code codes.Code `json:"code"`
	Msg  string     `json:"msg"`
}

var MsgFlags = map[codes.Code]string{
	Success:                    "success",
	Error:                      "fail",
	InvalidParams:              "参数错误",
	ErrorAuthToken:             "token 认证失败",
	ErrorAuthCheckTokenTimeout: "token时效已过，请重新登录",
}

//GstMag 获取状态码对应信息

func GetMsg(code codes.Code) string {
	msg, ok := MsgFlags[code]
	if !ok {
		return MsgFlags[1000]
	}
	return msg
}

// NewError Grpc错误封装
func NewError(code codes.Code) error {
	c := &CustomError{
		Code: code,
		Msg:  GetMsg(code),
	}

	// 使用status.Newf函数来创建一个新的status.Status类型的错误，并传入c作为格式化参数
	st := status.Newf(code, "%v", c)
	return st.Err()
}

/*
调用方解析

	if err != nil {
			// 将错误转换为status.Status
			st, _ := status.FromError(err)
			// 获取错误码和错误信息
			code := st.Code()
			msg := st.Message()
			fmt.Println("code:", code)
			fmt.Println("msg:", msg)
		}
*/
func (e CustomError) Error() string {
	return e.Msg
}
