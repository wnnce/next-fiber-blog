package manage

import (
	"github.com/gofiber/fiber/v3"
	"go-fiber-ent-web-layout/internal/middleware/auth"
	"go-fiber-ent-web-layout/internal/tools/res"
	"go-fiber-ent-web-layout/internal/usercase"
	"strconv"
)

type UserApi struct {
	service usercase.ISysUserService
}

func NewUserApi(service usercase.ISysUserService) *UserApi {
	return &UserApi{
		service: service,
	}
}

func (u *UserApi) Save(ctx fiber.Ctx) error {
	user := new(usercase.SysUser)
	if err := ctx.Bind().JSON(user); err != nil {
		return err
	}
	if err := u.service.SaveUser(user); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}

func (u *UserApi) Update(ctx fiber.Ctx) error {
	user := new(usercase.SysUser)
	if err := ctx.Bind().JSON(user); err != nil {
		return err
	}
	if err := u.service.UpdateUser(user); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}

func (self *UserApi) UpdateSelective(ctx fiber.Ctx) error {
	form := &usercase.SysUserUpdateForm{}
	if err := ctx.Bind().JSON(form); err != nil {
		return nil
	}
	if err := self.service.UpdateSelectiveUser(form); err != nil {
		return nil
	}
	return ctx.JSON(res.SimpleOK())
}

func (u *UserApi) UserInfo(ctx fiber.Ctx) error {
	loginUser := fiber.Locals[auth.LoginUser](ctx, "loginUser")
	userInfo, err := u.service.QueryUserInfo(loginUser.GetUserId())
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(userInfo))
}

func (u *UserApi) Page(ctx fiber.Ctx) error {
	query := new(usercase.SysUserQueryForm)
	if err := ctx.Bind().JSON(query); err != nil {
		return err
	}
	page, err := u.service.Page(query)
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(page))
}

func (u *UserApi) UpdatePassword(ctx fiber.Ctx) error {
	form := new(usercase.UpdatePasswordForm)
	if err := ctx.Bind().JSON(form); err != nil {
		return err
	}
	loginUser := fiber.Locals[auth.LoginUser](ctx, "loginUser")
	form.UserId = loginUser.GetUserId()
	if err := u.service.UpdatePassword(form); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}

func (u *UserApi) Delete(ctx fiber.Ctx) error {
	id, _ := strconv.ParseInt(ctx.Params("id"), 10, 0)
	if err := u.service.Delete(id); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}

func (u *UserApi) Login(ctx fiber.Ctx) error {
	form := new(usercase.LoginForm)
	if err := ctx.Bind().JSON(form); err != nil {
		return err
	}
	ip := ctx.IP()
	ua := ctx.Get(fiber.HeaderUserAgent)
	token, err := u.service.Login(form, ip, ua)
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(token))
}

func (self *UserApi) Logout(ctx fiber.Ctx) error {
	loginUser := fiber.Locals[auth.LoginUser](ctx, "loginUser")
	self.service.Logout(loginUser.GetUserId())
	return ctx.JSON(res.SimpleOK())
}
