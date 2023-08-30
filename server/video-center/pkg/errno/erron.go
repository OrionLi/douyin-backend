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
	CommentPostingCode         = 10009 // 发布评论失败
	DeleteCommentCode          = 10010 // 删除评论失败
	NoMyCommentCode            = 10011 // 不是自己的评论
	NoCommentExistsCode        = 10012 // 评论不存在
	FavListEmptyCode           = 10013 // 喜欢列表为空
	FavActionErrCode           = 10014 // 点赞操作失败
	FavCountErrCode            = 10015 // 获取点赞数量失败
	ListCommentCode            = 10016
	FailedToCallRpcCode        = 10017
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
	CommentPostingErr      = NewErrno(CommentPostingCode, "Failed to post a comment")
	DeleteCommentErr       = NewErrno(DeleteCommentCode, "Failed to delete comment")
	NoMyCommentErr         = NewErrno(NoMyCommentCode, "Not your comment")
	NoCommentExistsErr     = NewErrno(NoCommentExistsCode, "留下第一条评论吧！")
	FavListEmptyErr        = NewErrno(FavListEmptyCode, "Like the list to be empty")
	FavActionErr           = NewErrno(FavActionErrCode, "Like operation failed")
	FavCountErr            = NewErrno(FavCountErrCode, "Failed to get number of likes")
	ListCommentErr         = NewErrno(ListCommentCode, "Failed to query the comment list")
	FailedToCallRpcErr     = NewErrno(FailedToCallRpcCode, "Failed to call rpc")
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
