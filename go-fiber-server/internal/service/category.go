package service

import (
	"fmt"
	"go-fiber-ent-web-layout/internal/tools"
	"go-fiber-ent-web-layout/internal/usercase"
	"log/slog"
)

type CategoryService struct {
	repo usercase.ICategoryRepo
}

func NewCategoryService(repo usercase.ICategoryRepo) usercase.ICategoryService {
	return &CategoryService{
		repo: repo,
	}
}

func (c *CategoryService) CreateCategory(cat *usercase.Category) error {
	if err := c.repo.Save(cat); err != nil {
		slog.Error(fmt.Sprintf("分类保存失败，错误信息：%s", err))
		return tools.FiberServerError("分类保存失败")
	}
	return nil
}

func (c *CategoryService) UpdateCategory(cat *usercase.Category) error {
	if cat.CategoryId == 0 {
		return tools.FiberRequestError("分类Id不能为空")
	}
	if err := c.repo.Update(cat); err != nil {
		slog.Error(fmt.Sprintf("分类更新失败，错误信息：%s", err))
		return tools.FiberRequestError("分类更新失败")
	}
	return nil
}

func (c *CategoryService) ListCategory() ([]*usercase.Category, error) {
	categorys := c.repo.List()
	if len(categorys) <= 1 {
		return categorys, nil
	}
	return tools.BuilderTree(categorys), nil
}

func (c *CategoryService) ManageListCategory() ([]*usercase.Category, error) {
	categorys, err := c.repo.ManageList()
	if err != nil {
		slog.Error(fmt.Sprintf("获取分类列表失败，错误信息：%s", err))
		return nil, tools.FiberServerError("分类查询失败")
	}
	if len(categorys) <= 1 {
		return categorys, nil
	}
	return tools.BuilderTree(categorys), nil
}

func (c *CategoryService) QueryCategoryInfo(catId int) (*usercase.Category, error) {
	category, err := c.repo.SelectById(catId)
	if err != nil {
		slog.Error(fmt.Sprintf("分类获取失败，错误信息：%s", err))
		return nil, tools.FiberServerError("分类获取失败")
	}
	return category, nil
}

func (CategoryService) Delete(catId int) error {
	//TODO implement me
	panic("implement me")
}
