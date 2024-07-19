package service

import (
	"context"
	"github.com/jackc/pgx/v5"
	"go-fiber-ent-web-layout/internal/tools"
	"go-fiber-ent-web-layout/internal/usercase"
	"log/slog"
)

type SysDictService struct {
	repo usercase.ISysDictRepo
}

func NewSysDictService(repo usercase.ISysDictRepo) usercase.ISysDictService {
	return &SysDictService{
		repo: repo,
	}
}

func (self *SysDictService) PageDict(query *usercase.SysDictQueryForm) (*usercase.PageData[usercase.SysDict], error) {
	records, total, err := self.repo.PageDict(query)
	if err != nil {
		slog.Error("分页查询系统字典失败", "error", err.Error())
		return nil, tools.FiberServerError("查询失败")
	}
	return usercase.NewPageData(records, total, query.Page, query.Size), nil
}

func (self *SysDictService) SaveDict(dict *usercase.SysDict) error {
	if err := self.checkDictKeyIsActive(dict.DictKey, 0); err != nil {
		return err
	}
	if err := self.repo.SaveDict(dict); err != nil {
		slog.Error("系统字段保存失败", "error", err.Error())
		return tools.FiberServerError("保存失败")
	}
	return nil
}

func (self *SysDictService) UpdateDict(dict *usercase.SysDict) error {
	find, err := self.repo.SelectDictById(dict.DictId)
	if err != nil || find == nil {
		slog.Error("查询数据字典失败", "dictId", dict.DictId)
		return tools.FiberServerError("更新失败")
	}
	transactionErr := self.repo.Transaction(context.Background(), func(tx pgx.Tx) error {
		// 如果key不一致 那么检查key是否重复
		// 如果不重复那么先更新字典数据项的key
		if find.DictKey != dict.DictKey {
			if err = self.checkDictKeyIsActive(dict.DictKey, 0); err != nil {
				return err
			}
			updateDictValue := &usercase.SysDictValue{
				DictKey: dict.DictKey,
				DictId:  dict.DictId,
			}
			if err = self.repo.UpdateDictValueByDickId(updateDictValue, tx); err != nil {
				slog.Error("更新系统字典数据对应的Key失败", "error", err.Error())
				return err
			}
		}
		// 如果系统字典的状态发生了变化那么也需要更新字典数据的状态
		if *find.Status != *dict.Status {
			var valueStatus uint8
			if *dict.Status == 1 {
				valueStatus = 2
			}
			updateDictValue := &usercase.SysDictValue{
				DictId: dict.DictId,
			}
			updateDictValue.Status = &valueStatus
			if err = self.repo.UpdateDictValueByDickId(updateDictValue, tx); err != nil {
				slog.Error("更新系统字典数据状态失败", "error", err.Error())
				return err
			}
		}
		if err = self.repo.UpdateDict(dict, tx); err != nil {
			slog.Info("更新系统字典数据失败", "error", err.Error())
			return err
		}
		return nil
	})
	if transactionErr != nil {
		return tools.FiberServerError("更新失败")
	}
	return nil
}

func (self *SysDictService) UpdateSelectiveDict(form *usercase.SysDictSelectiveUpdateForm) error {
	dict := &usercase.SysDict{
		DictId: form.DictId,
	}
	dict.Status = form.Status
	err := self.repo.Transaction(context.Background(), func(tx pgx.Tx) error {
		var valueStatus uint8
		if *dict.Status == 1 {
			valueStatus = 2
		}
		updateDictValue := &usercase.SysDictValue{
			DictId: dict.DictId,
		}
		updateDictValue.Status = &valueStatus
		if err := self.repo.UpdateDictValueByDickId(updateDictValue, tx); err != nil {
			slog.Error("更新系统字典数据状态失败", "error", err.Error())
			return err
		}
		if err := self.repo.UpdateSelectiveDict(dict, nil); err != nil {
			slog.Error("更新系统字典状态失败", "error", err.Error())
			return err
		}
		return nil
	})
	if err != nil {
		return tools.FiberServerError("更新失败")
	}
	return nil
}

func (self *SysDictService) DeleteDict(dictId int64) error {
	if err := self.repo.DeleteDict(dictId, nil); err != nil {
		slog.Error("删除系统字典失败", "error", err.Error())
		return tools.FiberServerError("删除失败")
	}
	return nil
}

func (self *SysDictService) PageDictValue(query *usercase.SysDictValueQueryForm) (*usercase.PageData[usercase.SysDictValue], error) {
	records, total, err := self.repo.PageDictValue(query)
	if err != nil {
		slog.Error("分页查询系统字典数据失败", "error", err.Error())
	}
	return usercase.NewPageData(records, total, query.Page, query.Size), nil
}

func (self *SysDictService) SaveDictValue(value *usercase.SysDictValue) error {
	if err := self.checkDictValueIsActive(value.Value, value.DictId, 0); err != nil {
		return err
	}
	if err := self.repo.SaveDictValue(value); err != nil {
		slog.Error("保存系统字典数据失败", "error", err.Error())
		return tools.FiberServerError("保存失败")
	}
	return nil
}

func (self *SysDictService) UpdateDictValue(value *usercase.SysDictValue) error {
	if err := self.checkDictValueIsActive(value.Value, value.DictId, 0); err != nil {
		return err
	}
	if err := self.repo.UpdateDictValue(value); err != nil {
		slog.Error("保存系统字典数据失败", "error", err.Error())
		return tools.FiberServerError("保存失败")
	}
	return nil
}

func (self *SysDictService) UpdateSelectiveValue(form *usercase.SysDictValueSelectiveUpdateForm) error {
	dictValue := &usercase.SysDictValue{
		ID: form.ID,
	}
	dictValue.Status = form.Status
	if err := self.repo.UpdateSelectiveDictValue(dictValue, nil); err != nil {
		slog.Error("更新系统字典数据状态失败", "error", err.Error())
		return tools.FiberServerError("保存失败")
	}
	return nil
}

func (self *SysDictService) DeleteDictValue(valueId int64) error {
	if err := self.repo.DeleteDictValue(valueId); err != nil {
		slog.Error("删除系统字典数据失败", "error", err.Error())
		return tools.FiberServerError("删除失败")
	}
	return nil
}

func (self *SysDictService) checkDictKeyIsActive(dictKey string, dictId uint64) error {
	count, err := self.repo.CountByKey(dictKey, dictId)
	if err != nil || count > 0 {
		slog.Error("检查数据字典Key是否可用失败", "dictKey", dictKey, "dictId", dictId)
		return tools.FiberServerError("检查Key是否可用失败")
	}
	return nil
}

func (self *SysDictService) checkDictValueIsActive(value string, dictId uint64, valueId uint64) error {
	count, err := self.repo.CountValueById(value, dictId, valueId)
	if err != nil || count > 0 {
		slog.Error("检查数据字典Key是否可用失败", "dictId", dictId, "valueId", valueId)
		return tools.FiberServerError("检查value是否重复失败")
	}
	return nil
}
