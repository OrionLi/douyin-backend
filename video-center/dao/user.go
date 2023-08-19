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

func (u *User) TableName() string {
	return "user"
}

type follow struct {
	UserId   uint
	followId uint
}

func (u *follow) TableName() string {
	return "follows"
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
	var fw follow
	result := DB.WithContext(context.Background()).Where("user_id = ? and follow_id = ?", userID1, userID2).First(&fw)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}
