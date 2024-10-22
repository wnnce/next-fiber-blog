package comment

import (
	"github.com/gofiber/fiber/v3"
	"go-fiber-ent-web-layout/internal/middleware/auth"
	"go-fiber-ent-web-layout/internal/tools/res"
	"go-fiber-ent-web-layout/internal/usercase"
)

type HttpApi struct {
	service usercase.ICommentService
}

func NewHttpApi(service usercase.ICommentService) *HttpApi {
	return &HttpApi{
		service: service,
	}
}

func (self *HttpApi) Total(ctx fiber.Ctx) error {
	query := &usercase.CommentQueryForm{}
	if err := ctx.Bind().JSON(query); err != nil {
		return err
	}
	total, err := self.service.TotalComment(query)
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(total))
}

func (self *HttpApi) Save(ctx fiber.Ctx) error {
	comment := &usercase.Comment{}
	if err := ctx.Bind().JSON(comment); err != nil {
		return err
	}
	loginUser := fiber.Locals[auth.ClassicLoginUser](ctx, "classicUser")
	comment.UserId = loginUser.GetUserId()
	comment.CommentIp = ctx.IP()
	comment.CommentUa = ctx.Get(fiber.HeaderUserAgent)
	if err := self.service.SaveComment(comment); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}

func (self *HttpApi) Page(ctx fiber.Ctx) error {
	query := &usercase.CommentQueryForm{}
	if err := ctx.Bind().JSON(query); err != nil {
		return err
	}
	page, err := self.service.Page(query)
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(page))
}

func (self *HttpApi) UpdateSelective(ctx fiber.Ctx) error {
	form := &usercase.CommentUpdateForm{}
	if err := ctx.Bind().JSON(form); err != nil {
		return err
	}
	if err := self.service.UpdateSelectiveComment(form); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}

func (self *HttpApi) Delete(ctx fiber.Ctx) error {
	commentId := fiber.Params[int64](ctx, "id")
	if err := self.service.Delete(commentId); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}

func (self *HttpApi) ManagePage(ctx fiber.Ctx) error {
	query := &usercase.CommentQueryForm{}
	if err := ctx.Bind().JSON(query); err != nil {
		return err
	}
	page, err := self.service.ManagePage(query)
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(page))
}
