package api

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/wire"
	"go-fiber-ent-web-layout/api/category/v1"
	"go-fiber-ent-web-layout/api/concat/v1"
	"go-fiber-ent-web-layout/api/link/v1"
	"go-fiber-ent-web-layout/api/manage/v1"
	"go-fiber-ent-web-layout/api/other/v1"
	"go-fiber-ent-web-layout/api/tag/v1"
	"go-fiber-ent-web-layout/internal/middleware/auth"
)

var InjectSet = wire.NewSet(tag.NewHttpApi, category.NewHttpApi, concat.NewHttpApi, link.NewHttpApi, manage.NewMenuApi,
	manage.NewConfigApi, other.NewHttpApi, manage.NewRoleApi, manage.NewUserApi, manage.NewDictApi)

// RegisterRoutes 全局路由绑定处理函数 在newApp函数中调用 不然wire无法处理依赖注入
func RegisterRoutes(app *fiber.App, tagApi *tag.HttpApi, catApi *category.HttpApi, conApi *concat.HttpApi, linkApi *link.HttpApi,
	menuApi *manage.MenuApi, cfgApi *manage.ConfigApi, oApi *other.HttpApi, roleApi *manage.RoleApi, userApi *manage.UserApi,
	dictApi *manage.DictApi) {
	// 系统接口路由
	sysRoute := app.Group("/system", auth.ManageAuth, auth.VerifyRoles("admin"))
	sysRoute.Post("/record/login", oApi.PageLoginRecord)
	sysRoute.Post("/record/access", oApi.PageAccessRecord)
	// 菜单管理接口
	menuRoute := sysRoute.Group("/menu")
	menuRoute.Post("/", menuApi.Save)
	menuRoute.Put("/", menuApi.Update)
	menuRoute.Get("/manage/tree", menuApi.ManageTree)
	menuRoute.Delete("/:id<int;min<1>>", menuApi.Delete)
	// 参数管理接口
	cfgRoute := sysRoute.Group("/config")
	cfgRoute.Post("/", cfgApi.Save)
	cfgRoute.Put("/", cfgApi.Update)
	cfgRoute.Post("/page", cfgApi.ManagePage)
	cfgRoute.Delete("/:id<int:;min<1>>", cfgApi.Delete)
	// 角色管理接口
	roleRoute := sysRoute.Group("/role")
	roleRoute.Post("/", roleApi.Save)
	roleRoute.Put("/", roleApi.Update)
	roleRoute.Post("/page", roleApi.Page)
	roleRoute.Get("/list", roleApi.List)
	roleRoute.Delete("/:id<int;min=<1>>", roleApi.Delete)
	// 用户管理接口
	userRoute := sysRoute.Group("/user")
	userRoute.Post("/", userApi.Save)
	userRoute.Put("/", userApi.Update)
	userRoute.Post("/page", userApi.Page)
	userRoute.Delete("/:id<int;min<1>>", userApi.Delete)
	// 字典管理接口
	dictRoute := sysRoute.Group("/dict")
	dictRoute.Post("/", dictApi.SaveDict)
	dictRoute.Put("/", dictApi.UpdateDict)
	dictRoute.Put("/status", dictApi.UpdateDictStatus)
	dictRoute.Post("/page", dictApi.PageDict)
	dictRoute.Delete("/:id<int:min<1>>", dictApi.DeleteDict)
	dictRoute.Post("/value", dictApi.SaveDictValue)
	dictRoute.Put("/value", dictApi.UpdateDictValue)
	dictRoute.Put("/value/status", dictApi.UpdateDictValueStatus)
	dictRoute.Post("/value/page", dictApi.PageDictValue)
	dictRoute.Delete("/value/:id<int:min<1>>", dictApi.DeleteDictValue)

	// 系统公共接口路由 只需登录即可访问
	baseRoute := app.Group("/base", auth.ManageAuth)
	// 查询用户菜单
	baseRoute.Get("/menu-tree", menuApi.Tree)
	// 获取用户详情
	baseRoute.Get("/user-info", userApi.UserInfo)
	// 修改密码
	baseRoute.Put("/re-password", userApi.UpdatePassword)
	// 退出登录
	baseRoute.Post("/logout", userApi.Logout)
	// 图片上传
	baseRoute.Post("/upload/image", oApi.UploadImage)

	// 标签管理接口
	tagRoute := app.Group("/tag", auth.ManageAuth)
	tagRoute.Get("/:id<int;min<1>>", tagApi.QueryInfo)
	tagRoute.Get("/list", tagApi.List)
	tagRoute.Post("/page", tagApi.Page)
	tagRoute.Post("/", tagApi.Sava)
	tagRoute.Put("/", tagApi.Update)
	tagRoute.Delete("/:id<int;min<1>>", tagApi.Delete)

	// 分类管理接口
	catRoute := app.Group("/category", auth.ManageAuth)
	catRoute.Get("/:id<int;min=<1>>", catApi.QueryInfo)
	catRoute.Post("/", catApi.Save)
	catRoute.Put("/", catApi.Update)
	catRoute.Get("/list", catApi.List)
	catRoute.Get("/manage/list", catApi.ManageList)
	catRoute.Delete("/:id<int;min=<1>>", catApi.Delete)

	// 联系方式接口
	conRoute := app.Group("/concat", auth.ManageAuth)
	conRoute.Get("/list", conApi.List)
	conRoute.Post("/manage/list", conApi.ManageList)
	conRoute.Delete("/:id<int;min<1>>", conApi.Delete)
	conRoute.Post("/", conApi.Save)
	conRoute.Put("/", conApi.Update)

	// 友情链接接口
	linkRoute := app.Group("/link", auth.ManageAuth)
	linkRoute.Get("/list", linkApi.Page)
	linkRoute.Post("/page", linkApi.ManagePage)
	linkRoute.Post("/", linkApi.Save)
	linkRoute.Put("/", linkApi.Update)
	linkRoute.Delete("/:id<int;min=<1>>", linkApi.Delete)

	// 开放接口
	openRoute := app.Group("/open")
	openRoute.Get("/logger/stream/:interval<int;min<10>>", manage.LoggerPush)
	openRoute.Post("/login", userApi.Login)
	openRoute.Get("/trace/access", oApi.AccessTrace)
	openRoute.Get("/dict/:dictKey", dictApi.ListDictValue)
}
