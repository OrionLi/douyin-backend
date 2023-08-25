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

func ActionFavoriteCache(videoId int64, actionType int32) error {
	lockKey := fmt.Sprintf("lock:fav:vid:%d", videoId)
	favoriteKey := fmt.Sprintf("favorite:%d", videoId)
	lock, err := RedisLock(fmt.Sprintf(lockKey, videoId), 3*time.Second)
	if err != nil || !lock {
		return err
	}
	defer RedisUnlock(fmt.Sprintf(lockKey, videoId))
	// 查询 Redis 中的值
	favoriteCount, err := GetFavoriteCountCache(videoId)
	if err != nil {
		if err != redis.Nil {
			return err
		}
		count, err := dao.GetSingleVideoFavoriteCount(context.Background(), videoId)
		if err != nil {
			return err
		}
		favoriteCount = int64(count)
	}
	switch actionType {
	case 1:
		favoriteCount++
	case 2:
		favoriteCount--
	default:
		return errors.New("actionType error")
	}
	RedisClient.Set(context.Background(), favoriteKey, favoriteCount, 3*time.Minute)
	updateDone := make(chan error)
	// FIXME 异步更新 MySQL 中的值，此过程可能会有并发问题，考虑更改为定时更新
	go dao.UpdateMySQLFavoriteCount(videoId, favoriteCount, updateDone)
	// 等待异步更新完成，如果有错误，返回错误
	if err := <-updateDone; err != nil {
		util.LogrusObj.Error("<Redis-FavoriteAction>, Async Update MySQL failed", err)
		return err
	}
	return nil
}

func GetFavoriteCountCache(videoId int64) (int64, error) {
	favoriteKey := fmt.Sprintf("favorite:%d", videoId)
	return RedisClient.Get(context.Background(), favoriteKey).Int64()
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
