// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"context"
	"github.com/gofiber/fiber/v3"
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
	"go-fiber-ent-web-layout/internal/data"
	"go-fiber-ent-web-layout/internal/service"
)

// Injectors from wire.go:

// wireApp generate inject code
func wireApp(contextContext context.Context, confData *conf.Data, jwt *conf.Jwt, server *conf.Server) (*fiber.App, func(), error) {
	dataData, cleanup, err := data.NewData(confData)
	if err != nil {
		return nil, nil, err
	}
	iTagRepo := data.NewTagRepo(dataData)
	iArticleRepo := data.NewArticleRepo(dataData)
	redisTemplate := data.NewRedisTemplate(dataData)
	iTagService := service.NewTagService(iTagRepo, iArticleRepo, redisTemplate)
	httpApi := tag.NewHttpApi(iTagService)
	iCategoryRepo := data.NewCategoryRepo(dataData)
	iCategoryService := service.NewCategoryService(iCategoryRepo, iArticleRepo, redisTemplate)
	categoryHttpApi := category.NewHttpApi(iCategoryService)
	iConcatRepo := data.NewConcatRepo(dataData)
	iConcatService := service.NewConcatService(iConcatRepo, redisTemplate)
	concatHttpApi := concat.NewHttpApi(iConcatService)
	iLinkRepo := data.NewLinkRepo(dataData)
	iLinkService := service.NewLinkService(iLinkRepo)
	linkHttpApi := link.NewHttpApi(iLinkService, redisTemplate)
	iSysMenuRepo := data.NewSysMenuRepo(dataData)
	iSysMenuService := service.NewMenuService(iSysMenuRepo)
	menuApi := manage.NewMenuApi(iSysMenuService)
	iSysConfigRepo := data.NewSysConfigRepo(dataData)
	iSysConfigService := service.NewSysConfigService(iSysConfigRepo)
	configApi := manage.NewConfigApi(iSysConfigService)
	iOtherRepo := data.NewOtherRepo(dataData)
	iOtherService := service.NewOtherService(iOtherRepo, redisTemplate)
	otherHttpApi := other.NewHttpApi(iOtherService)
	iSysRoleRepo := data.NewSysRoleRepo(dataData)
	iSysUserRepo := data.NewSysUserRepo(dataData)
	iSysRoleService := service.NewSysRoleService(iSysRoleRepo, iSysUserRepo)
	roleApi := manage.NewRoleApi(iSysRoleService)
	iSysUserService := service.NewSysUserService(iSysUserRepo, iSysRoleRepo, iOtherService)
	userApi := manage.NewUserApi(iSysUserService)
	iSysDictRepo := data.NewSysDictRepo(dataData)
	iSysDictService := service.NewSysDictService(iSysDictRepo, redisTemplate)
	dictApi := manage.NewDictApi(iSysDictService)
	iNoticeRepo := data.NewNoticeRepo(dataData)
	iNoticeService := service.NewNoticeService(iNoticeRepo, redisTemplate)
	noticeApi := manage.NewNoticeApi(iNoticeService)
	iArticleService := service.NewArticleService(iArticleRepo, redisTemplate)
	articleHttpApi := article.NewHttpApi(iArticleService)
	iTopicRepo := data.NewTopicRepo(dataData)
	iTopicService := service.NewTopicService(iTopicRepo)
	topicHttpApi := topic.NewHttpApi(iTopicService)
	iUserRepo := data.NewUserRepo(dataData)
	iUserService := service.NewUserService(iUserRepo)
	userHttpApi := user.NewHttpApi(iUserService)
	iCommentRepo := data.NewCommentRepo(dataData)
	iCommentService := service.NewCommentService(iCommentRepo, iUserService)
	commentHttpApi := comment.NewHttpApi(iCommentService)
	app := newApp(contextContext, server, httpApi, categoryHttpApi, concatHttpApi, linkHttpApi, menuApi, configApi, otherHttpApi, roleApi, userApi, dictApi, noticeApi, articleHttpApi, topicHttpApi, userHttpApi, commentHttpApi)
	return app, func() {
		cleanup()
	}, nil
}
