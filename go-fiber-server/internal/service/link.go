package service

import (
	"fmt"
	"go-fiber-ent-web-layout/internal/tools"
	"go-fiber-ent-web-layout/internal/usercase"
	"log/slog"
	"math"
)

type LinkService struct {
	repo usercase.ILinkRepo
}

func NewLinkService(repo usercase.ILinkRepo) usercase.ILinkService {
	return &LinkService{
		repo: repo,
	}
}

func (self *LinkService) CreateLink(link *usercase.Link) error {
	if err := self.repo.Save(link); err != nil {
		slog.Error(fmt.Sprintf("新增友情链接失败，错误信息：%s", err))
		return tools.FiberServerError("保存失败")
	}
	return nil
}

func (self *LinkService) UpdateLink(link *usercase.Link) error {
	if err := self.repo.Update(link); err != nil {
		slog.Error(fmt.Sprintf("更新友情链接失败，错误信息：%s", err))
		return tools.FiberServerError("更新失败")
	}
	return nil
}

func (self *LinkService) UpdateSelectiveLink(form *usercase.LinkUpdateForm) error {
	if err := self.repo.UpdateSelective(form); err != nil {
		slog.Error("快捷更新友情链接失败", "error", err.Error())
		return tools.FiberServerError("更新失败")
	}
	return nil
}

func (self *LinkService) List() ([]*usercase.Link, error) {
	links, err := self.repo.List()
	if err != nil {
		slog.Error("获取友情链接列表失败", "error", err.Error())
		return nil, tools.FiberServerError("查询失败")
	}
	return links, nil
}

func (self *LinkService) ManagePageLink(query *usercase.LinkQueryForm) (*usercase.PageData[usercase.Link], error) {
	links, total, err := self.repo.ManagePage(query)
	if err != nil {
		slog.Error(fmt.Sprintf("管理端分页查询友情链接失败，错误信息：%s", err))
		return nil, tools.FiberServerError("查询失败")
	}
	pages := int(math.Ceil(float64(total) / float64(query.Size)))
	return &usercase.PageData[usercase.Link]{
		Current: query.Page,
		Pages:   pages,
		Size:    query.Size,
		Total:   total,
		Records: links,
	}, nil
}

func (self *LinkService) Delete(linkId int64) error {
	err := self.repo.DeleteById(linkId)
	if err != nil {
		slog.Error(fmt.Sprintf("删除友情链接失败，错误信息：%s", err))
		return tools.FiberServerError("删除失败")
	}
	return nil
}
