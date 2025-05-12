package cache

import (
	"app/internal/tools/safemap"
	"time"
)

type cacheValue struct {
	Data      interface{}
	CreatedAt time.Time
	ExpiredAt time.Time
}

var cacheMap = safemap.New[cacheValue]()

// Remember 从缓存中获取数据，如果不存在或已过期则从数据库中获取并缓存
func Remember[V any](key string, ttl int, callback func() (*V, error)) (*V, error) {
	// 尝试从缓存获取
	if val, ok := cacheMap.Get(key); ok {
		if time.Now().Before(val.ExpiredAt) {
			if data, ok := val.Data.(*V); ok {
				return data, nil
			}
			// 类型不匹配，删除无效缓存
			cacheMap.Delete(key)
		}
	}

	// 执行回调获取新数据
	data, err := callback()
	if err != nil {
		return nil, err
	}

	// 设置缓存时间
	now := time.Now()
	cacheMap.Set(key, cacheValue{
		Data:      data,
		CreatedAt: now,
		ExpiredAt: now.Add(time.Duration(ttl) * time.Second),
	})

	return data, nil
}

// GC 执行垃圾回收，清理过期缓存
func GC() {
	now := time.Now()
	for _, key := range cacheMap.Keys() {
		if val, ok := cacheMap.Get(key); ok {
			if now.After(val.ExpiredAt) {
				cacheMap.Delete(key)
			}
		}
	}
}
