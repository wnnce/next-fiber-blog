package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"go-fiber-ent-web-layout/api"
	"go-fiber-ent-web-layout/api/article/v1"
	"go-fiber-ent-web-layout/api/category/v1"
	"go-fiber-ent-web-layout/api/comment/v1"
	"go-fiber-ent-web-layout/api/concat/v1"
	"go-fiber-ent-web-layout/api/link/v1"
	"go-fiber-ent-web-layout/api/manage/v1"
	"go-fiber-ent-web-layout/api/other/v1"
	"go-fiber-ent-web-layout/api/tag/v1"
	"go-fiber-ent-web-layout/api/topic/v1"
	"go-fiber-ent-web-layout/api/user/v1"
	"go-fiber-ent-web-layout/internal/conf"
	"go-fiber-ent-web-layout/internal/middleware/cors"
	"go-fiber-ent-web-layout/internal/middleware/limiter"
	"go-fiber-ent-web-layout/internal/middleware/timeout"
	"go-fiber-ent-web-layout/internal/tools"
	"go-fiber-ent-web-layout/internal/tools/clog"
	"go-fiber-ent-web-layout/internal/tools/hand"
	"log/slog"
)

var confPath string

// 创建fiber app 包含注入中间件、错误处理、路由绑定等操作
func newApp(ctx context.Context, cf *conf.Server, tagApi *tag.HttpApi, catApi *category.HttpApi, conApi *concat.HttpApi,
	linkApi *link.HttpApi, menuApi *manage.MenuApi, cfgApi *manage.ConfigApi, oApi *other.HttpApi, roleApi *manage.RoleApi,
	userApi *manage.UserApi, dictApi *manage.DictApi, noticeApi *manage.NoticeApi, articleApi *article.HttpApi,
	topicApi *topic.HttpApi, classicUserApi *user.HttpApi, commentApi *comment.HttpApi) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:         cf.Name,                        // 应用名称
		ErrorHandler:    hand.CustomErrorHandler,        // 自定义错误处理器
		JSONDecoder:     sonic.Unmarshal,                // 使用sonic进行Json序列化
		JSONEncoder:     sonic.Marshal,                  // 使用sonic进行Json解析
		StructValidator: tools.DefaultStructValidator(), // 结构体参数验证
	})
	// 防止程序panic 使用自定义的处理器 记录异常
	app.Use(recover.New(recover.Config{
		EnableStackTrace:  true,
		StackTraceHandler: hand.StackTraceHandler,
	}))
	// 使用跨域中间件
	app.Use(cors.NewMiddleware(cors.DefaultCorsConfig))
	// 使用超时中间件
	app.Use(timeout.NewMiddleware(cf.Timeout))
	// 使用限流中间件
	app.Use(limiter.NewMiddleware(limiter.Config{
		KeyGenerate:     limiter.Md5KeyGenerate(),
		CallbackHandler: limiter.DefaultCallbackHandler,
		Sliding:         cf.Limiter.Sliding,
		TokenBucket:     cf.Limiter.TokenBucket,
	}, ctx))
	api.RegisterRoutes(app, tagApi, catApi, conApi, linkApi, menuApi, cfgApi, oApi, roleApi, userApi, dictApi, noticeApi,
		articleApi, topicApi, classicUserApi, commentApi)
	return app
}

func init() {
	flag.StringVar(&confPath, "conf", "/configs/config-prod.yaml", "config path, eg: -conf config-prod.yaml")
}

func main() {
	flag.Parse()
	config := conf.ReadConfig(confPath)
	// 初始化日志
	writer := &clog.CustomSlogWriter{}
	// 日志SSE端口推送
	writer.RegisterWriter(clog.GetSSEWriter())
	handler := slog.NewJSONHandler(writer, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelInfo,
	})
	slog.SetDefault(slog.New(handler).With("app-name", config.Server.Name))
	conf.IssuedConfig(config)
	ctx, cancel := context.WithCancel(context.Background())
	app, cleanup, err := wireApp(ctx, &config.Data, &config.Jwt, &config.Server)
	if err != nil {
		panic(err)
	}
	defer func() {
		cancel()
		cleanup()
	}()
	if err = app.Listen(fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)); err != nil {
		panic(err)
	}
}
