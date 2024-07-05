package service

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
	"go-fiber-ent-web-layout/internal/tools"
	"go-fiber-ent-web-layout/internal/usercase"
	"log/slog"
	"math"
	"strconv"
	"strings"
)

type TagService struct {
	repo usercase.ITagRepo
}

func NewTagService(repo usercase.ITagRepo) usercase.ITagService {
	return &TagService{
		repo: repo,
	}
}

func (t *TagService) CreateTag(form *usercase.TagForm) error {
	if err := t.checkTagNameIsActive(form.TagName, 0); err != nil {
		return err
	}
	if err := t.repo.Save(form); err != nil {
		slog.Error(fmt.Sprintf("保存标签失败，错误信息：%v", err))
		return tools.FiberServerError("标签新增失败")
	}
	return nil
}

func (t *TagService) UpdateTag(form *usercase.TagForm) error {
	if err := t.checkTagNameIsActive(form.TagName, form.TagId); err != nil {
		return err
	}
	if err := t.repo.Update(form); err != nil {
		slog.Error(fmt.Sprintf("更新标签失败，错误信息：%v", err))
		return tools.FiberServerError("标签更新失败")
	}
	return nil
}

func (t *TagService) QueryTagInfo(tagId int) (*usercase.Tag, error) {
	tag, err := t.repo.SelectById(tagId)
	if err != nil {
		slog.Error(fmt.Sprintf("获取标签详情失败，错误信息：%v", err))
		return nil, tools.FiberServerError("获取标签失败")
	}
	if tag == nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "标签不存在")
	}
	// 异步更新查看次数
	go func() {
		_ = t.repo.UpdateViewNum(tagId, 1)
	}()
	return tag, nil
}

func (t *TagService) PageTag(form *usercase.TagQueryForm) (*usercase.PageData[usercase.Tag], error) {
	tags, total, err := t.repo.Page(form)
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

func (t *TagService) AllTag() []*usercase.Tag {
	list, err := t.repo.List()
	if err != nil {
		slog.Error("获取标签列表失败，错误信息：" + err.Error())
		return make([]*usercase.Tag, 0)
	}
	return list
}

func (t *TagService) Delete(tagId int) error {
	if err := t.repo.DeleteById(tagId); err != nil {
		slog.Error(fmt.Sprintf("删除标签失败，tagId:%d,错误信息：%v", tagId, err))
		return tools.FiberServerError("删除标签失败")
	}
	return nil
}

func (t *TagService) BatchDelete(ids string) error {
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
	result, err := t.repo.DeleteByIds(idList)
	if err != nil {
		slog.Error(fmt.Sprintf("批量删除标签失败，tagIdLen:%d,错误信息：%v", len(ids), err))
		return tools.FiberServerError("批量删除标签失败")
	}
	slog.Info(fmt.Sprintf("批量删除标签完成，tagIdLen:%d,row:%d", len(ids), result))
	return nil
}

func (t *TagService) checkTagNameIsActive(name string, tagId uint) error {
	num, err := t.repo.CountByTagName(name, tagId)
	if err != nil {
		slog.Error(fmt.Sprintf("检查标签名称是否可用失败，错误信息：%v", err))
		return tools.FiberServerError("标签新增失败")
	}
	if num > 0 {
		return tools.FiberRequestError("标签名称已经存在")
	}
	return nil
}
