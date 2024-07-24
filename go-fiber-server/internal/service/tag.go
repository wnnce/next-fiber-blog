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
	"strconv"
	"strings"
	"time"
)

const tagListCacheKey = "BLOG:tag:list"

type TagService struct {
	repo          usercase.ITagRepo
	redisTemplate *data.RedisTemplate
}

func NewTagService(repo usercase.ITagRepo, redisTemplate *data.RedisTemplate) usercase.ITagService {
	return &TagService{
		repo:          repo,
		redisTemplate: redisTemplate,
	}
}

func (self *TagService) CreateTag(form *usercase.TagForm) error {
	if err := self.checkTagNameIsActive(form.TagName, 0); err != nil {
		return err
	}
	if err := self.repo.Save(form); err != nil {
		slog.Error(fmt.Sprintf("保存标签失败，错误信息：%v", err))
		return tools.FiberServerError("标签新增失败")
	}
	self.deleteRedisTags()
	return nil
}

func (self *TagService) UpdateTag(form *usercase.TagForm) error {
	if err := self.checkTagNameIsActive(form.TagName, form.TagId); err != nil {
		return err
	}
	if err := self.repo.Update(form); err != nil {
		slog.Error(fmt.Sprintf("更新标签失败，错误信息：%v", err))
		return tools.FiberServerError("标签更新失败")
	}
	self.deleteRedisTags()
	return nil
}

func (self *TagService) UpdateSelectiveTag(form *usercase.TagUpdateForm) error {
	if err := self.repo.UpdateSelective(form); err != nil {
		slog.Error("快捷更新标签失败", "error", err.Error())
		return tools.FiberServerError("更新失败")
	}
	self.deleteRedisTags()
	return nil
}

func (self *TagService) QueryTagInfo(tagId int) (*usercase.Tag, error) {
	tag, err := self.repo.SelectById(tagId)
	if err != nil {
		slog.Error(fmt.Sprintf("获取标签详情失败，错误信息：%v", err))
		return nil, tools.FiberServerError("获取标签失败")
	}
	if tag == nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "标签不存在")
	}
	// 异步更新查看次数
	go func() {
		_ = self.repo.UpdateViewNum(tagId, 1)
	}()
	return tag, nil
}

func (self *TagService) PageTag(form *usercase.TagQueryForm) (*usercase.PageData[usercase.Tag], error) {
	tags, total, err := self.repo.Page(form)
	if err != nil {
		slog.Error(fmt.Sprintf("获取标签列表失败，错误信息：%v", err))
		return nil, tools.FiberServerError("获取标签列表失败")
	}
	pages := int(math.Ceil(float64(total) / float64(form.Size)))
	return &usercase.PageData[usercase.Tag]{
		Current: form.Page,
		Size:    form.Size,
		Pages:   pages,
		Total:   total,
		Records: tags,
	}, nil
}

func (self *TagService) AllTag() []*usercase.Tag {
	tags, err := data.RedisGetSlice[*usercase.Tag](context.Background(), tagListCacheKey, self.redisTemplate.Client())
	if err == nil && len(tags) > 0 {
		return tags
	}
	tags, err = self.repo.List()
	if err != nil {
		slog.Error("获取标签列表失败，错误信息：" + err.Error())
		return make([]*usercase.Tag, 0)
	}
	pool.Go(func() {
		if setErr := self.redisTemplate.Set(context.Background(), tagListCacheKey, tags, time.Duration(math.MaxInt64)); err != nil {
			slog.Error("标签列表添加redis缓存失败", "error", setErr.Error())
		}
	})
	return tags
}

func (self *TagService) Delete(tagId int) error {
	// TODO 删除前验证是否还存在文章
	if err := self.repo.DeleteById(tagId); err != nil {
		slog.Error(fmt.Sprintf("删除标签失败，tagId:%d,错误信息：%v", tagId, err))
		return tools.FiberServerError("删除标签失败")
	}
	self.deleteRedisTags()
	return nil
}

func (self *TagService) BatchDelete(ids string) error {
	idList := make([]int, 0)
	if strings.Contains(ids, ",") {
		for _, idStr := range strings.Split(ids, ",") {
			id, err := strconv.ParseInt(idStr, 10, 0)
			if err != nil {
				return tools.FiberRequestError("参数错误")
			}
			idList = append(idList, int(id))
		}
	} else {
		id, err := strconv.ParseInt(ids, 10, 0)
		if err != nil {
			return tools.FiberRequestError("参数错误")
		}
		idList = append(idList, int(id))
	}
	// TODO 删除前验证是否还存在文章
	result, err := self.repo.DeleteByIds(idList)
	if err != nil {
		slog.Error(fmt.Sprintf("批量删除标签失败，tagIdLen:%d,错误信息：%v", len(ids), err))
		return tools.FiberServerError("批量删除标签失败")
	}
	slog.Info(fmt.Sprintf("批量删除标签完成，tagIdLen:%d,row:%d", len(ids), result))
	self.deleteRedisTags()
	return nil
}

func (self *TagService) checkTagNameIsActive(name string, tagId uint) error {
	num, err := self.repo.CountByTagName(name, tagId)
	if err != nil {
		slog.Error(fmt.Sprintf("检查标签名称是否可用失败，错误信息：%v", err))
		return tools.FiberServerError("标签新增失败")
	}
	if num > 0 {
		return tools.FiberRequestError("标签名称已经存在")
	}
	return nil
}

func (self *TagService) deleteRedisTags() {
	pool.Go(func() {
		if err := self.redisTemplate.Delete(context.Background(), tagListCacheKey); err != nil {
			slog.Error("删除标签Redis缓存失败", "error", err.Error())
		}
	})
}
