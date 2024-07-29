package category

import (
	"github.com/gofiber/fiber/v3"
	"go-fiber-ent-web-layout/internal/tools/res"
	"go-fiber-ent-web-layout/internal/usercase"
	"strconv"
)

type HttpApi struct {
	service usercase.ICategoryService
}

func NewHttpApi(service usercase.ICategoryService) *HttpApi {
	return &HttpApi{
		service: service,
	}
}

// Save 保存分类
func (h *HttpApi) Save(ctx fiber.Ctx) error {
	category := &usercase.Category{}
	if err := ctx.Bind().JSON(category); err != nil {
		return err
	}
	if err := h.service.CreateCategory(category); err != nil {
		return err
	}
	return ctx.JSON(res.OkByMessage("ok"))
}

// Update 更新分类
func (h *HttpApi) Update(ctx fiber.Ctx) error {
	category := &usercase.Category{}
	if err := ctx.Bind().JSON(category); err != nil {
		return err
	}
	if err := h.service.UpdateCategory(category); err != nil {
		return err
	}
	return ctx.JSON(res.OkByMessage("ok"))
}

func (self *HttpApi) UpdateSelective(ctx fiber.Ctx) error {
	form := &usercase.CategoryUpdateForm{}
	if err := ctx.Bind().JSON(form); err != nil {
		return err
	}
	if err := self.service.UpdateSelectiveCategory(form); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}

// Tree 博客端查询分类列表 tree
func (h *HttpApi) Tree(ctx fiber.Ctx) error {
	categorys, err := h.service.TreeCategory()
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(categorys))
}

// ManageTree 管理端查询分类列表 tree
func (h *HttpApi) ManageTree(ctx fiber.Ctx) error {
	categorys, err := h.service.ManageTreeCategory()
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(categorys))
}

// QueryInfo 查询分类详情
func (h *HttpApi) QueryInfo(ctx fiber.Ctx) error {
	id, _ := strconv.ParseInt(ctx.Params("id"), 10, 0)
	category, err := h.service.QueryCategoryInfo(int(id))
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(category))
}

// Delete 删除分类
func (h *HttpApi) Delete(ctx fiber.Ctx) error {
	id, _ := strconv.ParseInt(ctx.Params("id"), 10, 0)
	if err := h.service.Delete(int(id)); err != nil {
		return err
	}
	return ctx.JSON(res.OkByMessage("ok"))
}
