package service

import (
	"context"
	"github.com/go-redis/redis/v8"
	"user-center/cache"
	"user-center/dao"
	"user-center/model"
)

type RelationService struct {
	dao *dao.RelationDao
}

func NewRelationService(dao *dao.RelationDao) *RelationService {
	return &RelationService{dao: dao}
}

// Follow 关注
func (s *RelationService) Follow(ctx context.Context, userId, toUserId int64) error {
	relationCache := cache.NewRelationCache(ctx)
	relationDao := dao.NewRelationDao(ctx)
	err := relationDao.Follow(userId, toUserId)
	if err != nil {
		return err
	}
	// 更新关注数缓存
	return relationCache.IncrFollowCountCache(userId, 1)
}

// Unfollow 取消关注
func (s *RelationService) Unfollow(ctx context.Context, userId, toUserId int64) error {
	relationDao := dao.NewRelationDao(ctx)
	relationCache := cache.NewRelationCache(ctx)
	err := relationDao.Unfollow(userId, toUserId)
	if err != nil {
		return err
	}
	// 更新关注数缓存
	return relationCache.IncrFollowCountCache(userId, -1)
}

// GetFollowList 获取关注列表
func (s *RelationService) GetFollowList(ctx context.Context, userId int64) ([]*model.User, error) {
	relationCache := cache.NewRelationCache(ctx)
	relationDao := dao.NewRelationDao(ctx)
	// 先从缓存中获取关注数
	followCount, err := relationCache.GetFollowCountCache(userId)
	if err != nil && err != redis.Nil {
		return nil, err
	}

	users, err := relationDao.GetFollowList(userId)
	if err != nil {
		return nil, err
	}

	// 如果缓存不存在,更新缓存
	if followCount == 0 {
		if err = relationCache.IncrFollowCountCache(userId, int64(len(users))); err != nil {
			return nil, err
		}
	}

	return users, nil
}

// GetFollowerList 获取粉丝列表
func (s *RelationService) GetFollowerList(ctx context.Context, userId int64) ([]*model.User, error) {
	relationDao := dao.NewRelationDao(ctx)
	relationCache := cache.NewRelationCache(ctx)

	// 先从缓存中获取粉丝数
	followerCount, err := relationCache.GetFollowerCountCache(userId)
	if err != nil && err != redis.Nil {
		return nil, err
	}

	followers, err := relationDao.GetFollowerList(userId)
	if err != nil {
		return nil, err
	}

	// 如果缓存不存在,更新缓存
	if followerCount == 0 {
		if err = relationCache.IncrFollowerCountCache(userId, int64(len(followers))); err != nil {
			return nil, err
		}
	}

	return followers, nil
}

// GetFriendList 获取好友列表
func (s *RelationService) GetFriendList(ctx context.Context, userId int64) ([]*model.User, error) {
	relationDao := dao.NewRelationDao(ctx)
	relationCache := cache.NewRelationCache(ctx)

	// 先从缓存中获取粉丝数
	followerCount, err := relationCache.GetFollowerCountCache(userId)
	if err != nil && err != redis.Nil {
		return nil, err
	}

	friends, err := relationDao.GetFriendList(userId)
	if err != nil {
		return nil, err
	}

	// 好友列表与粉丝列表大小相同,直接使用粉丝数缓存
	if followerCount == 0 {
		if err = relationCache.IncrFollowerCountCache(userId, int64(len(friends))); err != nil {
			return nil, err
		}
	}

	return friends, nil
}
