package model

import "gorm.io/gorm"

const (
	CelebrityStandard = 0 // 设定粉丝超过300为网红，进行缓存处理
)

type User struct {
	gorm.Model
	Username    string `gorm:"not null;unique;index" `
	Password    string `gorm:"not null"`
	FollowCount int64
	FanCount    int64

	// 多对多
	Follows []*User `gorm:"many2many:follows;"`                         // 关注列表
	Fans    []*User `gorm:"many2many:follows;joinForeignKey:follow_id"` // 粉丝列表
}

func (user *User) IsCelebrity() bool {
	return user.FanCount >= CelebrityStandard
}
