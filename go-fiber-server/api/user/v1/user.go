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

func (self *HttpApi) Logout(ctx fiber.Ctx) error {
	classicUser := fiber.Locals[auth.ClassicLoginUser](ctx, "classicUser")
	if err := self.service.Logout(classicUser.GetUserId()); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}

func (self *HttpApi) Page(ctx fiber.Ctx) error {
	query := &usercase.UserQueryForm{}
	if err := ctx.Bind().JSON(query); err != nil {
		return err
	}
	page, err := self.service.PageUser(query)
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(page))
}

func (self *HttpApi) Update(ctx fiber.Ctx) error {
	user := &usercase.User{}
	if err := ctx.Bind().JSON(user); err != nil {
		return err
	}
	if err := self.service.UpdateUser(user); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}

func (self *HttpApi) PageExpertise(ctx fiber.Ctx) error {
	query := &usercase.ExpertiseQueryForm{}
	if err := ctx.Bind().JSON(query); err != nil {
		return err
	}
	page, err := self.service.PageExpertise(query)
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(page))
}
