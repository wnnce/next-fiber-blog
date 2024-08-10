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

// UploadImage 图片上传
func (self *HttpApi) UploadImage(ctx fiber.Ctx) error {
	fileHeader, err := ctx.FormFile("image")
	if err != nil {
		return err
	}
	url, err := self.service.UploadImage(fileHeader)
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(url))
}

// AccessTrace 记录访问请求
func (self *HttpApi) AccessTrace(ctx fiber.Ctx) error {
	ip := ctx.IP()
	referer := ctx.Get(fiber.HeaderReferer)
	ua := ctx.Get(fiber.HeaderUserAgent)
	self.service.TraceAccess(referer, ip, ua)
	return ctx.SendStatus(fiber.StatusOK)
}

// PageLoginRecord 分页查询登录日志
func (self *HttpApi) PageLoginRecord(ctx fiber.Ctx) error {
	query := &usercase.LoginLogQueryForm{}
	if err := ctx.Bind().JSON(query); err != nil {
		return err
	}
	page, err := self.service.PageLogin(query)
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(page))
}

func (self *HttpApi) PageAccessRecord(ctx fiber.Ctx) error {
	query := &usercase.AccessLogQueryForm{}
	if err := ctx.Bind().JSON(query); err != nil {
		return err
	}
	page, err := self.service.PageAccess(query)
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(page))
}

func (self *HttpApi) QuerySiteConfiguration(ctx fiber.Ctx) error {
	return ctx.JSON(res.OkByData(self.service.SiteConfiguration()))
}

func (self *HttpApi) UpdateSiteConfiguration(ctx fiber.Ctx) error {
	config := make(map[string]usercase.SiteConfigurationItem)
	if err := ctx.Bind().JSON(&config); err != nil {
		return err
	}
	if err := self.service.UpdateSiteConfiguration(config); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}
