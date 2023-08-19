package dao

import (
	"context"
	"gorm.io/gorm"
	"math/rand"
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
func TotalCount() (int64, error) {
	var count int64
	tx := DB.Model(&Video{}).Count(&count)
	if tx.Error != nil {
		return 0, nil
	}
	return count, nil
}
func QueryVideosByCurrentTime(ctx context.Context, lastTime int64, page int64, pagesize int64) ([]*Video, error) {
	res := make([]*Video, 0)

	var curTime time.Time
	if lastTime <= 0 {
		curTime = time.Now()
	} else {
		curTime = time.Unix(lastTime, 0)
	}

	count, err := TotalCount()
	if err != nil {
		return nil, err
	}

	// 计算当前偏移量
	curOffset := int(page) * int(pagesize)

	// 计算剩余数据量
	remainingData := count - int64(curOffset)

	// 创建随机生成器
	seed := curTime.UnixNano()
	rng := rand.New(rand.NewSource(seed))

	// 如果剩余数据量不足30条，将偏移量设置为0，返回所有数据
	if remainingData < pagesize {
		curOffset = 0
	} else {
		// 如果剩余数据量足够多，生成随机偏移量
		maxOffset := remainingData - pagesize
		if maxOffset > 0 {
			newPage := rng.Intn(int(maxOffset))
			curOffset += newPage
		}
	}

	tx := DB.WithContext(ctx).Order("created_at DESC").Where("created_at <= ?", curTime).Offset(curOffset).Limit(int(pagesize)).Find(&res)
	if tx.Error != nil {
		return nil, err
	}
	return res, nil
}
