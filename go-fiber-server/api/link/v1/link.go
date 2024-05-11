package link

import (
	"github.com/gofiber/fiber/v3"
	"go-fiber-ent-web-layout/internal/tools/res"
	"go-fiber-ent-web-layout/internal/usercase"
	"strconv"
)

type HttpApi struct {
	service usercase.ILinkService
}

func NewHttpApi(service usercase.ILinkService) *HttpApi {
	return &HttpApi{
		service: service,
	}
}

// Save 新增友情链接
func (h *HttpApi) Save(ctx fiber.Ctx) error {
	link := &usercase.Link{}
	if err := ctx.Bind().Body(link); err != nil {
		return err
	}
	if err := h.service.CreateLike(link); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}

// Update 更新友情链接
func (h *HttpApi) Update(ctx fiber.Ctx) error {
	link := &usercase.Link{}
	if err := ctx.Bind().Body(link); err != nil {
		return err
	}
	if err := h.service.UpdateLike(link); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}

// Page 博客端查询
func (h *HttpApi) Page(ctx fiber.Ctx) error {
	query := &usercase.PageQueryForm{}
	if err := ctx.Bind().Query(query); err != nil {
		return err
	}
	page, err := h.service.PageLike(query)
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(page))
}

// ManagePage 管理端查询
func (h *HttpApi) ManagePage(ctx fiber.Ctx) error {
	query := &usercase.LinkQueryForm{}
	if err := ctx.Bind().Body(query); err != nil {
		return err
	}
	page, err := h.service.ManagePageLike(query)
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(page))
}

// Delete 删除友情链接
func (h *HttpApi) Delete(ctx fiber.Ctx) error {
	id, _ := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err := h.service.Delete(id); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}
