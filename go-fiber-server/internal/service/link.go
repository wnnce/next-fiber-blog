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

func (l *LinkService) CreateLike(link *usercase.Link) error {
	if err := l.repo.Save(link); err != nil {
		slog.Error(fmt.Sprintf("新增友情链接失败，错误信息：%s", err))
		return tools.FiberServerError("保存失败")
	}
	return nil
}

func (l *LinkService) UpdateLike(like *usercase.Link) error {
	if err := l.repo.Update(like); err != nil {
		slog.Error(fmt.Sprintf("更新友情链接失败，错误信息：%s", err))
		return tools.FiberServerError("更新失败")
	}
	return nil
}

func (l *LinkService) PageLike(query *usercase.PageQueryForm) (*usercase.PageData[usercase.Link], error) {
	links, total, err := l.repo.Page(query)
	if err != nil {
		slog.Error(fmt.Sprintf("分页查询友情链接失败，错误信息：%s", err))
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

func (l *LinkService) ManagePageLike(query *usercase.LinkQueryForm) (*usercase.PageData[usercase.Link], error) {
	links, total, err := l.repo.ManagePage(query)
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

func (l *LinkService) Delete(linkId int64) error {
	err := l.repo.DeleteById(linkId)
	if err != nil {
		slog.Error(fmt.Sprintf("删除友情链接失败，错误信息：%s", err))
		return tools.FiberServerError("删除失败")
	}
	return nil
}
