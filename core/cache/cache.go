package cache

import (
	"context"
	"sync"
	"time"
)

// item 缓存项
type item struct {
	value     interface{}
	expireAt  time.Time
	hasExpire bool
}

// isExpired 检查是否过期
func (i *item) isExpired() bool {
	if !i.hasExpire {
		return false
	}
	return time.Now().After(i.expireAt)
}

// MemoryCache 基于内存的缓存实现
type MemoryCache struct {
	data    map[string]*item
	mu      sync.RWMutex
	cleaner *time.Ticker
	stopCh  chan struct{}
}

// NewMemoryCache 创建新的内存缓存实例
func NewMemoryCache() *MemoryCache {
	c := &MemoryCache{
		data:   make(map[string]*item),
		stopCh: make(chan struct{}),
	}
	// 启动定期清理过期key的goroutine
	c.startCleaner(time.Minute)
	return c
}

// startCleaner 启动清理器
func (c *MemoryCache) startCleaner(interval time.Duration) {
	c.cleaner = time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-c.cleaner.C:
				c.deleteExpired()
			case <-c.stopCh:
				c.cleaner.Stop()
				return
			}
		}
	}()
}

// deleteExpired 删除过期的key
func (c *MemoryCache) deleteExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, v := range c.data {
		if v.isExpired() {
			delete(c.data, k)
		}
	}
}

// Close 关闭缓存，停止清理器
func (c *MemoryCache) Close() {
	close(c.stopCh)
}

// Set 设置缓存（无过期时间）
func (c *MemoryCache) Set(key string, value interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = &item{
		value:     value,
		hasExpire: false,
	}
	return nil
}

// SetCtx 设置缓存（带context，无过期时间）
func (c *MemoryCache) SetCtx(ctx context.Context, key string, value interface{}) error {
	return c.Set(key, value)
}

// Setex 设置缓存（带过期时间，单位秒）
func (c *MemoryCache) Setex(key string, value interface{}, expireSeconds int) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = &item{
		value:     value,
		expireAt:  time.Now().Add(time.Duration(expireSeconds) * time.Second),
		hasExpire: true,
	}
	return nil
}

// SetexCtx 设置缓存（带context和过期时间，单位秒）
func (c *MemoryCache) SetexCtx(ctx context.Context, key string, value interface{}, expireSeconds int) error {
	return c.Setex(key, value, expireSeconds)
}

// Get 获取缓存值
func (c *MemoryCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	it, ok := c.data[key]
	if !ok {
		return nil, false
	}
	if it.isExpired() {
		return nil, false
	}
	return it.value, true
}

// GetCtx 获取缓存值（带context）
func (c *MemoryCache) GetCtx(ctx context.Context, key string) (interface{}, bool) {
	return c.Get(key)
}

// GetString 获取字符串类型的缓存值
func (c *MemoryCache) GetString(key string) (string, bool) {
	val, ok := c.Get(key)
	if !ok {
		return "", false
	}
	str, ok := val.(string)
	return str, ok
}

// GetStringCtx 获取字符串类型的缓存值（带context）
func (c *MemoryCache) GetStringCtx(ctx context.Context, key string) (string, bool) {
	return c.GetString(key)
}

// Exists 检查key是否存在
func (c *MemoryCache) Exists(key string) (bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	it, ok := c.data[key]
	if !ok {
		return false, nil
	}
	if it.isExpired() {
		return false, nil
	}
	return true, nil
}

// ExistsCtx 检查key是否存在（带context）
func (c *MemoryCache) ExistsCtx(ctx context.Context, key string) (bool, error) {
	return c.Exists(key)
}

// Del 删除缓存
func (c *MemoryCache) Del(keys ...string) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	count := 0
	for _, key := range keys {
		if _, ok := c.data[key]; ok {
			delete(c.data, key)
			count++
		}
	}
	return count, nil
}

// DelCtx 删除缓存（带context）
func (c *MemoryCache) DelCtx(ctx context.Context, keys ...string) (int, error) {
	return c.Del(keys...)
}

