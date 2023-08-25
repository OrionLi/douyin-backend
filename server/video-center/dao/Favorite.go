package dao

import (
	"context"
	"gorm.io/gorm"
	"video-center/cache"
	"video-center/pkg/util"
)

type Fav struct {
	UserId  uint
	VideoId int64
	Video   Video `gorm:"foreignKey:video_id"`
}

func (u *Fav) TableName() string {
	return "FavRelation"
}

// CreateFav 点赞
func CreateFav(ctx context.Context, videoId int64, userId int64) error {
	return DB.WithContext(ctx).Model(&Fav{}).Create(&Fav{UserId: uint(userId), VideoId: videoId}).Error
}

// DeleteFav 取消点赞
func DeleteFav(ctx context.Context, videoId int64, userId int64) error {
	return DB.WithContext(ctx).Model(&Fav{}).Where("user_id = ? and video_id = ?", userId, videoId).Delete(&Fav{}).Error
}

// IsFavorite 判断是否点赞
func IsFavorite(ctx context.Context, videoId int64, userId int64) (bool, error) {
	var fav Fav
	result := DB.WithContext(ctx).Where("user_id = ? and video_id = ?", userId, videoId).First(&fav)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

// GetFavoriteCount 获取用户点赞数量(给别人的赞)
func GetFavoriteCount(ctx context.Context, userId int64) (int32, error) {
	var favCount int64
	result := DB.WithContext(ctx).Model(&Fav{}).Where("user_id = ?", userId).Count(&favCount)
	if result.Error != nil {
		return 0, result.Error
	}
	return int32(favCount), nil
}

// GetSingleVideoFavoriteCount 获取单个视频点赞数量
func GetSingleVideoFavoriteCount(ctx context.Context, videoId int64) (int32, error) {
	var favCount int32
	result := DB.WithContext(ctx).Table("videos").Where("video_id = ?", videoId).Select("favorite_count").Scan(&favCount)
	if result.Error != nil {
		return 0, result.Error
	}
	return favCount, nil
}

// UpdateMySQLFavoriteCount 异步更新mysql中的值
func UpdateMySQLFavoriteCount(videoID int64, favoriteCount int64) {
	err := DB.Model(&Video{}).Where("id = ?", videoID).Update("favorite_count", favoriteCount).Error
	if err != nil {
		util.LogrusObj.Error("<Favorite Count Update failed> ", "videoId:", videoID, "err:", err)
	}
	err = cache.DeleteVideoIdFromFavoriteUpdateSet(videoID)
	if err != nil {
		util.LogrusObj.Error("<Favorite Count Update failed> : Failed to delete video id in Redis", "videoId:", videoID, "err:", err)
	}
}

// ListFav 获取用户喜欢列表
func ListFav(ctx context.Context, userId int64) []Video {
	var favs []Fav
	DB.WithContext(ctx).Where("user_id = ? ", userId).Find(&favs)
	var videoIDs []int64
	for _, rel := range favs {
		videoIDs = append(videoIDs, rel.VideoId)
	}
	if len(videoIDs) == 0 {
		return []Video{}
	}
	var videos []Video
	DB.WithContext(ctx).Find(&videos, videoIDs)
	return videos
}
