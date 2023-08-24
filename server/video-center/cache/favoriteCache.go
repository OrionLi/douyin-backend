package cache

func ActionFavoriteCache(videoId int64, actionType int32) error {
	// TODO 通过videoId创建键值对： key为favorite:videoId value为favorite_count
	// TODO 当actionType为1时，value+1；当actionType为2时，value-1
	// TODO 操作该值时，需要加锁，防止并发操作
	// TODO 如果查询key为空，需要查询mysql中的值，然后设置到redis中
	// TODO 如果查询key不为空，直接操作redis中的值
	// TODO 操作完成后，需要解锁
	// TODO 数据一致性策略： 异步更新mysql中的值
	return nil
}

func RedisLock() error {
	// TODO 加锁
	return nil
}

func RedisUnlock() error {
	// TODO 解锁
	return nil
}
