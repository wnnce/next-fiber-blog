package api

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/wire"
	"go-fiber-ent-web-layout/api/category/v1"
	"go-fiber-ent-web-layout/api/concat/v1"
	"go-fiber-ent-web-layout/api/link/v1"
	"go-fiber-ent-web-layout/api/manage/manage"
	"go-fiber-ent-web-layout/api/other/v1"
	"go-fiber-ent-web-layout/api/tag/v1"
)

var InjectSet = wire.NewSet(tag.NewHttpApi, category.NewHttpApi, concat.NewHttpApi, link.NewHttpApi, manage.NewMenuApi, manage.NewConfigApi,
	other.NewHttpApi)

// RegisterRoutes 全局路由绑定处理函数 在newApp函数中调用 不然wire无法处理依赖注入
func RegisterRoutes(app *fiber.App, tagApi *tag.HttpApi, catApi *category.HttpApi, conApi *concat.HttpApi, linkApi *link.HttpApi,
	menuApi *manage.MenuApi, cfgApi *manage.ConfigApi, oApi *other.HttpApi) {
	sysRoute := app.Group("/system")
	sysRoute.Get("/logger/sse/:interval<int;min<10>>", manage.LoggerPush)
	menuRoute := sysRoute.Group("/menu")
	menuRoute.Post("/", menuApi.Save)
	menuRoute.Put("/", menuApi.Update)
	menuRoute.Get("/tree", menuApi.Tree)
	menuRoute.Get("/manage/tree", menuApi.ManageTree)
	menuRoute.Delete("/:id<int;min<1>>", menuApi.Delete)
	cfgRoute := sysRoute.Group("/config")
	cfgRoute.Post("/", cfgApi.Save)
	cfgRoute.Put("/", cfgApi.Update)
	cfgRoute.Post("/page", cfgApi.ManagePage)
	cfgRoute.Delete("/:id<int:;min<1>>", cfgApi.Delete)

	tagRoute := app.Group("/tag")
	tagRoute.Get("/:id<int;min<1>>", tagApi.QueryInfo)
	tagRoute.Get("/list", tagApi.List)
	tagRoute.Post("/manage/list", tagApi.ManageList)
	tagRoute.Post("/", tagApi.Sava)
	tagRoute.Put("/", tagApi.Update)
	tagRoute.Delete("/:id<int;min<1>>", tagApi.Delete)

	catRoute := app.Group("/category")
	catRoute.Get("/:id<int;min=<1>>", catApi.QueryInfo)
	catRoute.Post("/", catApi.Save)
	catRoute.Put("/", catApi.Update)
	catRoute.Get("/list", catApi.List)
	catRoute.Get("/manage/list", catApi.ManageList)
	catRoute.Delete("/:id<int;min=<1>>", catApi.Delete)

	conRoute := app.Group("/concat")
	conRoute.Get("/list", conApi.List)
	conRoute.Post("/manage/list", conApi.ManageList)
	conRoute.Delete("/:id<int;min<1>>", conApi.Delete)
	conRoute.Post("/", conApi.Save)
	conRoute.Put("/", conApi.Update)

	linkRoute := app.Group("/link")
	linkRoute.Get("/list", linkApi.Page)
	linkRoute.Post("/manage/list", linkApi.ManagePage)
	linkRoute.Post("/", linkApi.Save)
	linkRoute.Put("/", linkApi.Update)
	linkRoute.Delete("/:id<int;min=<1>>", linkApi.Delete)

	othRoute := app.Group("/other")
	othRoute.Post("/upload/image", oApi.UploadImage)
	othRoute.Get("/access", oApi.AccessTrace)
}
