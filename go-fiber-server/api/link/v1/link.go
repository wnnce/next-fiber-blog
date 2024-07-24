package link

import (
	"github.com/gofiber/fiber/v3"
	"go-fiber-ent-web-layout/internal/data"
	"go-fiber-ent-web-layout/internal/tools/res"
	"go-fiber-ent-web-layout/internal/usercase"
	"strconv"
)

type HttpApi struct {
	service       usercase.ILinkService
	redisTemplate *data.RedisTemplate
}

func NewHttpApi(service usercase.ILinkService, redisTemplate *data.RedisTemplate) *HttpApi {
	return &HttpApi{
		service:       service,
		redisTemplate: redisTemplate,
	}
}

// Save 新增友情链接
func (h *HttpApi) Save(ctx fiber.Ctx) error {
	link := &usercase.Link{}
	if err := ctx.Bind().Body(link); err != nil {
		return err
	}
	if err := h.service.CreateLink(link); err != nil {
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
	if err := h.service.UpdateLink(link); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}

func (self *HttpApi) UpdateSelective(ctx fiber.Ctx) error {
	form := &usercase.LinkUpdateForm{}
	if err := ctx.Bind().JSON(form); err != nil {
		return err
	}
	if err := self.service.UpdateSelectiveLink(form); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}

// List 博客端查询友情链接列表
func (h *HttpApi) List(ctx fiber.Ctx) error {
	links, err := h.service.List()
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(links))
}

// ManagePage 管理端查询
func (h *HttpApi) ManagePage(ctx fiber.Ctx) error {
	query := &usercase.LinkQueryForm{}
	if err := ctx.Bind().Body(query); err != nil {
		return err
	}
	page, err := h.service.ManagePageLink(query)
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
