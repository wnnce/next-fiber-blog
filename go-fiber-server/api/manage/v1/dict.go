package manage

import (
	"github.com/gofiber/fiber/v3"
	"go-fiber-ent-web-layout/internal/tools/res"
	"go-fiber-ent-web-layout/internal/usercase"
	"strconv"
)

type DictApi struct {
	service usercase.ISysDictService
}

func NewDictApi(service usercase.ISysDictService) *DictApi {
	return &DictApi{
		service: service,
	}
}

func (self *DictApi) PageDict(ctx fiber.Ctx) error {
	query := &usercase.SysDictQueryForm{}
	if err := ctx.Bind().JSON(query); err != nil {
		return err
	}
	page, err := self.service.PageDict(query)
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(page))
}

func (self *DictApi) SaveDict(ctx fiber.Ctx) error {
	dict := &usercase.SysDict{}
	if err := ctx.Bind().JSON(dict); err != nil {
		return err
	}
	if err := self.service.SaveDict(dict); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}

func (self *DictApi) UpdateDict(ctx fiber.Ctx) error {
	dict := &usercase.SysDict{}
	if err := ctx.Bind().JSON(dict); err != nil {
		return err
	}
	if err := self.service.UpdateDict(dict); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}

func (self *DictApi) UpdateDictStatus(ctx fiber.Ctx) error {
	form := &usercase.SysDictSelectiveUpdateForm{}
	if err := ctx.Bind().JSON(form); err != nil {
		return err
	}
	if err := self.service.UpdateSelectiveDict(form); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}

func (self *DictApi) DeleteDict(ctx fiber.Ctx) error {
	dictId, _ := strconv.ParseInt(ctx.Params("id"), 10, 0)
	if err := self.service.DeleteDict(dictId); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}

func (self *DictApi) PageDictValue(ctx fiber.Ctx) error {
	query := &usercase.SysDictValueQueryForm{}
	if err := ctx.Bind().JSON(query); err != nil {
		return err
	}
	page, err := self.service.PageDictValue(query)
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(page))
}

func (self *DictApi) SaveDictValue(ctx fiber.Ctx) error {
	dict := &usercase.SysDictValue{}
	if err := ctx.Bind().JSON(dict); err != nil {
		return err
	}
	if err := self.service.SaveDictValue(dict); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}

func (self *DictApi) UpdateDictValue(ctx fiber.Ctx) error {
	dict := &usercase.SysDictValue{}
	if err := ctx.Bind().JSON(dict); err != nil {
		return err
	}
	if err := self.service.UpdateDictValue(dict); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}

func (self *DictApi) UpdateDictValueStatus(ctx fiber.Ctx) error {
	form := &usercase.SysDictValueSelectiveUpdateForm{}
	if err := ctx.Bind().JSON(form); err != nil {
		return err
	}
	if err := self.service.UpdateSelectiveValue(form); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}

func (self *DictApi) DeleteDictValue(ctx fiber.Ctx) error {
	valueId, _ := strconv.ParseInt(ctx.Params("id"), 10, 0)
	if err := self.service.DeleteDictValue(valueId); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}
