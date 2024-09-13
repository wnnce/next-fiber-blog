package api

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/wire"
	"go-fiber-ent-web-layout/api/article/v1"
	"go-fiber-ent-web-layout/api/category/v1"
	"go-fiber-ent-web-layout/api/concat/v1"
	"go-fiber-ent-web-layout/api/link/v1"
	"go-fiber-ent-web-layout/api/manage/v1"
	"go-fiber-ent-web-layout/api/other/v1"
	"go-fiber-ent-web-layout/api/tag/v1"
	"go-fiber-ent-web-layout/api/topic/v1"
	"go-fiber-ent-web-layout/api/user/v1"
	"go-fiber-ent-web-layout/internal/middleware/auth"
)

var InjectSet = wire.NewSet(tag.NewHttpApi, category.NewHttpApi, concat.NewHttpApi, link.NewHttpApi, article.NewHttpApi,
	manage.NewMenuApi, manage.NewConfigApi, other.NewHttpApi, manage.NewRoleApi, manage.NewUserApi, manage.NewDictApi,
	manage.NewNoticeApi, topic.NewHttpApi, user.NewHttpApi)

// RegisterRoutes 全局路由绑定处理函数 在newApp函数中调用 不然wire无法处理依赖注入
func RegisterRoutes(app *fiber.App, tagApi *tag.HttpApi, catApi *category.HttpApi, conApi *concat.HttpApi, linkApi *link.HttpApi,
	menuApi *manage.MenuApi, cfgApi *manage.ConfigApi, oApi *other.HttpApi, roleApi *manage.RoleApi, userApi *manage.UserApi,
	dictApi *manage.DictApi, noticeApi *manage.NoticeApi, articleApi *article.HttpApi, topicApi *topic.HttpApi,
	classicUserApi *user.HttpApi) {
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
	cfgRoute.Delete("/:id<int;min<1>>", cfgApi.Delete)
	// 角色管理接口
	roleRoute := sysRoute.Group("/role")
	roleRoute.Post("/", roleApi.Save)
	roleRoute.Put("/", roleApi.Update)
	roleRoute.Put("/status", roleApi.UpdateSelective)
	roleRoute.Post("/page", roleApi.Page)
	roleRoute.Get("/list", roleApi.List)
	roleRoute.Delete("/:id<int;min=<1>>", roleApi.Delete)
	// 用户管理接口
	userRoute := sysRoute.Group("/user")
	userRoute.Post("/", userApi.Save)
	userRoute.Put("/", userApi.Update)
	userRoute.Put("/status", userApi.UpdateSelective)
	userRoute.Post("/page", userApi.Page)
	userRoute.Delete("/:id<int;min<1>>", userApi.Delete)
	// 字典管理接口
	dictRoute := sysRoute.Group("/dict")
	dictRoute.Post("/", dictApi.SaveDict)
	dictRoute.Put("/", dictApi.UpdateDict)
	dictRoute.Put("/status", dictApi.UpdateDictStatus)
	dictRoute.Post("/page", dictApi.PageDict)
	dictRoute.Delete("/:id<int;min<1>>", dictApi.DeleteDict)
	dictRoute.Post("/value", dictApi.SaveDictValue)
	dictRoute.Put("/value", dictApi.UpdateDictValue)
	dictRoute.Put("/value/status", dictApi.UpdateDictValueStatus)
	dictRoute.Post("/value/page", dictApi.PageDictValue)
	dictRoute.Delete("/value/:id<int;min<1>>", dictApi.DeleteDictValue)
	// 通知公告管理接口
	noticeRoute := sysRoute.Group("/notice")
	noticeRoute.Post("/", noticeApi.Save)
	noticeRoute.Put("/", noticeApi.Update)
	noticeRoute.Post("/page", noticeApi.Page)
	noticeRoute.Delete("/:id<int<min:1>>", noticeApi.Delete)

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
	// 获取管理端通知公告
	baseRoute.Post("/notice/admin", noticeApi.ListAdminNotice)
	// 更新站点配置
	baseRoute.Put("/site/configuration", oApi.UpdateSiteConfiguration)

	// 标签管理接口
	tagRoute := app.Group("/tag", auth.ManageAuth)
	tagRoute.Get("/:id", tagApi.QueryInfo)
	tagRoute.Get("/list", tagApi.List)
	tagRoute.Post("/page", tagApi.Page)
	tagRoute.Post("/", tagApi.Sava)
	tagRoute.Put("/", tagApi.Update)
	tagRoute.Put("/status", tagApi.UpdateStatus)
	tagRoute.Delete("/:id<int;min<1>>", tagApi.Delete)

	// 分类管理接口
	catRoute := app.Group("/category", auth.ManageAuth)
	catRoute.Get("/:id", catApi.QueryInfo)
	catRoute.Post("/", catApi.Save)
	catRoute.Put("/", catApi.Update)
	catRoute.Put("/status", catApi.UpdateSelective)
	catRoute.Get("/manage/tree", catApi.ManageTree)
	catRoute.Delete("/:id<int;min=<1>>", catApi.Delete)

	// 联系方式接口
	conRoute := app.Group("/concat", auth.ManageAuth)
	conRoute.Get("/list", conApi.List)
	conRoute.Post("/manage/list", conApi.ManageList)
	conRoute.Delete("/:id<int;min<1>>", conApi.Delete)
	conRoute.Post("/", conApi.Save)
	conRoute.Put("/", conApi.Update)
	conRoute.Put("/status", conApi.UpdateSelective)

	// 友情链接接口
	linkRoute := app.Group("/link", auth.ManageAuth)
	linkRoute.Post("/page", linkApi.ManagePage)
	linkRoute.Post("/", linkApi.Save)
	linkRoute.Put("/", linkApi.Update)
	linkRoute.Put("/status", linkApi.UpdateSelective)
	linkRoute.Delete("/:id<int;min<1>>", linkApi.Delete)

	// 文章管理接口
	articleRoute := app.Group("/article", auth.ManageAuth)
	articleRoute.Post("/", articleApi.Save)
	articleRoute.Put("/", articleApi.Update)
	articleRoute.Put("/status", articleApi.UpdateSelective)
	articleRoute.Post("/manage/page", articleApi.ManagePage)
	articleRoute.Delete("/:id", articleApi.Delete)
	articleRoute.Get("/info/:id", articleApi.ManageQueryInfo)

	// 动态管理接口
	topicRoute := app.Group("/topic", auth.ManageAuth)
	topicRoute.Post("/", topicApi.Save)
	topicRoute.Put("/", topicApi.Update)
	topicRoute.Put("/status", topicApi.UpdateSelective)
	topicRoute.Post("/page", topicApi.ManagePage)
	topicRoute.Delete("/:id", topicApi.Delete)

	classicUserRoute := app.Group("/user")
	classicUserRoute.Get("/info", classicUserApi.UserInfo, auth.ClassicAuth)

	// 开放接口
	openRoute := app.Group("/open")
	// 日志推送接口
	openRoute.Get("/logger/stream/:interval<int;min<10>>", manage.LoggerPush)
	// 管理端登录接口
	openRoute.Post("/login", userApi.Login)
	// 博客端Github登录接口
	openRoute.Get("/classic/login/github", classicUserApi.LoginWithGithub)
	// 博客请求记录接口
	openRoute.Get("/trace/access", oApi.AccessTrace)
	// 获取字典数据接口
	openRoute.Get("/dict/:dictKey", dictApi.ListDictValue)
	// 获取站点配置接口
	openRoute.Get("/site/configuration", oApi.QuerySiteConfiguration)
	// 获取站点统计数据
	openRoute.Get("/site/stats", oApi.SiteStats)
	// 获取弹窗通知
	openRoute.Get("/notice/index", noticeApi.ListIndexNotice)
	// 获取公告通知
	openRoute.Get("/notice/public", noticeApi.ListPublicNotice)
	// 获取分类树形列表
	openRoute.Get("/category/list", catApi.Tree)
	// 获取分类详情
	openRoute.Get("/category/:id", catApi.QueryInfo)
	// 获取标签列表
	openRoute.Get("/tag/list", tagApi.List)
	// 获取标签详情
	openRoute.Get("/tag/:id", tagApi.QueryInfo)
	// 获取联系方式列表
	openRoute.Get("/concat/list", conApi.List)
	// 获取友情链接列表
	openRoute.Get("/link/list", linkApi.List)
	// 分页查询博客文章列表
	openRoute.Post("/article/page", articleApi.Page)
	// 分页查询分类和标签的文章列表
	openRoute.Post("/article/label/page", articleApi.PageByLabel)
	// 分页查询动态列表
	openRoute.Post("/topic/page", topicApi.Page)
}
