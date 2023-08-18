package dao

import (
	"context"
	"gorm.io/gorm"
)

type User struct {
	Id          int64  `json:"id,omitempty"`             // 用户id
	Username    string `json:"name,omitempty"`           // 用户名称
	FollowCount int64  `json:"follow_count,omitempty"`   // 关注总数
	FanCount    int64  `json:"follower_count,omitempty"` // 粉丝总数
	IsFollow    bool   `json:"is_follow,omitempty"`      // true-已关注，false-未关注
}
type TUser struct {
	gorm.Model
	Username    string `gorm:"not null;unique;index" `
	Password    string `gorm:"not null"`
	FollowCount int64
	FanCount    int64

	// 多对多
	Follows []*TUser `gorm:"many2many:follows;"`                         // 关注列表
	Fans    []*TUser `gorm:"many2many:follows;joinForeignKey:follow_id"` // 粉丝列表
}

func (u *TUser) TableName() string {
	return "user"
}

func (u *User) TableName() string {
	return "user"
}

func QueryUserByID(ctx context.Context, userId int64) (*User, error) {
	user := new(User)
	err := DB.WithContext(ctx).Where("id = ?", userId).First(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}
func ExistID(ctx context.Context, userId int64) bool {
	user := new(User)
	err := DB.WithContext(ctx).Where("id = ?", userId).First(&user).Error
	if err != nil {
		return false
	}
	return true
}
func IsFanOf(userID1 int64, userID2 uint) (bool, error) {
	var user1, user2 TUser

	// 查询用户1
	if err := DB.First(&user1, userID1).Error; err != nil {
		return false, err
	}

	// 查询用户2
	if err := DB.First(&user2, userID2).Error; err != nil {
		return false, err
	}

	// 遍历用户1的粉丝列表，检查是否存在用户2
	for _, fan := range user1.Fans {
		if fan.ID == userID2 {
			return true, nil
		}
	}

	return false, nil
}
