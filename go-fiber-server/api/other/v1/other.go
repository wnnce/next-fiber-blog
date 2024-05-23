package other

import (
	"github.com/gofiber/fiber/v3"
	"go-fiber-ent-web-layout/internal/tools/res"
	"go-fiber-ent-web-layout/internal/usercase"
)

type HttpApi struct {
	service usercase.IOtherService
}

func NewHttpApi(service usercase.IOtherService) *HttpApi {
	return &HttpApi{
		service: service,
	}
}

func (h *HttpApi) UploadImage(ctx fiber.Ctx) error {
	fileHeader, err := ctx.FormFile("image")
	if err != nil {
		return err
	}
	url, err := h.service.UploadImage(fileHeader)
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(url))
}

func (h *HttpApi) AccessTrace(ctx fiber.Ctx) error {
	ip := ctx.IP()
	referer := ctx.Get(fiber.HeaderReferer)
	ua := ctx.Get(fiber.HeaderUserAgent)
	h.service.TraceAccess(referer, ip, ua)
	return ctx.SendStatus(fiber.StatusOK)
}
