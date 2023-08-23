package dao

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"log"
	"strconv"
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

	go CacheChangeUserCount(currentUserId, 1, "follow")
	go CacheChangeUserCount(targetUserId, 1, "follower")

	return nil
}

// UnFollowAction 取消关注用户
func (u *RelationDao) UnFollowAction(currentUserId, targetUserId int64) error {
	err := u.Model(&model.User{Model: gorm.Model{ID: uint(currentUserId)}}).Association("Follows").Delete(&model.User{Model: gorm.Model{ID: uint(targetUserId)}}) // 从关注列表中删除
	if err != nil {
		return err
	}

	go CacheChangeUserCount(currentUserId, -1, "follow")
	go CacheChangeUserCount(targetUserId, -1, "follower")

	return nil
}

// GetFollowList 获取我关注的博主
func (u *RelationDao) GetFollowList(userID int64) ([]*model.User, error) {
	var user model.User
	err := u.Preload("Follows").First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return user.Follows, nil
}

// GetFollowerList 获取关注我的粉丝
func (u *RelationDao) GetFollowerList(userID int64) ([]*model.User, error) {
	var user model.User
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

// CacheChangeUserCount 更新缓存中用户的关注或粉丝数量
func CacheChangeUserCount(userID int64, count int, category string) {
	cache := redis.NewClient(&redis.Options{
		Addr:     "r-bp12xmzrbjr36iq7lepd.redis.rds.aliyuncs.com:6379", // Redis 服务器地址
		Password: "Rh2004==",                                           // Redis 密码
		DB:       0,                                                    // Redis 数据库索引
	})

	// 根据 category 构建缓存键
	key := fmt.Sprintf("user:%d:%s", userID, category)

	// 获取当前缓存中的数量
	currentCountStr, err := cache.Get(context.Background(), key).Result()
	if err == redis.Nil {
		// 如果缓存中不存在该键，则默认设置为 0
		currentCountStr = "0"
	} else if err != nil {
		// 处理其他 Redis 错误
		log.Println("Redis error:", err)
		return
	}

	// 将字符串转换为整数
	currentCount, _ := strconv.Atoi(currentCountStr)

	// 更新数量
	newCount := currentCount + count

	// 更新到缓存中
	err = cache.Set(context.Background(), key, strconv.Itoa(newCount), 0).Err()
	if err != nil {
		log.Println("Redis error:", err)
		return
	}
}
