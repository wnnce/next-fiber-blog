package cors

import (
	"github.com/gofiber/fiber/v3"
	"strconv"
	"strings"
)

// NewMiddleware 创建跨域中间件
// config 跨域配置
func NewMiddleware(config CorsConfig) fiber.Handler {
	corsConfigDefault(&config)
	return func(ctx fiber.Ctx) error {
		if config.UseOrigin {
			ctx.Set(fiber.HeaderAccessControlAllowOrigin, ctx.Get(fiber.HeaderOrigin, ""))
		} else {
			ctx.Set(fiber.HeaderAccessControlAllowOrigin, config.AllowOrigin)
		}
		ctx.Set(fiber.HeaderAccessControlAllowCredentials, strconv.FormatBool(config.AllowCredentials))
		ctx.Set(fiber.HeaderAccessControlAllowMethods, strings.Join(config.AllowMethods, ","))
		ctx.Set(fiber.HeaderAccessControlAllowHeaders, strings.Join(config.AllowHeaders, ","))
		if ctx.Method() == fiber.MethodOptions {
			ctx.Set(fiber.HeaderAccessControlMaxAge, strconv.FormatInt(config.OptionMaxAge, 10))
			return ctx.SendStatus(fiber.StatusNoContent)
		}
		return ctx.Next()
	}
}
