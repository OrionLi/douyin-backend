package errno

import (
	"errors"
	"fmt"
)

const (
	SuccessCode                = 0
	ServiceErrCode             = 10001
	ParamErrCode               = 10002
	UserAlreadyExistErrCode    = 10003
	AuthorizationFailedErrCode = 10004
	LoginErrCode               = 10005
	UserNotExistErrCode        = 10006
	TokenErrCode               = 10007
	FollowCode                 = 10008
)

type Errno struct {
	ErrCode int64
	ErrMsg  string
}

func (e Errno) Error() string {
	return fmt.Sprintf("err_code=%d, err_msg=%s", e.ErrCode, e.ErrMsg)
}

func NewErrno(code int64, msg string) Errno {
	return Errno{code, msg}
}

func (e Errno) WithMessage(msg string) Errno {
	e.ErrMsg = msg
	return e
}

var (
	Success                = NewErrno(SuccessCode, "Success")
	ServiceErr             = NewErrno(ServiceErrCode, "Service is unable to start successfully")
	ParamErr               = NewErrno(ParamErrCode, "Wrong Parameter has been given")
	AuthorizationFailedErr = NewErrno(AuthorizationFailedErrCode, "Authorization failed")
	UserNotExistErr        = NewErrno(UserNotExistErrCode, "User does not exists")
	TokenErr               = NewErrno(TokenErrCode, "Token confirm wrong")
)

func ConvertErr(err error) Errno {
	Err := Errno{}
	if errors.As(err, &Err) {
		return Err
	}

	s := ServiceErr
	s.ErrMsg = err.Error()
	return s
}
