package manage

import (
	"github.com/gofiber/fiber/v3"
	"go-fiber-ent-web-layout/internal/tools/res"
	"go-fiber-ent-web-layout/internal/usercase"
	"strconv"
)

type RoleApi struct {
	service usercase.ISysRoleService
}

func NewRoleApi(service usercase.ISysRoleService) *RoleApi {
	return &RoleApi{
		service: service,
	}
}

func (self *RoleApi) Save(ctx fiber.Ctx) error {
	role := new(usercase.SysRole)
	if err := ctx.Bind().JSON(role); err != nil {
		return err
	}
	if err := self.service.SaveRole(role); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}

func (self *RoleApi) Update(ctx fiber.Ctx) error {
	role := new(usercase.SysRole)
	if err := ctx.Bind().JSON(role); err != nil {
		return err
	}
	if err := self.service.UpdateRole(role); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}

func (self *RoleApi) UpdateSelective(ctx fiber.Ctx) error {
	form := &usercase.SysRoleUpdateForm{}
	if err := ctx.Bind().JSON(form); err != nil {
		return err
	}
	if err := self.service.UpdateSelectiveRole(form); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}

func (self *RoleApi) List(ctx fiber.Ctx) error {
	list, err := self.service.List()
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(list))
}

func (self *RoleApi) Page(ctx fiber.Ctx) error {
	query := new(usercase.SysRoleQueryForm)
	if err := ctx.Bind().JSON(query); err != nil {
		return err
	}
	page, err := self.service.Page(query)
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(page))
}

func (self *RoleApi) Delete(ctx fiber.Ctx) error {
	id, _ := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err := self.service.Delete(int(id)); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}
