package cache

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
	"go-fiber-ent-web-layout/internal/data"
	"time"
)

// RedisOptional Redis操作类
type RedisOptional struct {
	Rc *redis.Client
}

func NewRedisOptional(data *data.Data) *RedisOptional {
	return &RedisOptional{
		Rc: data.Rc,
	}
}

// Set 添加
func (r *RedisOptional) Set(ctx context.Context, key string, value any, expire time.Duration) error {
	_, err := r.Rc.Set(ctx, key, value, expire).Result()
	return err
}

// GetRaw 获取原始值
func (r *RedisOptional) GetRaw(ctx context.Context, key string) (string, error) {
	return r.Rc.Get(ctx, key).Result()
}

// Get 获取序列化之后的值，不支持泛型方法 使用指针
func (r *RedisOptional) Get(ctx context.Context, key string, value any) error {
	result, err := r.Rc.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return sonic.UnmarshalString(result, value)
}

// Remove 删除
func (r *RedisOptional) Remove(ctx context.Context, keys ...string) error {
	_, err := r.Rc.Del(ctx, keys...).Result()
	return err
}

// Client 获取原始Client对象
func (r *RedisOptional) Client() *redis.Client {
	return r.Rc
}
