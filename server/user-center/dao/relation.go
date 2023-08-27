package dao

import (
	"context"
	"user-center/cache"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	"user-center/model"
)

type RelationDao struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewRelationDao(db *gorm.DB, rdb *redis.Client) *RelationDao {
	return &RelationDao{db: db, rdb: rdb}
}

// Follow 关注
func (d *RelationDao) Follow(ctx context.Context, userId, toUserId uint) error {
	return d.db.Model(&model.User{}).Where("id = ?", userId).Association("Follows").Append(&model.User{
		Model: gorm.Model{ID: toUserId},
	})
}

// Unfollow 取消关注
func (d *RelationDao) Unfollow(ctx context.Context, userId, toUserId uint) error {
	return d.db.Model(&model.User{}).Where("id = ?", userId).Association("Follows").Delete(&model.User{
		Model: gorm.Model{ID: toUserId},
	})
}

// GetFollowList 获取关注列表
func (d *RelationDao) GetFollowList(ctx context.Context, userId uint) ([]*model.User, error) {
	var users []*model.User
	if err := d.db.Model(&model.User{}).Where("id = ?", userId).Association("Follows").Find(&users); err != nil {
		return nil, err
	}
	return users, nil
}

// GetFollowerList 获取粉丝列表
func (d *RelationDao) GetFollowerList(ctx context.Context, userId uint) ([]*model.User, error) {
	var users []*model.User
	if err := d.db.Model(&model.User{}).Where("id = ?", userId).Association("Fans").Find(&users); err != nil {
		return nil, err
	}
	return users, nil
}

// GetFriendList 获取好友列表
func (d *RelationDao) GetFriendList(ctx context.Context, userId uint) ([]*model.User, error) {
	// 通过双向关联查询获取互相关注好友
	var friends []*model.User
	if err := d.db.Model(&model.User{}).
		Where("id IN (?)", d.db.Model(&model.User{}).Where("id = ?", userId).Select("follow_id")).
		Where("follow_id IN (?)", d.db.Model(&model.User{}).Where("id = ?", userId).Select("id")).
		Find(&friends).Error; err != nil {
		return nil, err
	}
	return friends, nil
}

// IncrFollowCountCache 更新关注数缓存
func (d *RelationDao) IncrFollowCountCache(userId, incr uint) error {
	return d.rdb.IncrBy(context.Background(), cache.GenUserInfoCacheKey(userId), int64(incr)).Err()
}

// IncrFollowerCountCache 更新粉丝数缓存
func (d *RelationDao) IncrFollowerCountCache(userId, incr uint) error {
	return d.rdb.IncrBy(context.Background(), cache.GenUserInfoCacheKey(userId), int64(incr)).Err()
}

// GetFollowCountCache 获取关注数缓存
func (d *RelationDao) GetFollowCountCache(userId uint) (int64, error) {
	return d.rdb.Get(context.Background(), cache.GenUserInfoCacheKey(userId)).Int64()
}

// GetFollowerCountCache 获取粉丝数缓存
func (d *RelationDao) GetFollowerCountCache(userId uint) (int64, error) {
	return d.rdb.Get(context.Background(), cache.GenUserInfoCacheKey(userId)).Int64()
}
