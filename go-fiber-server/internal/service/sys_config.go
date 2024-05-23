package service

import (
	"go-fiber-ent-web-layout/internal/tools"
	"go-fiber-ent-web-layout/internal/usercase"
	"log/slog"
	"math"
)

type SysConfigService struct {
	repo usercase.ISysConfigRepo
}

func NewSysConfigService(repo usercase.ISysConfigRepo) usercase.ISysConfigService {
	return &SysConfigService{
		repo: repo,
	}
}

func (ss *SysConfigService) CreateConfig(cfg *usercase.SysConfig) error {
	num, err := ss.repo.CountByKey(cfg.ConfigKey, 0)
	if err != nil {
		slog.Error("查询configKey数量失败", "err", err.Error())
		return tools.FiberServerError("保存失败")
	}
	if num > 0 {
		slog.Info("系统配置Key已存在", "key", cfg.ConfigKey)
		return tools.FiberRequestError("该configKey已存在")
	}
	if err = ss.repo.Save(cfg); err != nil {
		slog.Error("保存系统配置信息失败", "err", err.Error())
		return tools.FiberServerError("保存失败")
	}
	return nil
}

func (ss *SysConfigService) UpdateConfig(cfg *usercase.SysConfig) error {
	num, err := ss.repo.CountByKey(cfg.ConfigKey, cfg.ConfigId)
	if err != nil {
		slog.Error("查询configKey数量失败", "err", err.Error())
		return tools.FiberServerError("保存失败")
	}
	if num > 0 {
		slog.Info("系统配置Key已存在", "key", cfg.ConfigKey)
		return tools.FiberRequestError("该configKey已存在")
	}
	if err = ss.repo.Update(cfg); err != nil {
		slog.Error("更新系统配置信息失败", "err", err.Error())
		return tools.FiberServerError("更新失败")
	}
	return nil
}

func (ss *SysConfigService) ManageList(query *usercase.SysConfigQueryForm) (*usercase.PageData[usercase.SysConfig], error) {
	list, total, err := ss.repo.ManagePage(query)
	if err != nil {
		slog.Error("获取系统配置分页列表失败", "err", err.Error())
		return nil, tools.FiberServerError("查询失败")
	}
	pages := int(math.Ceil(float64(total) / float64(query.Size)))
	records := make([]*usercase.SysConfig, 0, len(list))
	for i := 0; i < len(list); i++ {
		records = append(records, &list[i])
	}
	return &usercase.PageData[usercase.SysConfig]{
		Current: query.Page,
		Size:    query.Size,
		Pages:   pages,
		Total:   total,
		Records: records,
	}, nil
}

func (ss *SysConfigService) Delete(cid int64) error {
	if err := ss.repo.DeleteById(cid); err != nil {
		slog.Error("删除系统配置失败", "err", err.Error())
		return tools.FiberServerError("删除失败")
	}
	return nil
}
