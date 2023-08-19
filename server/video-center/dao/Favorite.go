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

func IsFavorite(ctx context.Context, videoId int64, userId int64) (bool, error) {
	var fav Fav
	result := DB.WithContext(ctx).Where("user_id = ? and video_id = ?", userId, videoId).First(&fav)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}
