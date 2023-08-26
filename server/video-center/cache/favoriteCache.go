package cache

import (
	"context"
	"errors"
	"fmt"
	"time"
)

const FavoriteUpdateSetKey = "fav_update_set"

// ActionFavoriteCache 点赞缓存
// 通过原子性操作解决并发问题
func ActionFavoriteCache(videoId int64, actionType int32) error {
	favoriteKey := fmt.Sprintf("favorite:%d", videoId)
	switch actionType {
	case 1:
		// 使用 RedisClient 执行 INCR 命令
		_, err := RedisClient.Incr(context.Background(), favoriteKey).Result()
		if err != nil {
			return err
		}
	case 2:
		_, err := RedisClient.Decr(context.Background(), favoriteKey).Result()
		if err != nil {
			return err
		}
	default:
		return errors.New("actionType error")
	}
	// 更新集合
	go RedisClient.SAdd(context.Background(), FavoriteUpdateSetKey, videoId)
	return nil
}

// SetFavoriteCountCache 设置缓存中的某个视频点赞数量
func SetFavoriteCountCache(videoId int64, favoriteCount int64) error {
	favoriteKey := fmt.Sprintf("favorite:%d", videoId)
	return RedisClient.Set(context.Background(), favoriteKey, favoriteCount, 7*24*time.Hour).Err()
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

// GetMemberFromFavoriteUpdateSet 获取更新集合中的所有视频ID
func GetMemberFromFavoriteUpdateSet() ([]string, error) {
	return RedisClient.SMembers(context.Background(), FavoriteUpdateSetKey).Result()
}
