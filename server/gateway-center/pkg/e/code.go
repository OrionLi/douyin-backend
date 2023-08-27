package e

const (
	Success = 0
	Error   = 404

	ErrorAuthToken             = 408
	ErrorAuthCheckTokenTimeout = 409

	InvalidParams          = 501
	ServiceErr             = 10001
	ParamErr               = 10002
	AuthorizationFailedErr = 10004
	UserNotExistErr        = 10006
	TokenErr               = 10007
	CommentPosting         = 10009 // 发布评论失败
	DeleteComment          = 10010 // 删除评论失败
	NoMyComment            = 10011 // 不是自己的评论
	NoCommentExists        = 10012 // 评论不存在
	FavListEmpty           = 10013 // 喜欢列表为空
	FavActionErr           = 10014 // 点赞操作失败
	FavCountErr            = 10015 // 获取点赞数量失败
	ListComment            = 10016
	FailedToCallRpc        = 10017
)
