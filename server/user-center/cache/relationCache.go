package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
	"user-center/pkg/util"
)

func NewRelationCache(ctx context.Context) *RedisCache {
	return &RedisCache{NewRedisClient(ctx)}
}

// CacheChangeUserCount 更新缓存中用户的关注或粉丝数量
func CacheChangeUserCount(userID int64, count int, category string) {
	cache := NewRelationCache(context.Background())

	// 根据 category 构建缓存键
	key := CacheChangeUserCountKey(userID, category)

	// 获取当前缓存中的数量
	currentCountStr, err := cache.Get(context.Background(), key).Result()
	if err == redis.Nil {
		// 如果缓存中不存在该键，则默认设置为 0
		currentCountStr = "0"
	} else if err != nil {
		// 处理其他 Redis 错误
		util.LogrusObj.Error("<relationCache> : ", err)
		return
	}

	// 将字符串转换为整数
	currentCount := util.StrToInt64(currentCountStr)

	// 更新数量
	newCount := int(currentCount) + count

	// 更新到缓存中
	err = cache.Set(context.Background(), key, strconv.Itoa(newCount), 0).Err()
	if err != nil {
		util.LogrusObj.Error("<relationCache> : ", err)
		return
	}
}

// 取关时从缓存删除Key
func DelCacheFollow(uId, followId uint) error {
	cache := NewRelationCache(context.Background())
	return cache.
		Del(context.Background(), GenFollowUserCacheKey(uId, followId)).
		Err()
}
