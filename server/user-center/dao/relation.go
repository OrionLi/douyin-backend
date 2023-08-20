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

// followExists 检查是否已经存在关注关系
func followExists(followID, followerID int64) (bool, error) {
	err := db.Where("follow_id = ? and follower_id = ?", followID, followerID).First(&model.Relation{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

// FollowAction 关注用户
func FollowAction(selfUserID, toUserID int64) error {
	if selfUserID == toUserID {
		return errors.New("you cannot follow yourself")
	}

	exists, err := followExists(toUserID, selfUserID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("you have followed this user")
	}

	relation := model.Relation{
		Follow:   toUserID,
		Follower: selfUserID,
	}
	err = db.Create(&relation).Error
	if err != nil {
		return err
	}

	go CacheChangeUserCount(selfUserID, 1, "follow")
	go CacheChangeUserCount(toUserID, 1, "follower")

	return nil
}

// UnFollowAction 取消关注用户
func UnFollowAction(selfUserID, toUserID int64) error {
	err := db.Where("follow_id = ? and follower_id = ?", toUserID, selfUserID).Delete(&model.Relation{}).Error
	if err != nil {
		return err
	}

	go CacheChangeUserCount(selfUserID, -1, "follow")
	go CacheChangeUserCount(toUserID, -1, "follower")

	return nil
}

// GetFollowList 获取我关注的博主
func GetFollowList(userID int64) ([]*model.Relation, error) {
	relationList := make([]*model.Relation, 0)
	err := db.Where("follower_id = ?", userID).Find(&relationList).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return relationList, nil
	} else if err != nil {
		return nil, err
	}
	return relationList, nil
}

// GetFollowerList 获取关注我的粉丝
func GetFollowerList(userID int64) ([]*model.Relation, error) {
	relationList := make([]*model.Relation, 0)
	err := db.Where("follow_id = ?", userID).Find(&relationList).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return relationList, nil
	} else if err != nil {
		return nil, err
	}
	return relationList, nil
}

// IsFollow 判断是否已经关注过某用户
func IsFollow(selfUserID, toUserID int64) (bool, error) {
	if selfUserID == toUserID {
		return true, nil
	}

	exists, err := followExists(toUserID, selfUserID)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// CacheChangeUserCount 更新缓存中用户的关注或粉丝数量
func CacheChangeUserCount(userID int64, count int, category string) {
	// 假设缓存使用 Redis 存储，使用第三方库进行操作，这里仅作示例
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
