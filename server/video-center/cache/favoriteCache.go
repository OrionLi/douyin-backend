package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
	"video-center/dao"
	"video-center/pkg/util"
)

const FavoriteUpdateSetKey = "fav_update_set"
const retryInterval = 5 * time.Millisecond
const maxRetryInterval = 50 * time.Millisecond

// ActionFavoriteCache 点赞缓存
func ActionFavoriteCache(videoId int64, actionType int32) error {
	lockKey := fmt.Sprintf("lock:fav:vid:%d", videoId)
	favoriteKey := fmt.Sprintf("favorite:%d", videoId)
	lock, err := RedisLock(fmt.Sprintf(lockKey, videoId), 1*time.Second)
	if err != nil {
		return err
	}
	if !lock {
		// FIXME 重试机制优化
		var retryDelay = retryInterval
		for {
			lock, err := RedisLock(fmt.Sprintf(lockKey, videoId), 1*time.Second)
			if err != nil {
				return err
			}
			if lock {
				break // 成功获取锁，退出重试循环
			}
			// 获取锁失败，等待一段时间后重试
			time.Sleep(retryDelay)
			if retryDelay < maxRetryInterval {
				retryDelay += 10 * time.Millisecond
			}
		}
	}
	defer RedisUnlock(fmt.Sprintf(lockKey, videoId))
	// 查询 Redis 中的值
	favoriteCount, err := GetFavoriteCountCache(videoId)
	if err != nil {
		if err == redis.Nil {
			count, err := dao.GetSingleVideoFavoriteCount(context.Background(), videoId)
			if err != nil {
				return err
			}
			favoriteCount = int64(count)
		} else {
			return err
		}
	}
	switch actionType {
	case 1:
		favoriteCount++
	case 2:
		favoriteCount--
	default:
		return errors.New("actionType error")
	}
	// 更新缓存
	RedisClient.Set(context.Background(), favoriteKey, favoriteCount, 7*24*time.Hour)
	// 更新集合
	RedisClient.SAdd(context.Background(), FavoriteUpdateSetKey, videoId)
	return nil
}

// SetFavoriteCountCache 设置缓存中的某个视频点赞数量
func SetFavoriteCountCache(videoId int64, favoriteCount int64) error {
	favoriteKey := fmt.Sprintf("favorite:%d", videoId)
	return RedisClient.Set(context.Background(), favoriteKey, favoriteCount, 7*24*time.Hour).Err()
}

// GetFavoriteCountCache 获取缓存中的某个视频点赞数量
func GetFavoriteCountCache(videoId int64) (int64, error) {
	favoriteKey := fmt.Sprintf("favorite:%d", videoId)
	return RedisClient.Get(context.Background(), favoriteKey).Int64()
}

// DeleteVideoIdFromFavoriteUpdateSet 从更新集合中删除某个视频ID
func DeleteVideoIdFromFavoriteUpdateSet(videoId int64) error {
	return RedisClient.SRem(context.Background(), FavoriteUpdateSetKey, videoId).Err()
}

// GetMemberFromFavoriteUpdateSet 获取更新集合中的所有视频ID
func GetMemberFromFavoriteUpdateSet() ([]string, error) {
	return RedisClient.SMembers(context.Background(), FavoriteUpdateSetKey).Result()
}

// RedisLock redis分布式锁
func RedisLock(lockKey string, lockTimeout time.Duration) (bool, error) {
	lockAcquired, err := RedisClient.SetNX(context.Background(), lockKey, "lock-true", lockTimeout).Result()
	if err != nil {
		return false, err
	}

	return lockAcquired, nil
}

// RedisUnlock redis分布式锁解锁
func RedisUnlock(lockKey string) {
	_, err := RedisClient.Del(context.Background(), lockKey).Result()
	if err != nil {
		util.LogrusObj.Error("<Redis-FavoriteAction>, Unlock failed", err)
	}
}
