package data

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
	"go-fiber-ent-web-layout/internal/tools/pool"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
)

type RedisTemplate struct {
	rdb *redis.Client
}

func NewRedisTemplate(data *Data) *RedisTemplate {
	return &RedisTemplate{
		rdb: data.Rc,
	}
}

// RedisGetStruct 查询在Redis中缓存的结构体
// 使用泛型指定返回结构体类型
func RedisGetStruct[T any](ctx context.Context, key string, client *redis.Client) (T, error) {
	value := new(T)
	result, err := client.Get(ctx, key).Result()
	if err != nil {
		return *value, err
	}
	err = sonic.UnmarshalString(result, value)
	return *value, err
}

// RedisGetSlice 查询在redis中缓存的切片
func RedisGetSlice[T any](ctx context.Context, key string, client *redis.Client) ([]T, error) {
	result, err := client.Get(ctx, key).Result()
	if err != nil {
		slog.Error("查询Redis缓存结构体失败", "error", err.Error(), "key", key)
		return nil, err
	}
	value := make([]T, 0)
	err = sonic.UnmarshalString(result, &value)
	return value, err
}

func (self *RedisTemplate) Set(ctx context.Context, key string, value any, expireTime time.Duration) error {
	bytes, err := sonic.Marshal(value)
	if err != nil {
		slog.Error("查询Redis缓存切片失败", "error", err.Error(), "key", key)
		return err
	}
	_, err = self.rdb.Set(ctx, key, bytes, expireTime).Result()
	return err
}

func (self *RedisTemplate) Get(ctx context.Context, key string) (string, error) {
	return self.rdb.Get(ctx, key).Result()
}

func (self *RedisTemplate) Delete(ctx context.Context, key string) error {
	_, err := self.rdb.Del(ctx, key).Result()
	return err
}

func (self *RedisTemplate) PatternDelete(ctx context.Context, pattern string) (int64, error) {
	keys, err := self.rdb.Keys(ctx, pattern).Result()
	if err != nil {
		return 0, err
	}
	var count int64
	var wg sync.WaitGroup
	for _, key := range keys {
		wg.Add(1)
		pool.Go(func() {
			defer wg.Done()
			_, err = self.rdb.Del(ctx, key).Result()
			if err != nil {
				slog.Error("批量删除Redis Key失败", "error", err.Error(), "pattern", pattern, "key", key)
				return
			}
			atomic.AddInt64(&count, 1)
		})
	}
	wg.Wait()
	return count, nil
}

func (self *RedisTemplate) ZSetAdd(ctx context.Context, key string, value any, scope float64) error {
	_, err := self.rdb.ZAdd(ctx, key, redis.Z{Member: value, Score: scope}).Result()
	return err
}

func (self *RedisTemplate) ZSetScope(ctx context.Context, key, value string) (float64, error) {
	return self.rdb.ZScore(ctx, key, value).Result()
}

func (self *RedisTemplate) ZIncrBy(ctx context.Context, key, value string, scope float64) error {
	_, err := self.rdb.ZIncrBy(ctx, key, scope, value).Result()
	return err
}

func (self *RedisTemplate) Client() *redis.Client {
	return self.rdb
}
