package timeout

import (
	"context"
	"github.com/gofiber/fiber/v3"
	"go-fiber-ent-web-layout/pkg/pool"
	"time"
)

// NewMiddleware 返回请求超时中间件
// 向请求ctx中set一个WithTimeout的Context
func NewMiddleware(timeout time.Duration) fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.UserContext(), timeout)
		defer cancel()

		c.SetUserContext(ctx)

		ch := make(chan struct{})

		var err error
		pool.Go(func() {
			err = c.Next()
			ch <- struct{}{}
		})

		select {
		// 如果请求正常完成那么直接返回
		case <-ch:
			return err
		// 返回请求超时错误
		case <-ctx.Done():
			return fiber.ErrRequestTimeout
		}
	}
}
