package manage

import (
	"github.com/gofiber/fiber/v3"
	"go-fiber-ent-web-layout/internal/middleware/auth"
	"go-fiber-ent-web-layout/internal/tools/res"
	"go-fiber-ent-web-layout/internal/usercase"
	"strconv"
)

type MenuApi struct {
	service usercase.ISysMenuService
}

func NewMenuApi(service usercase.ISysMenuService) *MenuApi {
	return &MenuApi{
		service: service,
	}
}

func (h *MenuApi) Save(ctx fiber.Ctx) error {
	menu := new(usercase.SysMenu)
	if err := ctx.Bind().JSON(menu); err != nil {
		return err
	}
	if err := h.service.CreateMenu(menu); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}

func (h *MenuApi) Update(ctx fiber.Ctx) error {
	menu := new(usercase.SysMenu)
	if err := ctx.Bind().JSON(menu); err != nil {
		return err
	}
	if err := h.service.UpdateMenu(menu); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}

func (h *MenuApi) Tree(ctx fiber.Ctx) error {
	loginUser := fiber.Locals[auth.LoginUser](ctx, "loginUser")
	menus, err := h.service.TreeMenu(loginUser.GetRoles())
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(menus))
}

func (h *MenuApi) ManageTree(ctx fiber.Ctx) error {
	menus, err := h.service.ManageTreeMenu()
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(menus))
}

func (h *MenuApi) Delete(ctx fiber.Ctx) error {
	id, _ := strconv.ParseInt(ctx.Params("id"), 10, 0)
	if err := h.service.Delete(int(id)); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}
