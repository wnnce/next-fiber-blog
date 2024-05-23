package service

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
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
	if err := c.checkCategoryName(cat.CategoryName, 0); err != nil {
		return err
	}
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
	if err := c.checkCategoryName(cat.CategoryName, cat.CategoryId); err != nil {
		return err
	}
	if err := c.repo.Update(cat); err != nil {
		slog.Error(fmt.Sprintf("分类更新失败，错误信息：%s", err))
		return tools.FiberRequestError("分类更新失败")
	}
	return nil
}

func (c *CategoryService) ListCategory() ([]*usercase.Category, error) {
	categorys, err := c.repo.List()
	if err != nil {
		slog.Error("获取分类列表失败，错误信息：" + err.Error())
		return make([]*usercase.Category, 0), tools.FiberServerError("获取标签列表失败")
	}
	if len(categorys) <= 1 {
		return categorys, nil
	}
	return tools.BuilderTree[uint](categorys), nil
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
	return tools.BuilderTree[uint](categorys), nil
}

func (c *CategoryService) QueryCategoryInfo(catId int) (*usercase.Category, error) {
	category, err := c.repo.SelectById(catId)
	if err != nil {
		slog.Error(fmt.Sprintf("分类获取失败，错误信息：%s", err))
		return nil, tools.FiberServerError("分类获取失败")
	}
	if category == nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "分类不存在")
	}
	return category, nil
}

func (c *CategoryService) Delete(catId int) error {
	if err := c.repo.DeleteById(catId); err != nil {
		slog.Error(fmt.Sprintf("分类删除失败，错误信息：%s", err))
		return tools.FiberServerError("分类删除失败")
	}
	return nil
}

func (c *CategoryService) checkCategoryName(name string, catId uint) error {
	total, err := c.repo.CountByName(name, catId)
	if err != nil {
		slog.Error(fmt.Sprintf("检查分类名称是否可用失败，错误信息：%s", err))
		return tools.FiberServerError("新增分类失败")
	}
	if total > 0 {
		return tools.FiberServerError("标签名称已经存在")
	}
	return nil
}
