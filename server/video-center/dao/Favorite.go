package dao

import (
	"context"
	"gorm.io/gorm"
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
	tx := DB.WithContext(ctx).Begin() // 开启事务
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	fav := Fav{
		UserId:  uint(userId),
		VideoId: videoId,
	}
	err := tx.Model(&Fav{}).Create(&fav).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Model(&Video{}).Where("id = ?", videoId).Update("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// DeleteFav 取消点赞
func DeleteFav(ctx context.Context, videoId int64, userId int64) error {
	tx := DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	err := tx.Model(&Fav{}).Where("user_id = ? and video_id = ?", userId, videoId).Delete(&Fav{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Model(&Video{}).Where("id = ?", videoId).Update("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
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
