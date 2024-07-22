package manage

import (
	"github.com/gofiber/fiber/v3"
	"go-fiber-ent-web-layout/internal/tools/res"
	"go-fiber-ent-web-layout/internal/usercase"
	"strconv"
)

type NoticeApi struct {
	service usercase.INoticeService
}

func NewNoticeApi(service usercase.INoticeService) *NoticeApi {
	return &NoticeApi{
		service: service,
	}
}

func (self *NoticeApi) Save(ctx fiber.Ctx) error {
	notice := &usercase.Notice{}
	if err := ctx.Bind().JSON(notice); err != nil {
		return err
	}
	if err := self.service.SaveNotice(notice); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}

func (self *NoticeApi) Update(ctx fiber.Ctx) error {
	notice := &usercase.Notice{}
	if err := ctx.Bind().JSON(notice); err != nil {
		return err
	}
	if err := self.service.UpdateNotice(notice); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}

func (self *NoticeApi) Page(ctx fiber.Ctx) error {
	query := &usercase.NoticeQueryForm{}
	if err := ctx.Bind().JSON(query); err != nil {
		return err
	}
	page, err := self.service.Page(query)
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(page))
}

func (self *NoticeApi) ListIndexNotice(ctx fiber.Ctx) error {
	notices, err := self.service.ListNoticeByType(1)
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(notices))
}

func (self *NoticeApi) ListPublicNotice(ctx fiber.Ctx) error {
	notices, err := self.service.ListNoticeByType(2)
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(notices))
}

func (self *NoticeApi) ListAdminNotice(ctx fiber.Ctx) error {
	notices, err := self.service.ListNoticeByType(3)
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(notices))
}

func (self *NoticeApi) Delete(ctx fiber.Ctx) error {
	noticeId, _ := strconv.ParseInt(ctx.Params("id"), 10, 0)
	if err := self.service.Delete(noticeId); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}
