// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"context"
	"github.com/gofiber/fiber/v3"
	"go-fiber-ent-web-layout/api/category/v1"
	"go-fiber-ent-web-layout/api/concat/v1"
	"go-fiber-ent-web-layout/api/tag/v1"
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
	iTagService := service.NewTagService(iTagRepo)
	httpApi := tag.NewHttpApi(iTagService)
	iCategoryRepo := data.NewCategoryRepo(dataData)
	iCategoryService := service.NewCategoryService(iCategoryRepo)
	categoryHttpApi := category.NewHttpApi(iCategoryService)
	iConcatRepo := data.NewConcatRepo(dataData)
	iConcatService := service.NewConcatService(iConcatRepo)
	concatHttpApi := concat.NewHttpApi(iConcatService)
	app := newApp(contextContext, server, httpApi, categoryHttpApi, concatHttpApi)
	return app, func() {
		cleanup()
	}, nil
}
