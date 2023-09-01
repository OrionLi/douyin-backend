package cache

import (
	"context"
	"time"
)

func NewRelationCache(ctx context.Context) *RedisCache {
	return &RedisCache{NewRedisClient(ctx)}
}

// UpdateFollowCountCache 更新关注数缓存
func (c *RedisCache) UpdateFollowCountCache(userId int64, num int64) error {
	err := c.Set(context.Background(), GenFollowCountCacheKey(uint(userId)), num, time.Hour).Err()
	return err
}

// UpdateFollowerCountCache 更新粉丝数缓存
func (c *RedisCache) UpdateFollowerCountCache(userId int64, num int64) error {
	err := c.Set(context.Background(), GenFollowerCountCacheKey(uint(userId)), num, time.Hour).Err()
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
