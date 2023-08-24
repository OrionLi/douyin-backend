package dao

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"user-center/cache"
	"user-center/model"
)

type RelationDao struct {
	*gorm.DB
}

func NewRelationDao(ctx context.Context) *RelationDao {
	return &RelationDao{NewDBClient(ctx)}
}

// followExists 检查是否已经存在关注关系
func followExists(followID, followerID int64) (bool, error) {
	err := db.Model(&model.User{}).Where("id = ? AND id IN (SELECT follow_id FROM follow WHERE follower_id = ?)", followID, followerID).First(&model.User{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

// FollowAction 关注用户
func (u *RelationDao) FollowAction(currentUserId, targetUserId int64) error {
	if currentUserId == targetUserId {
		return errors.New("you cannot follow yourself")
	}

	exists, err := followExists(targetUserId, currentUserId)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("you have followed this user")
	}

	err = u.Model(&model.User{Model: gorm.Model{ID: uint(currentUserId)}}).Association("Follows").Append(&model.User{Model: gorm.Model{ID: uint(targetUserId)}}) // 添加到关注列表
	if err != nil {
		return err
	}

	go cache.CacheChangeUserCount(currentUserId, 1, "follow")
	go cache.CacheChangeUserCount(targetUserId, 1, "follower")

	return nil
}

// UnFollowAction 取消关注用户
func (u *RelationDao) UnFollowAction(currentUserId, targetUserId int64) error {
	err := u.Model(&model.User{Model: gorm.Model{ID: uint(currentUserId)}}).Association("Follows").Delete(&model.User{Model: gorm.Model{ID: uint(targetUserId)}}) // 从关注列表中删除
	if err != nil {
		return err
	}
	go cache.DelCacheFollow(uint(currentUserId), uint(targetUserId))
	go cache.CacheChangeUserCount(currentUserId, -1, "follow")
	go cache.CacheChangeUserCount(targetUserId, -1, "follower")

	return nil
}

// GetFollowList 获取我关注的博主
func (u *RelationDao) GetFollowList(userID int64) ([]*model.User, error) {
	var user model.User
	// 可直接获取id
	// u.Model(&model.User{Model: gorm.Model{ID: uint(userID)}}).Association("Follows").Find(&ids)
	err := u.Preload("Follows").First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return user.Follows, nil
}

// GetFollowerList 获取关注我的粉丝
func (u *RelationDao) GetFollowerList(userID int64) ([]*model.User, error) {
	var user model.User
	//获取粉丝id
	// u.Model(&model.User{Model: gorm.Model{ID: userID}}).Association("Fans").Find(&ids)
	err := u.Preload("Fans").First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return user.Fans, nil
}

// IsFollow 判断是否已经关注过某用户
func (u *RelationDao) IsFollow(currentUserId, targetUserId int64) (bool, error) {
	if currentUserId == targetUserId {
		return true, nil
	}

	exists, err := followExists(targetUserId, currentUserId)
	if err != nil {
		return false, err
	}
	return exists, nil
}
