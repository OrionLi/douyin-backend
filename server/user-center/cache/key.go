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
	return "follow_user_" + util.UintToStr(userId) + "_" + util.UintToStr(followUserId)
}

func GenUserInfoCacheKey(userId uint) string {
	return "user_info_" + util.UintToStr(userId)
}

func GenFollowCountCacheKey(userId uint) string {
	return "count:follow:" + util.UintToStr(userId)
}

func GenFollowerCountCacheKey(userId uint) string {
	return "count:follower:" + util.UintToStr(userId)
}

func VideoCacheCountKey(id int64) string {
	return fmt.Sprintf("publishlist:%d", id)
}
