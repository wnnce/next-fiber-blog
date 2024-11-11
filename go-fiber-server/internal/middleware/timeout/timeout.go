package timeout

import (
	"context"
	"github.com/gofiber/fiber/v3"
	"go-fiber-ent-web-layout/internal/tools"
	"go-fiber-ent-web-layout/pkg/pool"
	"log/slog"
	"time"
)

// NewMiddleware 返回请求超时中间件
// 向请求ctx中set一个WithTimeout的Context
func NewMiddleware(timeout time.Duration) fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.UserContext(), timeout)
		defer cancel()

		c.SetUserContext(ctx)

		ch := make(chan error)
		defer close(ch)

		pool.DoGo(context.Background(), func(ctx context.Context, err any) {
			slog.Error("协程池请求处理异常", "error", err, "uri", c.OriginalURL(), "method", c.Method())
			ch <- tools.FiberServerError("请求处理失败")
		}, func() {
			ch <- c.Next()
		})

		select {
		// 如果请求正常完成那么直接返回
		case err := <-ch:
			return err
		// 返回请求超时错误
		case <-ctx.Done():
			return fiber.ErrRequestTimeout
		}
	}
}
