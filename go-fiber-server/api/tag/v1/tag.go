package tag

import (
	"github.com/gofiber/fiber/v3"
	"go-fiber-ent-web-layout/internal/tools"
	"go-fiber-ent-web-layout/internal/tools/res"
	"go-fiber-ent-web-layout/internal/usercase"
	"strconv"
)

type HttpApi struct {
	service usercase.ITagService
}

func NewHttpApi(service usercase.ITagService) *HttpApi {
	return &HttpApi{
		service: service,
	}
}

func (h *HttpApi) Sava(ctx fiber.Ctx) error {
	form := &usercase.TagForm{}
	if err := ctx.Bind().JSON(form); err != nil {
		return err
	}
	if err := h.service.CreateTag(form); err != nil {
		return err
	}
	return ctx.JSON(res.OkByMessage("ok"))
}

func (h *HttpApi) Update(ctx fiber.Ctx) error {
	form := &usercase.TagForm{}
	if err := ctx.Bind().JSON(form); err != nil {
		return err
	}
	if err := h.service.UpdateTag(form); err != nil {
		return err
	}
	return ctx.JSON(res.OkByMessage("ok"))
}

func (h *HttpApi) QueryInfo(ctx fiber.Ctx) error {
	id, _ := strconv.ParseInt(ctx.Params("id"), 10, 0)
	info, err := h.service.QueryTagInfo(int(id))
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(info))
}

func (h *HttpApi) List(ctx fiber.Ctx) error {
	tags := h.service.AllTag()
	return ctx.JSON(res.OkByData(tags))
}

func (h *HttpApi) Page(ctx fiber.Ctx) error {
	form := &usercase.TagQueryForm{}
	if err := ctx.Bind().JSON(form); err != nil {
		return tools.FiberRequestError("参数错误")
	}
	tags, err := h.service.PageTag(form)
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(tags))
}

func (h *HttpApi) Delete(ctx fiber.Ctx) error {
	id, _ := strconv.ParseInt(ctx.Params("id"), 10, 0)
	if err := h.service.Delete(int(id)); err != nil {
		return err
	}
	return ctx.JSON(res.OkByMessage("ok"))
}
