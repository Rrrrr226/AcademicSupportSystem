package cache

import "context"

// 带Context的便捷方法

// SetCtx 设置缓存（带context，无过期时间）
func SetCtx(ctx context.Context, key string, value interface{}) error {
	return GetCache().SetCtx(ctx, key, value)
}

// SetexCtx 设置缓存（带context和过期时间，单位秒）
func SetexCtx(ctx context.Context, key string, value interface{}, expireSeconds int) error {
	return GetCache().SetexCtx(ctx, key, value, expireSeconds)
}

// GetCtx 获取缓存值（带context）
func GetCtx(ctx context.Context, key string) (interface{}, bool) {
	return GetCache().GetCtx(ctx, key)
}

// GetStringCtx 获取字符串类型的缓存值（带context）
func GetStringCtx(ctx context.Context, key string) (string, bool) {
	return GetCache().GetStringCtx(ctx, key)
}

// ExistsCtx 检查key是否存在（带context）
func ExistsCtx(ctx context.Context, key string) (bool, error) {
	return GetCache().ExistsCtx(ctx, key)
}

// DelCtx 删除缓存（带context）
func DelCtx(ctx context.Context, keys ...string) (int, error) {
	return GetCache().DelCtx(ctx, keys...)
}

// ExpireCtx 设置key的过期时间（带context，单位秒）
func ExpireCtx(ctx context.Context, key string, expireSeconds int) (bool, error) {
	return GetCache().ExpireCtx(ctx, key, expireSeconds)
}

// TTLCtx 获取key的剩余过期时间（带context，单位秒）
func TTLCtx(ctx context.Context, key string) (int, error) {
	return GetCache().TTLCtx(ctx, key)
}

// KeysCtx 获取所有匹配pattern的key（带context）
func KeysCtx(ctx context.Context, pattern string) ([]string, error) {
	return GetCache().KeysCtx(ctx, pattern)
}

// IncrCtx 将key中存储的数字增加1（带context）
func IncrCtx(ctx context.Context, key string) (int64, error) {
	return GetCache().IncrCtx(ctx, key)
}

// IncrByCtx 将key中存储的数字增加指定值（带context）
func IncrByCtx(ctx context.Context, key string, increment int64) (int64, error) {
	return GetCache().IncrByCtx(ctx, key, increment)
}

// FlushCtx 清空所有缓存（带context）
func FlushCtx(ctx context.Context) error {
	return GetCache().FlushCtx(ctx)
}
