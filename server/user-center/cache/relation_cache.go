package cache

import (
	"context"
	"time"
)

func NewRelationCache(ctx context.Context) *RedisCache {
	return &RedisCache{NewRedisClient(ctx)}
}

//// IncrFollowCountCache 更新关注数缓存
//func (c *RedisCache) IncrFollowCountCache(userId, incr int64) error {
//	return c.IncrBy(context.Background(), GenFollowCountCacheKey(uint(userId)), incr).Err()
//}
//
//// IncrFollowerCountCache 更新粉丝数缓存
//func (c *RedisCache) IncrFollowerCountCache(userId, incr int64) error {
//	return c.IncrBy(context.Background(), GenFollowerCountCacheKey(uint(userId)), incr).Err()
//}

// IncrFollowCountCache 更新关注数缓存
func (c *RedisCache) IncrFollowCountCache(userId int64, incr int64) error {
	key := GenFollowCountCacheKey(uint(userId))
	count, _ := c.GetFollowCountCache(userId)
	err := c.Set(context.Background(), key, count+incr, time.Hour).Err()
	return err
}

// IncrFollowerCountCache 更新粉丝数缓存
func (c *RedisCache) IncrFollowerCountCache(userId int64, incr int64) error {
	key := GenFollowerCountCacheKey(uint(userId))
	count, _ := c.GetFollowerCountCache(userId)
	err := c.Set(context.Background(), key, count+incr, time.Hour).Err()
	return err
}

// GetFollowCountCache 获取关注数缓存
func (c *RedisCache) GetFollowCountCache(userId int64) (int64, error) {
	return c.Get(context.Background(), GenFollowCountCacheKey(uint(userId))).Int64()
}

// GetFollowerCountCache 获取粉丝数缓存
func (c *RedisCache) GetFollowerCountCache(userId int64) (int64, error) {
	return c.Get(context.Background(), GenFollowerCountCacheKey(uint(userId))).Int64()
}
