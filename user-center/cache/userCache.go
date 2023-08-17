package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisCache struct {
	*redis.Client
}

func NewRedisCache(ctx context.Context) *RedisCache {
	return &RedisCache{NewRedisClient(ctx)}
}
func (userCache *RedisCache) IsFollow(ctx context.Context, uId, followId uint) bool {

	return userCache.
		Exists(ctx, GenFollowUserCacheKey(uId, followId)).
		Val() == 1
}
func (userCache *RedisCache) AddFollow(ctx context.Context, uId, followId uint) error {

	return userCache.
		Set(ctx, GenFollowUserCacheKey(uId, followId), 1, time.Hour).
		Err()
}
