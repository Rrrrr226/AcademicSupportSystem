package cache

import "sync"

var (
	// 全局缓存实例
	globalCache *MemoryCache
	once        sync.Once
)

// GetCache 获取全局缓存实例（单例模式）
func GetCache() *MemoryCache {
	once.Do(func() {
		globalCache = NewMemoryCache()
	})
	return globalCache
}

// 以下为便捷方法，直接使用全局缓存实例

// Set 设置缓存（无过期时间）
func Set(key string, value interface{}) error {
	return GetCache().Set(key, value)
}

// Setex 设置缓存（带过期时间，单位秒）
func Setex(key string, value interface{}, expireSeconds int) error {
	return GetCache().Setex(key, value, expireSeconds)
}

// Get 获取缓存值
func Get(key string) (interface{}, bool) {
	return GetCache().Get(key)
}

// GetString 获取字符串类型的缓存值
func GetString(key string) (string, bool) {
	return GetCache().GetString(key)
}

// Exists 检查key是否存在
func Exists(key string) (bool, error) {
	return GetCache().Exists(key)
}

// Del 删除缓存
func Del(keys ...string) (int, error) {
	return GetCache().Del(keys...)
}

// Expire 设置key的过期时间（单位秒）
func Expire(key string, expireSeconds int) (bool, error) {
	return GetCache().Expire(key, expireSeconds)
}

// TTL 获取key的剩余过期时间（单位秒）
func TTL(key string) (int, error) {
	return GetCache().TTL(key)
}

// Keys 获取所有匹配pattern的key
func Keys(pattern string) ([]string, error) {
	return GetCache().Keys(pattern)
}

// Incr 将key中存储的数字增加1
func Incr(key string) (int64, error) {
	return GetCache().Incr(key)
}

// IncrBy 将key中存储的数字增加指定值
func IncrBy(key string, increment int64) (int64, error) {
	return GetCache().IncrBy(key, increment)
}

// Flush 清空所有缓存
func Flush() error {
	return GetCache().Flush()
}

// Len 获取缓存中key的数量
func Len() int {
	return GetCache().Len()
}
