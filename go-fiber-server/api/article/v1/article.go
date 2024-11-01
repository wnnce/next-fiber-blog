package article

import (
	"github.com/gofiber/fiber/v3"
	"go-fiber-ent-web-layout/internal/tools"
	"go-fiber-ent-web-layout/internal/tools/res"
	"go-fiber-ent-web-layout/internal/usercase"
)

type HttpApi struct {
	service usercase.IArticleService
}

func NewHttpApi(service usercase.IArticleService) *HttpApi {
	return &HttpApi{
		service: service,
	}
}

func (self *HttpApi) Save(ctx fiber.Ctx) error {
	article := &usercase.Article{}
	if err := ctx.Bind().JSON(article); err != nil {
		return err
	}
	if err := self.service.SaveArticle(article); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}

func (self *HttpApi) Update(ctx fiber.Ctx) error {
	article := &usercase.Article{}
	if err := ctx.Bind().JSON(article); err != nil {
		return err
	}
	if err := self.service.UpdateArticle(article); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}

func (self *HttpApi) UpdateSelective(ctx fiber.Ctx) error {
	form := &usercase.ArticleUpdateForm{}
	if err := ctx.Bind().JSON(form); err != nil {
		return err
	}
	if err := self.service.UpdateSelectiveArticle(form); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}

func (self *HttpApi) Page(ctx fiber.Ctx) error {
	query := &usercase.ArticleQueryForm{}
	if err := ctx.Bind().JSON(query); err != nil {
		return err
	}
	var statusValue uint8 = 0
	query.Status = &statusValue
	page, err := self.service.Page(query)
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(page))
}

func (self *HttpApi) PageByLabel(ctx fiber.Ctx) error {
	query := &usercase.ArticleQueryForm{}
	if err := ctx.Bind().JSON(query); err != nil {
		return err
	}
	page, err := self.service.PageByLabel(query)
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(page))
}

func (self *HttpApi) ManagePage(ctx fiber.Ctx) error {
	query := &usercase.ArticleQueryForm{}
	query.IsAdmin = true
	if err := ctx.Bind().JSON(query); err != nil {
		return err
	}
	page, err := self.service.Page(query)
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(page))
}

func (self *HttpApi) ListTop(ctx fiber.Ctx) error {
	list, err := self.service.ListTopArticle()
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(list))
}

func (self *HttpApi) ListHot(ctx fiber.Ctx) error {
	list, err := self.service.ListHotArticle()
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(list))
}

func (self *HttpApi) Archives(ctx fiber.Ctx) error {
	archives, err := self.service.Archives()
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(archives))
}

func (self *HttpApi) ManageQueryInfo(ctx fiber.Ctx) error {
	articleId := fiber.Params[uint64](ctx, "id")
	article, err := self.service.SelectById(articleId, true)
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(article))
}

func (self *HttpApi) QueryInfo(ctx fiber.Ctx) error {
	articleId := fiber.Params[uint64](ctx, "id")
	article, err := self.service.SelectById(articleId, false)
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(article))
}

func (self *HttpApi) Delete(ctx fiber.Ctx) error {
	articleId := fiber.Params[uint64](ctx, "id")
	if err := self.service.DeleteArticleById(articleId); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}

func (self *HttpApi) VoteUp(ctx fiber.Ctx) error {
	articleId := fiber.Params[uint64](ctx, "id")
	if articleId == 0 {
		return tools.FiberRequestError("参数错误")
	}
	if err := self.service.ArticleVoteUp(articleId); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}
