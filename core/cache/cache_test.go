package cache

import (
	"context"
	"testing"
	"time"
)

func TestMemoryCache_SetAndGet(t *testing.T) {
	c := NewMemoryCache()
	defer c.Close()

	// 测试基本的 Set/Get
	err := c.Set("key1", "value1")
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}

	val, ok := c.Get("key1")
	if !ok {
		t.Error("Get failed: key not found")
	}
	if val != "value1" {
		t.Errorf("Get returned wrong value: got %v, want value1", val)
	}

	// 测试不存在的key
	_, ok = c.Get("nonexistent")
	if ok {
		t.Error("Get should return false for nonexistent key")
	}
}

func TestMemoryCache_Setex(t *testing.T) {
	c := NewMemoryCache()
	defer c.Close()

	// 设置1秒过期
	err := c.Setex("expiring_key", "value", 1)
	if err != nil {
		t.Errorf("Setex failed: %v", err)
	}

	// 立即获取应该成功
	val, ok := c.Get("expiring_key")
	if !ok {
		t.Error("Get failed immediately after Setex")
	}
	if val != "value" {
		t.Errorf("Got wrong value: %v", val)
	}

	// 等待过期
	time.Sleep(1100 * time.Millisecond)

	// 过期后获取应该失败
	_, ok = c.Get("expiring_key")
	if ok {
		t.Error("Get should return false for expired key")
	}
}

func TestMemoryCache_Exists(t *testing.T) {
	c := NewMemoryCache()
	defer c.Close()

	c.Set("existing_key", "value")

	exists, err := c.Exists("existing_key")
	if err != nil {
		t.Errorf("Exists failed: %v", err)
	}
	if !exists {
		t.Error("Exists should return true for existing key")
	}

	exists, err = c.Exists("nonexistent")
	if err != nil {
		t.Errorf("Exists failed: %v", err)
	}
	if exists {
		t.Error("Exists should return false for nonexistent key")
	}
}

func TestMemoryCache_Del(t *testing.T) {
	c := NewMemoryCache()
	defer c.Close()

	c.Set("key1", "value1")
	c.Set("key2", "value2")

	count, err := c.Del("key1", "key2", "nonexistent")
	if err != nil {
		t.Errorf("Del failed: %v", err)
	}
	if count != 2 {
		t.Errorf("Del returned wrong count: got %d, want 2", count)
	}

	_, ok := c.Get("key1")
	if ok {
		t.Error("key1 should be deleted")
	}
}

func TestMemoryCache_TTL(t *testing.T) {
	c := NewMemoryCache()
	defer c.Close()

	// 设置5秒过期
	c.Setex("ttl_key", "value", 5)

	ttl, err := c.TTL("ttl_key")
	if err != nil {
		t.Errorf("TTL failed: %v", err)
	}
	if ttl < 4 || ttl > 5 {
		t.Errorf("TTL returned unexpected value: %d", ttl)
	}

	// 不存在的key
	ttl, _ = c.TTL("nonexistent")
	if ttl != -2 {
		t.Errorf("TTL for nonexistent key should be -2, got %d", ttl)
	}

	// 没有过期时间的key
	c.Set("no_expire", "value")
	ttl, _ = c.TTL("no_expire")
	if ttl != -1 {
		t.Errorf("TTL for key without expiry should be -1, got %d", ttl)
	}
}

func TestMemoryCache_Incr(t *testing.T) {
	c := NewMemoryCache()
	defer c.Close()

	// 不存在的key，从0开始
	val, err := c.Incr("counter")
	if err != nil {
		t.Errorf("Incr failed: %v", err)
	}
	if val != 1 {
		t.Errorf("Incr returned wrong value: got %d, want 1", val)
	}

	// 再次增加
	val, _ = c.Incr("counter")
	if val != 2 {
		t.Errorf("Incr returned wrong value: got %d, want 2", val)
	}

	// IncrBy
	val, _ = c.IncrBy("counter", 10)
	if val != 12 {
		t.Errorf("IncrBy returned wrong value: got %d, want 12", val)
	}
}

func TestMemoryCache_Expire(t *testing.T) {
	c := NewMemoryCache()
	defer c.Close()

	c.Set("key", "value")

	ok, err := c.Expire("key", 1)
	if err != nil {
		t.Errorf("Expire failed: %v", err)
	}
	if !ok {
		t.Error("Expire should return true for existing key")
	}

	// 等待过期
	time.Sleep(1100 * time.Millisecond)

	_, found := c.Get("key")
	if found {
		t.Error("Key should be expired")
	}
}

func TestMemoryCache_Keys(t *testing.T) {
	c := NewMemoryCache()
	defer c.Close()

	c.Set("user:1", "a")
	c.Set("user:2", "b")
	c.Set("order:1", "c")

	keys, err := c.Keys("user:*")
	if err != nil {
		t.Errorf("Keys failed: %v", err)
	}
	if len(keys) != 2 {
		t.Errorf("Keys returned wrong count: got %d, want 2", len(keys))
	}

	keys, _ = c.Keys("*")
	if len(keys) != 3 {
		t.Errorf("Keys('*') returned wrong count: got %d, want 3", len(keys))
	}
}

func TestMemoryCache_Flush(t *testing.T) {
	c := NewMemoryCache()
	defer c.Close()

	c.Set("key1", "value1")
	c.Set("key2", "value2")

	err := c.Flush()
	if err != nil {
		t.Errorf("Flush failed: %v", err)
	}

	if c.Len() != 0 {
		t.Errorf("Cache should be empty after Flush, got %d items", c.Len())
	}
}

func TestMemoryCache_ContextMethods(t *testing.T) {
	c := NewMemoryCache()
	defer c.Close()

	ctx := context.Background()

	// 测试带context的方法
	err := c.SetCtx(ctx, "ctx_key", "ctx_value")
	if err != nil {
		t.Errorf("SetCtx failed: %v", err)
	}

	val, ok := c.GetCtx(ctx, "ctx_key")
	if !ok || val != "ctx_value" {
		t.Error("GetCtx failed")
	}

	exists, _ := c.ExistsCtx(ctx, "ctx_key")
	if !exists {
		t.Error("ExistsCtx failed")
	}

	_, err = c.DelCtx(ctx, "ctx_key")
	if err != nil {
		t.Errorf("DelCtx failed: %v", err)
	}
}

func TestGlobalCache(t *testing.T) {
	// 测试全局缓存函数
	err := Set("global_key", "global_value")
	if err != nil {
		t.Errorf("Global Set failed: %v", err)
	}

	val, ok := Get("global_key")
	if !ok || val != "global_value" {
		t.Error("Global Get failed")
	}

	exists, _ := Exists("global_key")
	if !exists {
		t.Error("Global Exists failed")
	}

	Del("global_key")
	exists, _ = Exists("global_key")
	if exists {
		t.Error("Global Del failed")
	}
}
