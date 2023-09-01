package cache

import (
	"context"
)

func NewRelationCache(ctx context.Context) *RedisCache {
	return &RedisCache{NewRedisClient(ctx)}
}

// UpdateFollowCountCache 更新关注数缓存
func (c *RedisCache) UpdateFollowCountCache(userId int64, num int64) error {
	err := c.HSet(context.Background(), GenUserInfoCacheKey(uint(userId)), FollowCount, num).Err()
	return err
}

// UpdateFollowerCountCache 更新粉丝数缓存
func (c *RedisCache) UpdateFollowerCountCache(userId int64, num int64) error {
	err := c.HSet(context.Background(), GenUserInfoCacheKey(uint(userId)), FanCount, num).Err()
	return err
}

// GetFollowCountCache 获取关注数缓存
func (c *RedisCache) GetFollowCountCache(userId int64) (int64, error) {
	return c.HGet(context.Background(), GenUserInfoCacheKey(uint(userId)), FollowCount).Int64()
}

// GetFollowerCountCache 获取粉丝数缓存
func (c *RedisCache) GetFollowerCountCache(userId int64) (int64, error) {
	return c.HGet(context.Background(), GenUserInfoCacheKey(uint(userId)), FanCount).Int64()
}
