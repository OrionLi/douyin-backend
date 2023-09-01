package dao

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"user-center/model"
)

type RelationDao struct {
	*gorm.DB
}

func NewRelationDao(ctx context.Context) *RelationDao {
	return &RelationDao{NewDBClient(ctx)}
}

// Follow 关注
func (d *RelationDao) Follow(userId, toUserId int64) error {
	fmt.Printf("userid：%v，toUserId：%v", userId, toUserId)
	return d.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Model(&model.User{Model: gorm.Model{ID: uint(userId)}}).
		Association("Follows").
		Append(&model.User{
			Model: gorm.Model{ID: uint(toUserId)},
		})
}

// Unfollow 取消关注
func (d *RelationDao) Unfollow(userId, toUserId int64) error {
	return d.DB.
		Model(&model.User{
			Model: gorm.Model{ID: uint(userId)},
		}).
		Association("Follows").
		Delete(&model.User{
			Model: gorm.Model{ID: uint(toUserId)},
		})
}

// GetFollowList 获取关注列表
func (d *RelationDao) GetFollowList(userId int64) ([]*model.User, error) {
	var user *model.User
	if err := d.
		Where("id = ?", userId).
		Preload("Follows").
		Find(&user).Error; err != nil {
		return nil, err
	}
	return user.Follows, nil
}

// GetFollowerList 获取粉丝列表
func (d *RelationDao) GetFollowerList(userId int64) ([]*model.User, error) {
	var user *model.User
	if err := d.
		Where("id = ?", userId).
		Preload("Fans").
		Find(&user).Error; err != nil {
		return nil, err
	}
	return user.Fans, nil
}

// GetFriendList 获取好友列表
func (d *RelationDao) GetFriendList(userId int64) ([]*model.User, error) {

	// 获取粉丝的交集
	var friends []*model.User
	db.
		Raw("SELECT * FROM user WHERE id IN (SELECT follow_id FROM follows WHERE user_id = ?) AND id IN (SELECT user_id FROM follows WHERE follow_id = ?)",
			userId, userId).
		Scan(&friends)
	return friends, nil
}

// UpdateFanCount 更新 FanCount
func (d *RelationDao) UpdateFanCount(userID int64, newFanCount int) error {
	user := &model.User{}
	if err := d.First(user, userID).Error; err != nil {
		return err
	}

	user.FanCount = int64(newFanCount)
	if err := d.Save(user).Error; err != nil {
		return err
	}

	return nil
}

// UpdateFollowCount 更新 FollowCount
func (d *RelationDao) UpdateFollowCount(userID int64, newFollowCount int) error {
	user := &model.User{}
	if err := d.First(user, userID).Error; err != nil {
		return err
	}

	user.FollowCount = int64(newFollowCount)
	if err := d.Save(user).Error; err != nil {
		return err
	}

	return nil
}
