package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type Video struct {
	gorm.Model
	Id            int64  `json:"id,omitempty"`             // 视频唯一标识
	AuthorID      int64  `json:"author_id,omitempty"`      // 视频作者userId
	PlayUrl       string `json:"play_url,omitempty"`       // 视频播放地址
	CoverUrl      string `json:"cover_url,omitempty"`      // 视频封面地址
	FavoriteCount int64  `json:"favorite_count,omitempty"` // 视频的点赞总数
	CommentCount  int64  `json:"comment_count,omitempty"`  // 视频的评论总数
	IsFavorite    bool   `json:"is_favorite,omitempty"`    // true-已点赞，false-未点赞
	Title         string `json:"title,omitempty"`          // 视频标题
}

// SaveVideo 保存video，其存储title、PlayURL、CoverURL、AuthorID
func SaveVideo(ctx context.Context, videos []*Video) error {
	return DB.WithContext(ctx).Create(videos).Error
}

// QueryVideoListByTitle 根据Title查询video集合
func QueryVideoListByTitle(ctx context.Context, title string) ([]*Video, error) {
	res := make([]*Video, 0)
	if err := DB.WithContext(ctx).Where("title like ?", "%"+title+"%").Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

// QueryVideoListByAuthorId 查询目标user下的所有video
func QueryVideoListByAuthorId(ctx context.Context, AuthorId int64) ([]*Video, error) {
	res := make([]*Video, 0)
	if err := DB.WithContext(ctx).Where("author_id = ?", AuthorId).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

// FeedByTime 根据传入的时间戳来返回videoList，如果没有，则从当前时间开始
func FeedByTime(ctx context.Context, lastTime int64) ([]*Video, error) {
	res := make([]*Video, 0)
	var curTime time.Time
	if lastTime <= 0 {
		curTime = time.Now()
	} else {
		curTime = time.Unix(lastTime, 0)
	}
	if err := DB.WithContext(ctx).Where("created_time<=", curTime).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}
