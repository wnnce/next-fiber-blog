package concat

import (
	"github.com/gofiber/fiber/v3"
	"go-fiber-ent-web-layout/internal/tools/res"
	"go-fiber-ent-web-layout/internal/usercase"
	"strconv"
)

type HttpApi struct {
	service usercase.IConcatService
}

func NewHttpApi(service usercase.IConcatService) *HttpApi {
	return &HttpApi{
		service: service,
	}
}

func (self *HttpApi) Save(ctx fiber.Ctx) error {
	concat := &usercase.Concat{}
	if err := ctx.Bind().Body(concat); err != nil {
		return err
	}
	if err := self.service.CreateConcat(concat); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}

func (self *HttpApi) Update(ctx fiber.Ctx) error {
	concat := &usercase.Concat{}
	if err := ctx.Bind().JSON(concat); err != nil {
		return err
	}
	if err := self.service.UpdateConcat(concat); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}

func (self *HttpApi) UpdateSelective(ctx fiber.Ctx) error {
	form := &usercase.ConcatUpdateForm{}
	if err := ctx.Bind().JSON(form); err != nil {
		return err
	}
	if err := self.service.UpdateSelectiveConcat(form); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}

func (self *HttpApi) List(ctx fiber.Ctx) error {
	concats, err := self.service.ListConcat()
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(concats))
}

func (self *HttpApi) ManageList(ctx fiber.Ctx) error {
	query := &usercase.ConcatQueryForm{}
	_ = ctx.Bind().JSON(query)
	concats, err := self.service.ManageListConcat(query)
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(concats))
}

func (self *HttpApi) Delete(ctx fiber.Ctx) error {
	id, _ := strconv.ParseInt(ctx.Params("id"), 10, 0)
	if err := self.service.Delete(int(id)); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}
