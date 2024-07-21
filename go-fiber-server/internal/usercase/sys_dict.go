package usercase

import (
	"context"
	"github.com/jackc/pgx/v5"
)

// SysDict 系统字典
type SysDict struct {
	DictId   uint64 `json:"dictId" db:"dict_id"`                         // 字典ID
	DictName string `json:"dictName" db:"dict_name" validate:"required"` // 字典名称
	DictKey  string `json:"dictKey" db:"dict_key" validate:"required"`   // 字典键值
	Remark   string `json:"remark" db:"remark"`                          // 字典备注
	BaseEntity
}

// SysDictQueryForm 系统字典查询表单
type SysDictQueryForm struct {
	DictName        string `json:"dictName"`
	DictKey         string `json:"dictKey"`
	CreateTimeBegin string `json:"createTimeBegin"`
	CreateTimeEnd   string `json:"createTimeEnd"`
	PageQueryForm
}

// SysDictSelectiveUpdateForm 系统字段快捷更新表单
type SysDictSelectiveUpdateForm struct {
	DictId uint64 `json:"dictId" validate:"required"`
	Status *uint8 `json:"status" validate:"required"`
}

// SysDictValue 系统字典数据
type SysDictValue struct {
	ID      uint64 `json:"id" db:"id"`                                          // 数据ID
	DictId  uint64 `json:"dictId,omitempty" db:"dict_Id" validate:"required"`   // 字典Id
	DictKey string `json:"dictKey,omitempty" db:"dict_key" validate:"required"` // 字典Key
	Label   string `json:"label" db:"label" validate:"required"`                // 数据名称
	Value   string `json:"value" db:"value" validate:"required"`                // 数据值
	Remark  string `json:"remark,omitempty" db:"remark"`                        // 备注
	BaseEntity
}

// SysDictValueQueryForm 系统字典数据查询表单
type SysDictValueQueryForm struct {
	DictId          uint64 `json:"dictId"`
	DictKey         string `json:"dictKey"`
	Label           string `json:"label"`
	CreateTimeBegin string `json:"createTimeBegin"`
	CreateTimeEnd   string `json:"createTimeEnd"`
	PageQueryForm
}

// SysDictValueSelectiveUpdateForm 字典数据快捷更新表单
type SysDictValueSelectiveUpdateForm struct {
	ID     uint64 `json:"id" validate:"required"`
	Status *uint8 `json:"status" validate:"required"`
}

type ISysDictRepo interface {
	Transaction(ctx context.Context, fn func(tx pgx.Tx) error) error

	SaveDict(dict *SysDict) error

	UpdateDict(dict *SysDict, tx pgx.Tx) error

	UpdateSelectiveDict(dict *SysDict, tx pgx.Tx) error

	CountByKey(key string, dictId uint64) (uint8, error)

	PageDict(query *SysDictQueryForm) ([]*SysDict, int64, error)

	SelectDictById(dictId uint64) (*SysDict, error)

	SelectDictKeyById(dictId int64) string

	DeleteDict(dictId int64, tx pgx.Tx) error

	SaveDictValue(value *SysDictValue) error

	UpdateDictValue(value *SysDictValue) error

	UpdateSelectiveDictValue(value *SysDictValue, tx pgx.Tx) error

	UpdateDictValueByDickId(value *SysDictValue, tx pgx.Tx) error

	CountValueById(value string, dictId uint64, valueId uint64) (uint8, error)

	SelectDictKeyByValueId(valueId int64) string

	PageDictValue(query *SysDictValueQueryForm) ([]*SysDictValue, int64, error)

	ListDictValueByDictKey(dictKey string) ([]SysDictValue, error)

	DeleteDictValue(valueId int64) error

	DeleteDictValueByDictId(dictId int64, tx pgx.Tx) error
}

type ISysDictService interface {
	PageDict(query *SysDictQueryForm) (*PageData[SysDict], error)

	SaveDict(dict *SysDict) error

	UpdateDict(dict *SysDict) error

	UpdateSelectiveDict(form *SysDictSelectiveUpdateForm) error

	DeleteDict(dictId int64) error

	PageDictValue(query *SysDictValueQueryForm) (*PageData[SysDictValue], error)

	SaveDictValue(value *SysDictValue) error

	UpdateDictValue(value *SysDictValue) error

	UpdateSelectiveValue(form *SysDictValueSelectiveUpdateForm) error

	DeleteDictValue(valueId int64) error

	ListDictValueByDictKey(dictKey string) ([]SysDictValue, error)
}
