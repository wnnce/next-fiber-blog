package manage

import (
	"github.com/gofiber/fiber/v3"
	"go-fiber-ent-web-layout/internal/tools/res"
	"go-fiber-ent-web-layout/internal/usercase"
	"strconv"
)

type ConfigApi struct {
	service usercase.ISysConfigService
}

func NewConfigApi(service usercase.ISysConfigService) *ConfigApi {
	return &ConfigApi{
		service: service,
	}
}

func (ca *ConfigApi) Save(ctx fiber.Ctx) error {
	cfg := new(usercase.SysConfig)
	if err := ctx.Bind().JSON(cfg); err != nil {
		return err
	}
	if err := ca.service.CreateConfig(cfg); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}

func (ca *ConfigApi) Update(ctx fiber.Ctx) error {
	cfg := new(usercase.SysConfig)
	if err := ctx.Bind().JSON(cfg); err != nil {
		return err
	}
	if err := ca.service.UpdateConfig(cfg); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}

func (ca *ConfigApi) ManagePage(ctx fiber.Ctx) error {
	query := new(usercase.SysConfigQueryForm)
	if err := ctx.Bind().JSON(query); err != nil {
		return err
	}
	page, err := ca.service.ManageList(query)
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(page))
}

func (ca *ConfigApi) Delete(ctx fiber.Ctx) error {
	id, _ := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err := ca.service.Delete(id); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}
