package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
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
	favoriteCount, err := RedisClient.Get(context.Background(), favoriteKey).Int64()
	if err != nil {
		if err != redis.Nil {
			return err
		}
		// TODO 查询mysql中的值
	}
	favoriteCount++
	RedisClient.Set(context.Background(), favoriteKey, favoriteCount, 3*time.Minute)
	// TODO 异步更新mysql中的值
	return nil
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
