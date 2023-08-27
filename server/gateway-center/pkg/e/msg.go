package e

import "errors"

// CustomError 自定义的错误详情结构体
type CustomError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

var MsgFlags = map[int]string{
	Success:                    "success",
	Error:                      "fail",
	InvalidParams:              "参数错误",
	ErrorAuthToken:             "token 认证失败",
	ErrorAuthCheckTokenTimeout: "token时效已过，请重新登录",
	ServiceErr:                 "Service is unable to start successfully",
	ParamErr:                   "Wrong Parameter has been given",
	AuthorizationFailedErr:     "Authorization failed",
	UserNotExistErr:            "User does not exists",
	TokenErr:                   "Token confirm wrong",
	CommentPosting:             "Failed to post a comment",
	DeleteComment:              "Failed to delete comment",
	NoMyComment:                "Not your comment",
	NoCommentExists:            "Comment does not exist",
	FavListEmpty:               "Like the list to be empty",
	FavActionErr:               "Like operation failed",
	FavCountErr:                "Failed to get number of likes",
	ListComment:                "Failed to query the comment list",
	FailedToCallRpc:            "Failed to call rpc",
}

//GstMag 获取状态码对应信息

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if !ok {
		return MsgFlags[1000]
	}
	return msg
}
func NewCustomError(code int64, msg string) CustomError {
	return CustomError{int(code), msg}
}
func ConvertErr(err error) CustomError {
	Err := CustomError{}
	if errors.As(err, &Err) {
		return Err
	}

	s := NewCustomError(ServiceErr, GetMsg(ServiceErr))
	s.Msg = err.Error()
	return s
}

/*// NewError Grpc错误封装
func NewError(code codes.Code) error {
	c := &CustomError{
		Code: code,
		Msg:  GetMsg(code),
	}

	// 使用status.Newf函数来创建一个新的status.Status类型的错误，并传入c作为格式化参数
	st := status.Newf(code, "%v", c)
	return st.Err()
}*/

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
