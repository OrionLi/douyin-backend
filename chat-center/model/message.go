package model

// Message 消息结构体
type message struct {
	Id         int64
	ToUserId   int64
	FromUserId int64
	Content    string
	CreateTime string
}
