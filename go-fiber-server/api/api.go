package api

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/wire"
	"go-fiber-ent-web-layout/api/manage"
)

var InjectSet = wire.NewSet()

// RegisterRoutes 全局路由绑定处理函数 在newApp函数中调用 不然wire无法处理依赖注入
func RegisterRoutes(app *fiber.App) {
	manageRoute := app.Group("/manage")
	manageRoute.Get("/logger/sse/:interval<int;min<100>>", manage.LoggerPush)
}
