package service

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"go-fiber-ent-web-layout/internal/data"
	"go-fiber-ent-web-layout/internal/tools"
	"go-fiber-ent-web-layout/internal/usercase"
	"go-fiber-ent-web-layout/pkg/pool"
	"log/slog"
	"math"
)

const categoryTreeListCacheKey = "BLOG:category:list:tree"

type CategoryService struct {
	repo          usercase.ICategoryRepo
	redisTemplate *data.RedisTemplate
}

func NewCategoryService(repo usercase.ICategoryRepo, redisTemplate *data.RedisTemplate) usercase.ICategoryService {
	return &CategoryService{
		repo:          repo,
		redisTemplate: redisTemplate,
	}
}

func (self *CategoryService) CreateCategory(cat *usercase.Category) error {
	if err := self.checkCategoryName(cat.CategoryName, 0); err != nil {
		return err
	}
	if err := self.repo.Save(cat); err != nil {
		slog.Error(fmt.Sprintf("分类保存失败，错误信息：%s", err))
		return tools.FiberServerError("分类保存失败")
	}
	self.deleteRedisCategory()
	return nil
}

func (self *CategoryService) UpdateCategory(cat *usercase.Category) error {
	if cat.CategoryId == 0 {
		return tools.FiberRequestError("分类Id不能为空")
	}
	if err := self.checkCategoryName(cat.CategoryName, cat.CategoryId); err != nil {
		return err
	}
	if err := self.repo.Update(cat); err != nil {
		slog.Error(fmt.Sprintf("分类更新失败，错误信息：%s", err))
		return tools.FiberRequestError("分类更新失败")
	}
	self.deleteRedisCategory()
	return nil
}

func (self *CategoryService) UpdateSelectiveCategory(form *usercase.CategoryUpdateForm) error {
	if err := self.repo.UpdateSelective(form); err != nil {
		slog.Error("快捷更新分类失败", "error", err.Error())
		return tools.FiberRequestError("更新失败")
	}
	self.deleteRedisCategory()
	return nil
}

func (self *CategoryService) ListCategory() ([]*usercase.Category, error) {
	categoryTree, err := data.RedisGetSlice[*usercase.Category](context.Background(), categoryTreeListCacheKey, self.redisTemplate.Client())
	if err == nil && len(categoryTree) > 0 {
		return categoryTree, nil
	}
	categorys, err := self.repo.List()
	if err != nil {
		slog.Error("获取分类列表失败，错误信息：" + err.Error())
		return nil, tools.FiberServerError("获取标签列表失败")
	}
	categoryTree = tools.BuilderTree[uint](categorys)
	pool.Go(func() {
		if setErr := self.redisTemplate.Set(context.Background(), categoryTreeListCacheKey, categoryTree, math.MaxInt64); err != nil {
			slog.Error("分类树形列表添加redis缓存失败", "error", setErr.Error())
		}
	})
	return categoryTree, err
}

func (self *CategoryService) ManageListCategory() ([]*usercase.Category, error) {
	categorys, err := self.repo.ManageList()
	if err != nil {
		slog.Error(fmt.Sprintf("获取分类列表失败，错误信息：%s", err))
		return nil, tools.FiberServerError("分类查询失败")
	}
	if len(categorys) <= 1 {
		return categorys, nil
	}
	return tools.BuilderTree[uint](categorys), nil
}

func (self *CategoryService) QueryCategoryInfo(catId int) (*usercase.Category, error) {
	category, err := self.repo.SelectById(catId)
	if err != nil {
		slog.Error(fmt.Sprintf("分类获取失败，错误信息：%s", err))
		return nil, tools.FiberServerError("分类获取失败")
	}
	if category == nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "分类不存在")
	}
	return category, nil
}

func (self *CategoryService) Delete(catId int) error {
	if err := self.repo.DeleteById(catId); err != nil {
		slog.Error(fmt.Sprintf("分类删除失败，错误信息：%s", err))
		return tools.FiberServerError("分类删除失败")
	}
	self.deleteRedisCategory()
	return nil
}

func (self *CategoryService) checkCategoryName(name string, catId uint) error {
	total, err := self.repo.CountByName(name, catId)
	if err != nil {
		slog.Error(fmt.Sprintf("检查分类名称是否可用失败，错误信息：%s", err))
		return tools.FiberServerError("新增分类失败")
	}
	if total > 0 {
		return tools.FiberServerError("标签名称已经存在")
	}
	return nil
}

// 将分类缓存从redis中删除 使用异步删除
func (self *CategoryService) deleteRedisCategory() {
	pool.Go(func() {
		if err := self.redisTemplate.Delete(context.Background(), categoryTreeListCacheKey); err != nil {
			slog.Error("删除分类Redis缓存失败", "error", err.Error())
		}
	})
}
