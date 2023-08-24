package cache

import (
	"fmt"
	"user-center/pkg/util"
)

/*
用于规范缓存键的格式
*/

func GenFollowUserCacheKey(userId, followUserId uint) string {
	return "follow_user_" + util.UintToStr(userId) + "_" + util.UintToStr(followUserId)
}
func GenUserInfoCacheKey(userId uint) string {
	return "user_info_" + util.UintToStr(userId)
}

func CacheChangeUserCountKey(userID int64, category string) string {
	return fmt.Sprintf("user:%d:%s", userID, category)
}
