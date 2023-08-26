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

// ActionFavoriteCache 点赞缓存
// 通过原子性操作解决并发问题
func ActionFavoriteCache(videoId int64, actionType int32) error {
	favoriteKey := fmt.Sprintf("favorite:%d", videoId)
	_, err := RedisClient.Get(context.Background(), favoriteKey).Int64()
	if err != nil {
		if err == redis.Nil {
			lock, err := RedisLock(fmt.Sprintf("fav:lock:vid:%d", videoId), 5*time.Second)
			defer RedisUnlock(fmt.Sprintf("fav:lock:vid:%d", videoId))
			if err != nil {
				return err
			}
			// 获得锁则设置键值对，其余请求等待
			if lock {
				count, err := dao.GetSingleVideoFavoriteCount(context.Background(), videoId)
				if err != nil {
					return err
				}
				// 设置缓存
				err = SetFavoriteCountCache(videoId, int64(count))
			} else {
				time.Sleep(100 * time.Millisecond)
			}
		} else {
			return err
		}
	}
	switch actionType {
	case 1:
		err := RedisClient.Incr(context.Background(), favoriteKey).Err()
		if err != nil {
			return err
		}
	case 2:
		err := RedisClient.Decr(context.Background(), favoriteKey).Err()
		if err != nil {
			return err
		}
	default:
		return errors.New("actionType error")
	}
	// 更新集合
	go RedisClient.SAdd(context.Background(), FavoriteUpdateSetKey, videoId)
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
