package e

import "google.golang.org/grpc/status"

// CustomError 自定义的错误详情结构体
type CustomError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

var MsgFlags = map[int]string{
	Error:                      "fail",
	InvalidParams:              "参数错误",
	ErrorExistUser:             "该用户名已存在",
	ErrorExistUserNotFound:     "用户不存在",
	ErrorNotCompare:            "密码错误",
	ErrorAuthToken:             "token 认证失败",
	ErrorAuthCheckTokenTimeout: "token 过期",
}

//GstMag 获取状态码对应信息

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if !ok {
		return MsgFlags[1000]
	}
	return msg
}
func NewError(code int) error {
	c := &CustomError{
		Code: code,
		Msg:  GetMsg(code),
	}
	st, _ := status.FromError(c)
	return st.Err()
}

/*
调用方解析

	if err != nil {
		// 将错误转换为status.Status
		st := status.FromError(err)
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
