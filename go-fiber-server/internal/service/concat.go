package service

import (
	"fmt"
	"go-fiber-ent-web-layout/internal/tools"
	"go-fiber-ent-web-layout/internal/usercase"
	"log/slog"
)

type ConcatService struct {
	repo usercase.IConcatRepo
}

func NewConcatService(repo usercase.IConcatRepo) usercase.IConcatService {
	return &ConcatService{
		repo: repo,
	}
}

// CreateConcat 新增联系方式
func (c *ConcatService) CreateConcat(concat *usercase.Concat) error {
	if err := c.checkConcatName(concat.Name, 0); err != nil {
		return err
	}
	if err := c.repo.Save(concat); err != nil {
		slog.Error(fmt.Sprintf("保存联系方式失败，错误信息：%s", err))
		return tools.FiberServerError("联系方式保存失败")
	}
	return nil
}

func (c *ConcatService) UpdateConcat(concat *usercase.Concat) error {
	if concat.ConcatId == 0 {
		return tools.FiberRequestError("联系方式Id不能为空")
	}
	if err := c.checkConcatName(concat.Name, concat.ConcatId); err != nil {
		return err
	}
	if err := c.repo.Update(concat); err != nil {
		slog.Error(fmt.Sprintf("更新联系方式失败，错误信息：%s", err))
		return tools.FiberServerError("联系方式更新失败")
	}
	return nil
}

func (c *ConcatService) ListConcat() ([]*usercase.Concat, error) {
	concats, err := c.repo.List()
	if err != nil {
		slog.Error(fmt.Sprintf("获取联系方式列表失败，错误信息：%s", err))
		return nil, tools.FiberServerError("查询失败")
	}
	return concats, nil
}

func (c *ConcatService) ManageListConcat(query *usercase.ConcatQueryForm) ([]*usercase.Concat, error) {
	concats, err := c.repo.ManageList(query)
	if err != nil {
		slog.Error(fmt.Sprintf("获取联系方式列表失败，错误信息：%s", err))
		return nil, tools.FiberServerError("查询失败")
	}
	return concats, nil
}

func (c *ConcatService) Delete(cid int) error {
	if err := c.repo.DeleteById(cid); err != nil {
		slog.Error(fmt.Sprintf("删除联系方式失败，错误信息：%s", err))
		return tools.FiberServerError("删除失败")
	}
	return nil
}

func (c *ConcatService) checkConcatName(name string, cid uint) error {
	total, err := c.repo.CountByName(name, cid)
	if err != nil {
		slog.Error(fmt.Sprintf("检查联系方式名称是否可用失败，错误信息：%s", err))
		return tools.FiberServerError("联系方式保存失败")
	}
	if total > 0 {
		return tools.FiberServerError("联系方式名称已经存在")
	}
	return nil
}
