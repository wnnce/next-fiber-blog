//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"github.com/gofiber/fiber/v3"
	"github.com/google/wire"
	"go-fiber-ent-web-layout/internal/conf"
)

// wireApp generate inject code
func wireApp(context.Context, *conf.Data, *conf.Jwt, *conf.Server) (*fiber.App, func(), error) {
	panic(wire.Build(newApp))
}
