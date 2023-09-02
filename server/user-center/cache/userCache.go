package cache

import (
	"context"
	"time"
)

func NewUserCache(ctx context.Context) *RedisCache {
	return &RedisCache{NewRedisClient(ctx)}
}

// IsFollow  是否关注
func (c *RedisCache) IsFollow(ctx context.Context, uId, followId uint) bool {
	return c.
		Exists(ctx, GenFollowUserCacheKey(uId, followId)).
		Val() == 1
}

// AddFollow 关注关系缓存
func (c *RedisCache) AddFollow(ctx context.Context, uId, followId uint) error {
	return c.
		Set(ctx, GenFollowUserCacheKey(uId, followId), 1, time.Hour).
		Err()
}

// DeleteFollow 取关关系缓存
func (c *RedisCache) DeleteFollow(ctx context.Context, uId, followId uint) error {
	return c.
		Del(ctx, GenFollowUserCacheKey(uId, followId)).
		Err()
}

// AddUser 用户信息缓存
func (c *RedisCache) AddUser(ctx context.Context, uId uint, m map[string]interface{}) error {
	err := c.HSet(ctx, GenUserInfoCacheKey(uId), m).Err()
	if err != nil {
		return err
	}
	// 设置键的过期时间为10分钟
	return c.Expire(ctx, GenUserInfoCacheKey(uId), 600*time.Second).Err()
}

// HasUser 判断Redis中是否存在某个用户信息缓存
func (c *RedisCache) HasUser(ctx context.Context, uId uint) (cacheData map[string]string, err error) {
	// 获取一个哈希表中的所有字段和值
	cacheData, err = c.HGetAll(ctx, GenUserInfoCacheKey(uId)).Result()
	if err != nil {
		return nil, err
	}
	// 判断用户是否存在，如果哈希表为空，表示不存在，如果不为空，表示存在
	return cacheData, nil
}
