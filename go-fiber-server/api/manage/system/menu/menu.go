package menu

import (
	"github.com/gofiber/fiber/v3"
	"go-fiber-ent-web-layout/internal/tools/res"
	usercase "go-fiber-ent-web-layout/internal/usercase/system"
	"strconv"
)

type HttpApi struct {
	service usercase.IMenuService
}

func NewHttpApi(service usercase.IMenuService) *HttpApi {
	return &HttpApi{
		service: service,
	}
}

func (h *HttpApi) Save(ctx fiber.Ctx) error {
	menu := new(usercase.Menu)
	if err := ctx.Bind().JSON(menu); err != nil {
		return err
	}
	if err := h.service.CreateMenu(menu); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}

func (h *HttpApi) Update(ctx fiber.Ctx) error {
	menu := new(usercase.Menu)
	if err := ctx.Bind().JSON(menu); err != nil {
		return err
	}
	if err := h.service.UpdateMenu(menu); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}

func (h *HttpApi) Tree(ctx fiber.Ctx) error {
	menus, err := h.service.TreeMenu()
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(menus))
}

func (h *HttpApi) ManageTree(ctx fiber.Ctx) error {
	menus, err := h.service.ManageTreeMenu()
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(menus))
}

func (h *HttpApi) Delete(ctx fiber.Ctx) error {
	id, _ := strconv.ParseInt(ctx.Params("id"), 10, 0)
	if err := h.service.Delete(int(id)); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}
