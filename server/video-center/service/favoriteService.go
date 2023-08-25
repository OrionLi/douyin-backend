package service

import (
	"context"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"github.com/go-redis/redis/v8"
	"time"
	"video-center/cache"
	"video-center/dao"
	"video-center/pkg/util"
)

type FavoriteService interface {
	CreateFav(videoId int64, userId int64) error
	DeleteFav(videoId int64, userId int64) error
	IsFav(videoId int64, userId int64) (bool, error)
	ListFav(userId int64) (bool, []*pb.Video)
	CountFav(userId int64) (int32, int32, error)
}

type FavoriteServiceImpl struct {
	ctx context.Context
}

// CreateFav 创建点赞记录
// TODO 考虑当redis宕机时，如何保证数据一致性
func (f FavoriteServiceImpl) CreateFav(videoId int64, userId int64) error {
	// TODO 验证是否已经点赞 此项应在压测通过后实现
	// HACK IsFav() 验证 重复点赞问题
	err := dao.CreateFav(f.ctx, videoId, userId)
	if err != nil {
		return err
	}
	return cache.ActionFavoriteCache(videoId, 1)
}

// DeleteFav 删除点赞记录
func (f FavoriteServiceImpl) DeleteFav(videoId int64, userId int64) error {
	// TODO 验证是否未点赞 此项应在压测通过后实现
	// HACK IsFav() 验证 未点赞试图取消点赞问题
	err := dao.DeleteFav(f.ctx, videoId, userId)
	if err != nil {
		return err
	}
	return cache.ActionFavoriteCache(videoId, 2)
}

// IsFav 判断是否点赞
func (f FavoriteServiceImpl) IsFav(videoId int64, userId int64) (bool, error) {
	return dao.IsFavorite(f.ctx, videoId, userId)
}

// ListFav 获取喜欢列表
func (f FavoriteServiceImpl) ListFav(userId int64) (bool, []*pb.Video) {
	//得到喜欢的视频集合
	favs := dao.ListFav(f.ctx, userId)
	if len(favs) == 0 {
		return false, []*pb.Video{}
	}
	//pb
	var favVideoList []*pb.Video
	for _, v := range favs {
		//todo 得到用户ID，然后调用rpc查询用户信息
		//var user := xxx(v.AuthorID) //然后修改Author
		video := &pb.Video{
			Id:            v.Id,
			Author:        nil,
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    true,
			Title:         v.Title,
		}
		favVideoList = append(favVideoList, video)
	}
	return true, favVideoList
}

// CountFav 获取用户点赞数量, 收到的点赞数量
// TODO 是否考虑缓存穿透问题
func (f FavoriteServiceImpl) CountFav(userId int64) (int32, int32, error) {
	favoriteCount, err := dao.GetFavoriteCount(f.ctx, userId)
	if err != nil {
		return 0, 0, err
	}
	VideoList, err := dao.QueryVideoListByAuthorId(f.ctx, userId)
	if err != nil {
		return 0, 0, err
	}
	var receivedFavoriteCount int64
	for _, video := range VideoList {
		count, err := cache.GetFavoriteCountCache(video.Id)
		if err != nil {
			if err == redis.Nil {
				count = video.FavoriteCount
				err := cache.SetFavoriteCountCache(video.Id, count)
				if err != nil {
					util.LogrusObj.Error("<Favorite Count In Redis Update failed> ", "videoId:", video.Id, "err:", err)
				}
			} else {
				return 0, 0, err
			}
		}
		receivedFavoriteCount += count
	}
	return favoriteCount, int32(receivedFavoriteCount), nil
}

// UpdateMySQLFavoriteCount 更新到MySQL
func UpdateMySQLFavoriteCount(videoId int64, favoriteCount int64) {
	err := dao.UpdateFavoriteCountByVideoId(videoId, favoriteCount)
	if err != nil {
		util.LogrusObj.Error("<Favorite Count Update failed> ", "videoId:", videoId, "err:", err)
	}
	err = cache.DeleteVideoIdFromFavoriteUpdateSet(videoId)
	if err != nil {
		util.LogrusObj.Error("<Favorite Count Update failed> : Failed to delete video id in Redis", "videoId:", videoId, "err:", err)
	}
}

func UpdateFavoriteCacheToMySQL() {
	favoriteUpdateSet, err := cache.GetMemberFromFavoriteUpdateSet()
	if err != nil {
		util.LogrusObj.Error("<Favorite Count Update failed>", ": Get list fail", err)
	}
	// 处理每个视频ID
	for _, videoIdStr := range favoriteUpdateSet {
		videoId := util.StringToInt64(videoIdStr)
		count, err := cache.GetFavoriteCountCache(videoId)
		if err != nil {
			util.LogrusObj.Error("<Favorite Count Update failed> ", "videoId:", videoId, "err:", err)
		}
		go UpdateMySQLFavoriteCount(videoId, count)
	}
}

// UpdateFavoriteCacheToMySQLAtRegularTime 更新到MySQL
func UpdateFavoriteCacheToMySQLAtRegularTime() {
	util.LogrusObj.Info("goroutine:UpdateToMySQL is running", time.Now())
	interval := 12 * time.Hour // 设置定时任务的时间间隔
	// 创建一个定时器
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			util.LogrusObj.Info("UpdateToMySQL task start at", time.Now())
			UpdateFavoriteCacheToMySQL()
			util.LogrusObj.Info("UpdateToMySQL task end at", time.Now())
		}
	}
}

func NewFavoriteService(context context.Context) FavoriteService {
	return &FavoriteServiceImpl{
		ctx: context,
	}
}
