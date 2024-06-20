package cors

import "github.com/gofiber/fiber/v3"

// CorsConfig 跨域中间件配置
type CorsConfig struct {
	UseOrigin           bool     // 是否直接使用origin做为跨域地址
	AllowOrigin         string   // 放行的域名 如果useOrigin为true 可以忽略
	AllowCredentials    bool     // 是否允许携带凭证
	AllowMethods        []string // 允许跨域的请求方法
	AllowHeaders        []string // 允许携带的请求头
	OptionMaxAge        int64    // 预检请求的缓存时间
	ReleaseOptionMethod bool     // 是否直接处理Option预检请求
}

// DefaultCorsConfig 默认的跨域中间件配置
var DefaultCorsConfig = CorsConfig{
	UseOrigin:           true,
	AllowOrigin:         "*",
	AllowCredentials:    true,
	AllowMethods:        []string{fiber.MethodGet, fiber.MethodPost, fiber.MethodPut, fiber.MethodDelete, fiber.MethodOptions},
	AllowHeaders:        []string{"*"},
	OptionMaxAge:        3600,
	ReleaseOptionMethod: true,
}

// 检查跨域中间件的配置
func corsConfigDefault(cfg *CorsConfig) {
	if cfg.AllowOrigin == "" {
		cfg.AllowOrigin = DefaultCorsConfig.AllowOrigin
	}
	if cfg.AllowMethods == nil || len(cfg.AllowMethods) == 0 {
		cfg.AllowMethods = DefaultCorsConfig.AllowMethods
	}
	if cfg.AllowHeaders == nil || len(cfg.AllowHeaders) == 0 {
		cfg.AllowHeaders = DefaultCorsConfig.AllowHeaders
	}
	if cfg.OptionMaxAge < 0 {
		cfg.OptionMaxAge = DefaultCorsConfig.OptionMaxAge
	}
}
