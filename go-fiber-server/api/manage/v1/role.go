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

func (r *RoleApi) Save(ctx fiber.Ctx) error {
	role := new(usercase.SysRole)
	if err := ctx.Bind().JSON(role); err != nil {
		return err
	}
	if err := r.service.SaveRole(role); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}

func (r *RoleApi) Update(ctx fiber.Ctx) error {
	role := new(usercase.SysRole)
	if err := ctx.Bind().JSON(role); err != nil {
		return err
	}
	if err := r.service.UpdateRole(role); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}

func (r *RoleApi) List(ctx fiber.Ctx) error {
	list, err := r.service.List()
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(list))
}

func (r *RoleApi) Page(ctx fiber.Ctx) error {
	query := new(usercase.SysRoleQueryForm)
	if err := ctx.Bind().JSON(query); err != nil {
		return err
	}
	page, err := r.service.Page(query)
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(page))
}

func (r *RoleApi) Delete(ctx fiber.Ctx) error {
	id, _ := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err := r.service.Delete(int(id)); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}
