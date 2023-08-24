package cache

import (
	"context"
	"time"
)

func NewUserCache(ctx context.Context) *RedisCache {
	return &RedisCache{NewRedisClient(ctx)}
}
func (userCache *RedisCache) IsFollow(ctx context.Context, uId, followId uint) bool {

	return userCache.
		Exists(ctx, GenFollowUserCacheKey(uId, followId)).
		Val() == 1
}

// AddFollow 关注关系缓存
func (userCache *RedisCache) AddFollow(ctx context.Context, uId, followId uint) error {

	return userCache.
		Set(ctx, GenFollowUserCacheKey(uId, followId), 1, time.Hour).
		Err()
}

// AddUser 用户信息缓存
func (userCache *RedisCache) AddUser(ctx context.Context, uId uint, m map[string]interface{}) error {
	err := userCache.HSet(ctx, GenUserInfoCacheKey(uId), m).Err()
	if err != nil {
		return err
	}
	// 设置键的过期时间为10天
	return userCache.Expire(ctx, GenUserInfoCacheKey(uId), 10*24*3600*time.Second).Err()
}

// HasUser 判断Redis中是否存在某个用户信息缓存
func (userCache *RedisCache) HasUser(ctx context.Context, uId uint) (cacheData map[string]string, err error) {
	// 获取一个哈希表中的所有字段和值
	cacheData, err = userCache.HGetAll(ctx, GenUserInfoCacheKey(uId)).Result()
	if err != nil {
		return nil, err
	}
	// 判断用户是否存在，如果哈希表为空，表示不存在，如果不为空，表示存在
	return cacheData, nil
}
