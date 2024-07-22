package service

import (
	"context"
	"go-fiber-ent-web-layout/internal/data"
	"go-fiber-ent-web-layout/internal/tools"
	"go-fiber-ent-web-layout/internal/tools/pool"
	"go-fiber-ent-web-layout/internal/usercase"
	"log/slog"
	"math"
	"strconv"
	"time"
)

const noticeCachePrefix = "SYSTEM_NOTICE:"

type NoticeService struct {
	repo          usercase.INoticeRepo
	redisTemplate *data.RedisTemplate
}

func NewNoticeService(repo usercase.INoticeRepo, redisTemplate *data.RedisTemplate) usercase.INoticeService {
	return &NoticeService{
		repo:          repo,
		redisTemplate: redisTemplate,
	}
}

func (self *NoticeService) SaveNotice(notice *usercase.Notice) error {
	if err := self.repo.Save(notice); err != nil {
		slog.Error("保存系统通知失败", "error", err.Error())
		return tools.FiberServerError("保存失败")
	}
	self.deleteRedisNoticeByType(notice.NoticeType)
	return nil
}

func (self *NoticeService) UpdateNotice(notice *usercase.Notice) error {
	if err := self.repo.Update(notice); err != nil {
		slog.Error("更新系统通知失败", "error", err.Error())
		return tools.FiberServerError("更新失败")
	}
	self.deleteRedisNoticeByType(notice.NoticeType)
	return nil
}

func (self *NoticeService) Page(query *usercase.NoticeQueryForm) (*usercase.PageData[usercase.Notice], error) {
	records, total, err := self.repo.ManagePage(query)
	if err != nil {
		slog.Error("分页查询系统通知失败", "error", err.Error())
		return nil, tools.FiberServerError("查询失败")
	}
	return usercase.NewPageData(records, total, query.Page, query.Size), nil
}

func (self *NoticeService) ListNoticeByType(noticeType int) ([]usercase.Notice, error) {
	cacheKey := noticeCachePrefix + "list:" + strconv.Itoa(noticeType)
	notices, err := data.RedisGetSlice[usercase.Notice](context.Background(), cacheKey, self.redisTemplate.Client())
	if err == nil && len(notices) > 0 {
		return notices, nil
	}
	notices, err = self.repo.ListByType(noticeType)
	if err != nil {
		slog.Error("查询指定类型通知列表失败", "error", err.Error())
		return nil, tools.FiberServerError("查询失败")
	}
	pool.Go(func() {
		if setErr := self.redisTemplate.Set(context.Background(), cacheKey, notices, time.Duration(math.MaxInt64)); setErr != nil {
			slog.Error("缓存系统通知列表失败", "error", setErr.Error())
		}
	})
	return notices, err
}

func (self *NoticeService) Delete(noticeId int64) error {
	if err := self.repo.DeleteById(noticeId); err != nil {
		slog.Error("删除系统通知失败", "error", err.Error(), "noticeId", noticeId)
		return tools.FiberServerError("删除失败")
	}
	pool.Go(func() {
		if noticeType := self.repo.QueryNoticeTypeById(noticeId); noticeType >= 0 {
			self.deleteRedisNoticeByType(noticeType)
		}
	})
	return nil
}

func (self *NoticeService) deleteRedisNoticeByType(noticeType int) {
	cacheKey := noticeCachePrefix + "list:" + strconv.Itoa(noticeType)
	pool.Go(func() {
		if err := self.redisTemplate.Delete(context.Background(), cacheKey); err != nil {
			slog.Error("异步删除系统通知redis缓存失败", "error", err.Error(), "key", cacheKey)
		}
	})
}
