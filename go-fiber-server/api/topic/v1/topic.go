package topic

import (
	"github.com/gofiber/fiber/v3"
	"go-fiber-ent-web-layout/internal/tools/res"
	"go-fiber-ent-web-layout/internal/usercase"
)

type HttpApi struct {
	service usercase.ITopicService
}

func NewHttpApi(service usercase.ITopicService) *HttpApi {
	return &HttpApi{
		service: service,
	}
}

func (self *HttpApi) Save(ctx fiber.Ctx) error {
	topic := &usercase.Topic{}
	if err := ctx.Bind().JSON(topic); err != nil {
		return err
	}
	if err := self.service.SaveTopic(topic); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}

func (self *HttpApi) Update(ctx fiber.Ctx) error {
	topic := &usercase.Topic{}
	if err := ctx.Bind().JSON(topic); err != nil {
		return err
	}
	if err := self.service.UpdateTopic(topic); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}

func (self *HttpApi) UpdateSelective(ctx fiber.Ctx) error {
	form := &usercase.TopicUpdateForm{}
	if err := ctx.Bind().JSON(form); err != nil {
		return err
	}
	if err := self.service.UpdateSelectiveTopic(form); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}

func (self *HttpApi) Page(ctx fiber.Ctx) error {
	query := &usercase.TopicQueryForm{}
	if err := ctx.Bind().JSON(query); err != nil {
		return err
	}
	page, err := self.service.PageTopic(query)
	if err != nil {
		return err
	}
	return ctx.JSON(res.OkByData(page))
}

func (self *HttpApi) Delete(ctx fiber.Ctx) error {
	topicId := fiber.Params[int64](ctx, "id")
	if err := self.service.DeleteTopicById(topicId); err != nil {
		return err
	}
	return ctx.JSON(res.SimpleOK())
}
