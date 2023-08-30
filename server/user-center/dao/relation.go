package dao

import (
	"context"
	"fmt"
	"gorm.io/gorm"

	"user-center/model"
)

// RelationDao TODO:去除无用ctx
type RelationDao struct {
	*gorm.DB
}

func NewRelationDao(ctx context.Context) *RelationDao {
	return &RelationDao{NewDBClient(ctx)}
}

// Follow 关注
func (d *RelationDao) Follow(ctx context.Context, userId, toUserId int64) error {
	fmt.Printf("userid：%v，toUserId：%v", userId, toUserId)
	return d.DB.Model(&model.User{}).Where("id = ?", userId).Association("Follows").Append(&model.User{
		Model: gorm.Model{ID: uint(toUserId)},
	})
}

// Unfollow 取消关注
func (d *RelationDao) Unfollow(userId, toUserId int64) error {
	return d.DB.Model(&model.User{}).Where("id = ?", userId).Association("Follows").Delete(&model.User{
		Model: gorm.Model{ID: uint(toUserId)},
	})
}

// GetFollowList 获取关注列表
func (d *RelationDao) GetFollowList(ctx context.Context, userId int64) ([]*model.User, error) {
	var users []*model.User
	if err := d.DB.Model(&model.User{}).Where("id = ?", userId).Association("Follows").Find(&users); err != nil {
		return nil, err
	}
	return users, nil
}

// GetFollowerList 获取粉丝列表
func (d *RelationDao) GetFollowerList(ctx context.Context, userId int64) ([]*model.User, error) {
	var users []*model.User
	if err := d.DB.Model(&model.User{}).Where("id = ?", userId).Association("Fans").Find(&users); err != nil {
		return nil, err
	}
	return users, nil
}

// GetFriendList 获取好友列表
func (d *RelationDao) GetFriendList(userId int64) ([]*model.User, error) {
	// 通过双向关联查询获取互相关注好友
	var friends []*model.User
	if err := d.DB.Model(&model.User{}).
		Where("id IN (?)", d.DB.Model(&model.User{}).Where("id = ?", userId).Select("follow_id")).
		Where("follow_id IN (?)", d.DB.Model(&model.User{}).Where("id = ?", userId).Select("id")).
		Find(&friends).Error; err != nil {
		return nil, err
	}
	return friends, nil
}

//// IncrFollowCountCache 更新关注数缓存
//func (d *RelationDao) IncrFollowCountCache(userId, incr int64) error {
//	return d.rdb.IncrBy(context.Background(), cache.GenUserInfoCacheKey(uint(userId)), incr).Err()
//}
//
//// IncrFollowerCountCache 更新粉丝数缓存
//func (d *RelationDao) IncrFollowerCountCache(userId, incr int64) error {
//	return d.rdb.IncrBy(context.Background(), cache.GenUserInfoCacheKey(uint(userId)), incr).Err()
//}
//
//// GetFollowCountCache 获取关注数缓存
//func (d *RelationDao) GetFollowCountCache(userId int64) (int64, error) {
//	return d.rdb.Get(context.Background(), cache.GenUserInfoCacheKey(uint(userId))).Int64()
//}
//
//// GetFollowerCountCache 获取粉丝数缓存
//func (d *RelationDao) GetFollowerCountCache(userId int64) (int64, error) {
//	return d.rdb.Get(context.Background(), cache.GenUserInfoCacheKey(uint(userId))).Int64()
//}
