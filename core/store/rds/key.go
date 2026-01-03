package rds

import "strings"

const (
	// DefaultKeyPrefix 默认的key前缀
	DefaultKeyPrefix = "app"
	// KeySeparator key分隔符
	KeySeparator = ":"
)

var keyPrefix = DefaultKeyPrefix

// SetKeyPrefix 设置全局key前缀
func SetKeyPrefix(prefix string) {
	keyPrefix = prefix
}

// GetKeyPrefix 获取当前key前缀
func GetKeyPrefix() string {
	return keyPrefix
}

// Key 生成带前缀的缓存key
// 例如: Key("oauth", "mark", "abc123") -> "app:oauth:mark:abc123"
func Key(parts ...string) string {
	if len(parts) == 0 {
		return keyPrefix
	}
	allParts := make([]string, 0, len(parts)+1)
	allParts = append(allParts, keyPrefix)
	allParts = append(allParts, parts...)
	return strings.Join(allParts, KeySeparator)
}

// KeyWithoutPrefix 生成不带前缀的缓存key
// 例如: KeyWithoutPrefix("oauth", "mark", "abc123") -> "oauth:mark:abc123"
func KeyWithoutPrefix(parts ...string) string {
	return strings.Join(parts, KeySeparator)
}

// KeyWithPrefix 使用指定前缀生成缓存key
// 例如: KeyWithPrefix("custom", "oauth", "mark") -> "custom:oauth:mark"
func KeyWithPrefix(prefix string, parts ...string) string {
	if len(parts) == 0 {
		return prefix
	}
	allParts := make([]string, 0, len(parts)+1)
	allParts = append(allParts, prefix)
	allParts = append(allParts, parts...)
	return strings.Join(allParts, KeySeparator)
}
