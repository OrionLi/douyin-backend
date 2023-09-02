package cache

import (
	"fmt"
	"user-center/pkg/util"
)

/*
用于规范缓存键的格式
*/

const FollowCount = "FollowCount"
const FanCount = "FanCount"

func GenFollowUserCacheKey(userId, followUserId uint) string {
	return "follow:" + util.UintToStr(userId) + ":" + util.UintToStr(followUserId)
}

func GenUserInfoCacheKey(userId uint) string {
	return "user:info:" + util.UintToStr(userId)
}

func VideoCacheCountKey(id int64) string {
	return fmt.Sprintf("publishlist:%d", id)
}