// Expire 设置key的过期时间（单位秒）
func (c *MemoryCache) Expire(key string, expireSeconds int) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	it, ok := c.data[key]
	if !ok {
		return false, nil
	}
	it.expireAt = time.Now().Add(time.Duration(expireSeconds) * time.Second)
	it.hasExpire = true
	return true, nil
}

// ExpireCtx 设置key的过期时间（带context，单位秒）
func (c *MemoryCache) ExpireCtx(ctx context.Context, key string, expireSeconds int) (bool, error) {
	return c.Expire(key, expireSeconds)
}

// TTL 获取key的剩余过期时间（单位秒）
// 返回 -1 表示key没有设置过期时间
// 返回 -2 表示key不存在
func (c *MemoryCache) TTL(key string) (int, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	it, ok := c.data[key]
	if !ok {
		return -2, nil
	}
	if !it.hasExpire {
		return -1, nil
	}
	if it.isExpired() {
		return -2, nil
	}
	ttl := int(time.Until(it.expireAt).Seconds())
	if ttl < 0 {
		return 0, nil
	}
	return ttl, nil
}

// TTLCtx 获取key的剩余过期时间（带context，单位秒）
func (c *MemoryCache) TTLCtx(ctx context.Context, key string) (int, error) {
	return c.TTL(key)
}

// Keys 获取所有匹配pattern的key（简单实现，仅支持*通配符）
func (c *MemoryCache) Keys(pattern string) ([]string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	var keys []string
	for k, v := range c.data {
		if !v.isExpired() {
			if pattern == "*" || matchPattern(pattern, k) {
				keys = append(keys, k)
			}
		}
	}
	return keys, nil
}

// KeysCtx 获取所有匹配pattern的key（带context）
func (c *MemoryCache) KeysCtx(ctx context.Context, pattern string) ([]string, error) {
	return c.Keys(pattern)
}

// Incr 将key中存储的数字增加1
func (c *MemoryCache) Incr(key string) (int64, error) {
	return c.IncrBy(key, 1)
}

// IncrCtx 将key中存储的数字增加1（带context）
func (c *MemoryCache) IncrCtx(ctx context.Context, key string) (int64, error) {
	return c.Incr(key)
}

// IncrBy 将key中存储的数字增加指定值
func (c *MemoryCache) IncrBy(key string, increment int64) (int64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	it, ok := c.data[key]
	if !ok {
		c.data[key] = &item{
			value:     increment,
			hasExpire: false,
		}
		return increment, nil
	}

	if it.isExpired() {
		c.data[key] = &item{
			value:     increment,
			hasExpire: false,
		}
		return increment, nil
	}

	var currentVal int64
	switch v := it.value.(type) {
	case int64:
		currentVal = v
	case int:
		currentVal = int64(v)
	case int32:
		currentVal = int64(v)
	default:
		currentVal = 0
	}

	newVal := currentVal + increment
	it.value = newVal
	return newVal, nil
}

// IncrByCtx 将key中存储的数字增加指定值（带context）
func (c *MemoryCache) IncrByCtx(ctx context.Context, key string, increment int64) (int64, error) {
	return c.IncrBy(key, increment)
}

// Flush 清空所有缓存
func (c *MemoryCache) Flush() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[string]*item)
	return nil
}

// FlushCtx 清空所有缓存（带context）
func (c *MemoryCache) FlushCtx(ctx context.Context) error {
	return c.Flush()
}

// Len 获取缓存中key的数量
func (c *MemoryCache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	count := 0
	for _, v := range c.data {
		if !v.isExpired() {
			count++
		}
	}
	return count
}

// matchPattern 简单的模式匹配（仅支持 * 通配符）
func matchPattern(pattern, key string) bool {
	if pattern == "*" {
		return true
	}
	// 简单实现：前缀匹配 "prefix*"
	if len(pattern) > 0 && pattern[len(pattern)-1] == '*' {
		prefix := pattern[:len(pattern)-1]
		return len(key) >= len(prefix) && key[:len(prefix)] == prefix
	}
	return pattern == key
}
