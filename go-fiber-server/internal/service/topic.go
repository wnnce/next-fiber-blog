package service

import (
	"go-fiber-ent-web-layout/internal/tools"
	"go-fiber-ent-web-layout/internal/usercase"
	"log/slog"
)

type TopicService struct {
	repo usercase.ITopicRepo
}

func NewTopicService(repo usercase.ITopicRepo) usercase.ITopicService {
	return &TopicService{
		repo: repo,
	}
}

func (self *TopicService) SaveTopic(topic *usercase.Topic) error {
	if err := self.repo.Save(topic); err != nil {
		slog.Error("保存博客动态失败", "error", err.Error())
		return tools.FiberServerError("保存失败")
	}
	return nil
}

func (self *TopicService) UpdateTopic(topic *usercase.Topic) error {
	if err := self.repo.Update(topic); err != nil {
		slog.Error("更新博客动态失败", "error", err.Error())
		return tools.FiberServerError("更新失败")
	}
	return nil
}

func (self *TopicService) UpdateSelectiveTopic(form *usercase.TopicUpdateForm) error {
	if err := self.repo.UpdateSelective(form); err != nil {
		slog.Error("快捷更新博客动态失败", "error", err.Error())
		return tools.FiberServerError("更新失败")
	}
	return nil
}

func (self *TopicService) PageTopic(query *usercase.TopicQueryForm) (*usercase.PageData[usercase.Topic], error) {
	topics, total, err := self.repo.Page(query)
	if err != nil {
		slog.Error("分页查询博客动态失败", "error", err.Error())
		return nil, tools.FiberServerError("查询失败")
	}
	return usercase.NewPageData[usercase.Topic](topics, total, query.Page, query.Size), nil
}

func (self *TopicService) DeleteTopicById(topicId int64) error {
	if err := self.repo.DeleteById(topicId); err != nil {
		slog.Error("删除博客动态失败", "error", err.Error())
		return tools.FiberServerError("删除失败")
	}
	return nil
}
