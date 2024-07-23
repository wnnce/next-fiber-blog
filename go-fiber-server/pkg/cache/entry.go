package cache

import (
	"encoding/binary"
	"fmt"
	"github.com/bytedance/sonic"
	"slices"
	"time"
)

var defaultCacheManage *cacheManage

type cacheManage struct {
	mcache    *memoryCache
	marshal   func(any) ([]byte, error)
	unmarshal func([]byte, any) error
}

func init() {
	defaultCacheManage = &cacheManage{
		mcache:    newMemoryCache(8),
		marshal:   sonic.Marshal,
		unmarshal: sonic.Unmarshal,
	}
}

// Set 添加缓存 会将待缓存的数据序列化为二进制数据并拼接上过期时间的二进制数据
// key 缓存的key
// value 待缓存的数据 任意数据
// expire 缓存的过期时间
func Set(key string, value any, expire time.Duration) error {
	bytes, err := defaultCacheManage.marshal(value)
	if err != nil {
		return err
	}
	expireTimestamp := expireTimestampBytes(&expire)
	finalBytes := slices.Concat(expireTimestamp, bytes)
	return defaultCacheManage.mcache.Set(key, finalBytes)
}

func Get[T any](key string) (result T, err error) {
	bytes, err := defaultCacheManage.mcache.Get(key)
	if err != nil {
		return
	}
	timestamp := int64(binary.BigEndian.Uint64(bytes[:8]))
	if timestamp <= time.Now().UnixMilli() {
		defaultCacheManage.mcache.Remove(key)
		return result, fmt.Errorf("cache expiration")
	}
	err = defaultCacheManage.unmarshal(bytes[8:], &result)
	return
}

func GetRaw(key string) ([]byte, error) {
	bytes, err := defaultCacheManage.mcache.Get(key)
	if err != nil {
		return nil, err
	}
	timestamp := int64(binary.BigEndian.Uint64(bytes[:8]))
	if timestamp <= time.Now().UnixMilli() {
		defaultCacheManage.mcache.Remove(key)
		return nil, fmt.Errorf("cache expiration")
	}
	value := make([]byte, 0, len(bytes)-8)
	copy(value, bytes)
	return value, nil
}

func Delete(key string) {
	defaultCacheManage.mcache.Remove(key)
}

func expireTimestampBytes(expire *time.Duration) []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, uint64(time.Now().Add(*expire).UnixMilli()))
	return bytes
}
