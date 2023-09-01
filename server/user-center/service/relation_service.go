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
	_, err = GetFollowCount(relationDao, relationCache, userId)
	return err
}

// Unfollow 取消关注
func (s *RelationService) Unfollow(ctx context.Context, userId, toUserId int64) error {
	relationDao := dao.NewRelationDao(ctx)
	relationCache := cache.NewRelationCache(ctx)
	err := relationDao.Unfollow(userId, toUserId)
	if err != nil {
		return err
	}
	_, err = GetFollowCount(relationDao, relationCache, userId)
	return err
}

// GetFollowList 获取关注列表
func (s *RelationService) GetFollowList(ctx context.Context, userId int64) ([]*model.User, error) {
	relationDao := dao.NewRelationDao(ctx)
	users, err := relationDao.GetFollowList(userId)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// GetFollowerList 获取粉丝列表
func (s *RelationService) GetFollowerList(ctx context.Context, userId int64) ([]*model.User, error) {
	relationDao := dao.NewRelationDao(ctx)
	followers, err := relationDao.GetFollowerList(userId)
	if err != nil {
		return nil, err
	}
	return followers, nil
}

// GetFriendList 获取好友列表
func (s *RelationService) GetFriendList(ctx context.Context, userId int64) ([]*model.User, error) {
	relationDao := dao.NewRelationDao(ctx)
	//relationCache := cache.NewRelationCache(ctx)

	// 先从缓存中获取粉丝数
	//followerCount, err := relationCache.GetFollowerCountCache(userId)
	//if err != nil && err != redis.Nil {
	//	return nil, err
	//}

	friends, err := relationDao.GetFriendList(userId)
	if err != nil {
		return nil, err
	}

	//// 好友列表与粉丝列表大小相同,直接使用粉丝数缓存
	//if followerCount == 0 {
	//	if err = relationCache.UpdateFollowerCountCache(userId, int64(len(friends))); err != nil {
	//		return nil, err
	//	}
	//}

	return friends, nil
}

// GetFollowerCount 获取粉丝数
func GetFollowerCount(relationDao *dao.RelationDao, relationCache *cache.RedisCache, userId int64) (int64, error) {
	// 先从缓存中获取粉丝数
	followerCount, err := relationCache.GetFollowerCountCache(userId)
	// 读到了
	if err != redis.Nil {
		return followerCount, nil
	}
	// 没读到，顺带更新user的count字段
	followers, err := relationDao.GetFollowerList(userId)
	followerCount = int64(len(followers))
	if err = relationDao.UpdateFanCount(userId, int(followerCount)); err != nil {
		return 0, err
	}
	if err = relationCache.UpdateFollowerCountCache(userId, followerCount); err != nil {
		return 0, err
	}
	return followerCount, nil
}

// GetFollowCount 获取关注数
func GetFollowCount(relationDao *dao.RelationDao, relationCache *cache.RedisCache, userId int64) (int64, error) {
	// 先从缓存中获取粉丝数
	followCount, err := relationCache.GetFollowCountCache(userId)
	// 读到了
	if err != redis.Nil {
		return followCount, nil
	}
	// 没读到，顺带更新user的count字段
	followers, err := relationDao.GetFollowList(userId)
	followCount = int64(len(followers))
	if err = relationDao.UpdateFollowCount(userId, int(followCount)); err != nil {
		return 0, err
	}
	if err = relationCache.UpdateFollowCountCache(userId, followCount); err != nil {
		return 0, err
	}
	return followCount, nil
}
