package category

import (
	"github.com/gofiber/fiber/v3"
	"go-fiber-ent-web-layout/internal/tools"
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

func (h *HttpApi) Save(ctx fiber.Ctx) error {
	category := &usercase.Category{}
	if err := ctx.Bind().JSON(&category); err != nil {
		return tools.FiberRequestError("参数错误")
	}
	if validation := tools.StructFieldValidation(category); validation != "" {
		return tools.FiberRequestError(validation)
	}
	if err := h.service.CreateCategory(category); err != nil {
		return err
	}
	return ctx.JSON(res.OkByMessage("ok"))
}

func (h *HttpApi) Update(ctx fiber.Ctx) error {
	category := &usercase.Category{}
	if err := ctx.Bind().JSON(&category); err != nil {
		return tools.FiberRequestError("参数错误")
	}
	if validation := tools.StructFieldValidation(category); validation != "" {
		return tools.FiberRequestError(validation)
	}
	if err := h.service.UpdateCategory(category); err != nil {
		return err
	}
	return ctx.JSON(res.OkByMessage("ok"))
}

func (h *HttpApi) List(ctx fiber.Ctx) error {
	categorys, err := h.service.ListCategory()
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(categorys))
}

func (h *HttpApi) ManageList(ctx fiber.Ctx) error {
	categorys, err := h.service.ManageListCategory()
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(categorys))
}

func (h *HttpApi) QueryInfo(ctx fiber.Ctx) error {
	id, _ := strconv.ParseInt(ctx.Params("id"), 10, 0)
	category, err := h.service.QueryCategoryInfo(int(id))
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(category))
}
