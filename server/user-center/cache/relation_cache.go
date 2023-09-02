package cache

import (
	"context"
	"time"
)

func NewRelationCache(ctx context.Context) *RedisCache {
	return &RedisCache{NewRedisClient(ctx)}
}

// UpdateFollowCountCache 更新关注数缓存，并设置过期时间为300秒
func (c *RedisCache) UpdateFollowCountCache(userId int64, num int64) error {
	// 设置缓存值
	err := c.HSet(context.Background(), GenUserInfoCacheKey(uint(userId)), FollowCount, num).Err()
	if err != nil {
		return err
	}

	// 设置过期时间为300秒
	key := GenUserInfoCacheKey(uint(userId))
	err = c.Expire(context.Background(), key, 10*time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}

// UpdateFollowerCountCache 更新粉丝数缓存，并设置过期时间为300秒
func (c *RedisCache) UpdateFollowerCountCache(userId int64, num int64) error {
	// 设置缓存值
	err := c.HSet(context.Background(), GenUserInfoCacheKey(uint(userId)), FanCount, num).Err()
	if err != nil {
		return err
	}

	// 设置过期时间为300秒
	key := GenUserInfoCacheKey(uint(userId))
	err = c.Expire(context.Background(), key, 10*time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}

// GetFollowCountCache 获取关注数缓存
func (c *RedisCache) GetFollowCountCache(userId int64) (int64, error) {
	return c.HGet(context.Background(), GenUserInfoCacheKey(uint(userId)), FollowCount).Int64()
}

// GetFollowerCountCache 获取粉丝数缓存
func (c *RedisCache) GetFollowerCountCache(userId int64) (int64, error) {
	return c.HGet(context.Background(), GenUserInfoCacheKey(uint(userId)), FanCount).Int64()
}
