package cache

import "strconv"

/*
用于规范key的格式
*/
func UintToStr(i uint) string {
	return strconv.FormatInt(int64(i), 10)
}
func GenFollowUserCacheKey(userId, followUserId uint) string {
	return "follow_user_" + UintToStr(userId) + "_" + UintToStr(followUserId)
}
