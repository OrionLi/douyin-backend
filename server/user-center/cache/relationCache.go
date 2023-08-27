// cache.go

package cache

import "time"

type Cache struct {
	data map[string]interface{}
	ttl  map[string]time.Time
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]interface{}),
		ttl:  make(map[string]time.Time),
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	// 实现同上
}

func (c *Cache) Set(key string, val interface{}, expire time.Duration) {
	// 实现同上
}

// 其他方法
