package user

import (
	"github.com/gofiber/fiber/v3"
	"go-fiber-ent-web-layout/internal/middleware/auth"
	"go-fiber-ent-web-layout/internal/tools"
	"go-fiber-ent-web-layout/internal/tools/res"
	"go-fiber-ent-web-layout/internal/usercase"
)

type HttpApi struct {
	service usercase.IUserService
}

func NewHttpApi(service usercase.IUserService) *HttpApi {
	return &HttpApi{
		service: service,
	}
}

func (self *HttpApi) LoginWithGithub(ctx fiber.Ctx) error {
	code := ctx.Query("code")
	if len(code) != 20 {
		return tools.FiberRequestError("code参数错误")
	}
	token, err := self.service.LoginWithGithub(code, ctx.IP())
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(token))
}

func (self *HttpApi) UserInfo(ctx fiber.Ctx) error {
	classicUser := fiber.Locals[auth.ClassicLoginUser](ctx, "classicUser")
	userinfo, err := self.service.UserInfo(classicUser.(*usercase.User))
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(userinfo))
}
