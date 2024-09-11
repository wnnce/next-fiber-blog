package service

import (
	"context"
	"fmt"
	"go-fiber-ent-web-layout/internal/data"
	"go-fiber-ent-web-layout/internal/tools"
	"go-fiber-ent-web-layout/internal/usercase"
	"go-fiber-ent-web-layout/pkg/pool"
	"log/slog"
	"math"
	"time"
)

const concatListCacheKey = "BLOG:concat:list"

type ConcatService struct {
	repo          usercase.IConcatRepo
	redisTemplate *data.RedisTemplate
}

func NewConcatService(repo usercase.IConcatRepo, redisTemplate *data.RedisTemplate) usercase.IConcatService {
	return &ConcatService{
		repo:          repo,
		redisTemplate: redisTemplate,
	}
}

// CreateConcat 新增联系方式
func (self *ConcatService) CreateConcat(concat *usercase.Concat) error {
	if err := self.checkConcatName(concat.Name, 0); err != nil {
		return err
	}
	if err := self.repo.Save(concat); err != nil {
		slog.Error(fmt.Sprintf("保存联系方式失败，错误信息：%s", err))
		return tools.FiberServerError("联系方式保存失败")
	}
	self.deleteRedisConcat()
	return nil
}

func (self *ConcatService) UpdateConcat(concat *usercase.Concat) error {
	if concat.ConcatId == 0 {
		return tools.FiberRequestError("联系方式Id不能为空")
	}
	if err := self.checkConcatName(concat.Name, concat.ConcatId); err != nil {
		return err
	}
	if err := self.repo.Update(concat); err != nil {
		slog.Error(fmt.Sprintf("更新联系方式失败，错误信息：%s", err))
		return tools.FiberServerError("联系方式更新失败")
	}
	self.deleteRedisConcat()
	return nil
}

func (self *ConcatService) UpdateSelectiveConcat(form *usercase.ConcatUpdateForm) error {
	if err := self.repo.UpdateSelective(form); err != nil {
		slog.Error("快捷更新联系方式失败", "error", err.Error())
		return tools.FiberServerError("更新失败")
	}
	self.deleteRedisConcat()
	return nil
}

func (self *ConcatService) ListConcat() ([]*usercase.Concat, error) {
	concats, err := data.RedisGetSlice[*usercase.Concat](context.Background(), concatListCacheKey)
	if err == nil && len(concats) > 0 {
		return concats, nil
	}
	concats, err = self.repo.List()
	if err != nil {
		slog.Error(fmt.Sprintf("获取联系方式列表失败，错误信息：%s", err))
		return nil, tools.FiberServerError("查询失败")
	}
	pool.Go(func() {
		if setErr := self.redisTemplate.Set(context.Background(), concatListCacheKey, concats, time.Duration(math.MaxInt64)); setErr != nil {
			slog.Error("添加联系方式列表Redis缓存失败", "error", err.Error())
		}
	})
	return concats, nil
}

func (self *ConcatService) ManageListConcat(query *usercase.ConcatQueryForm) ([]*usercase.Concat, error) {
	concats, err := self.repo.ManageList(query)
	if err != nil {
		slog.Error(fmt.Sprintf("获取联系方式列表失败，错误信息：%s", err))
		return nil, tools.FiberServerError("查询失败")
	}
	return concats, nil
}

func (self *ConcatService) Delete(cid int) error {
	if err := self.repo.DeleteById(cid); err != nil {
		slog.Error(fmt.Sprintf("删除联系方式失败，错误信息：%s", err))
		return tools.FiberServerError("删除失败")
	}
	self.deleteRedisConcat()
	return nil
}

func (self *ConcatService) checkConcatName(name string, cid uint) error {
	total, err := self.repo.CountByName(name, cid)
	if err != nil {
		slog.Error(fmt.Sprintf("检查联系方式名称是否可用失败，错误信息：%s", err))
		return tools.FiberServerError("联系方式保存失败")
	}
	if total > 0 {
		return tools.FiberRequestError("联系方式名称已经存在")
	}
	return nil
}

func (self *ConcatService) deleteRedisConcat() {
	pool.Go(func() {
		if err := self.redisTemplate.Delete(context.Background(), concatListCacheKey); err != nil {
			slog.Error("删除联系方式Redis缓存失败", "error", err.Error())
		}
	})
}
