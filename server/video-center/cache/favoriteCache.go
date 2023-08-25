package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
	"video-center/dao"
	"video-center/pkg/util"
)

var FavoriteUpdateSetKey = "fav_update_set"

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
	RedisClient.Set(context.Background(), favoriteKey, favoriteCount, -1)
	RedisClient.SAdd(context.Background(), FavoriteUpdateSetKey, videoId)
	return nil
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

// UpdateFavoriteCacheToMySQLAtRegularTime 更新到MySQL
func UpdateFavoriteCacheToMySQLAtRegularTime() {
	interval := 12 * time.Hour // 设置定时任务的时间间隔
	// 创建一个定时器
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			favoriteUpdateSet, err := RedisClient.SMembers(context.Background(), FavoriteUpdateSetKey).Result()
			if err != nil {
				util.LogrusObj.Error("<Favorite Count Update failed>", ": Get list fail", err)
			}
			// 处理每个视频ID
			for _, videoIdStr := range favoriteUpdateSet {
				videoId := util.StringToInt64(videoIdStr)
				count, err := GetFavoriteCountCache(videoId)
				if err != nil {
					util.LogrusObj.Error("<Favorite Count Update failed> ", "videoId:", videoId, "err:", err)
				}
				go dao.UpdateMySQLFavoriteCount(videoId, count)
			}
			log.Println("UpdateToMySQL task executed at", time.Now())
		}
	}
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
