package usercase

// SysDict 系统字典
type SysDict struct {
	DictId   uint64 `json:"dictId" db:"dict_id"`     // 字典ID
	DictName string `json:"dictName" db:"dict_name"` // 字典名称
	DictKey  string `json:"dictKey" db:"dict_key"`   // 字典键值
	Remark   string `json:"remark" db:"remark"`      // 字典备注
	BaseEntity
}

// SysDictQueryForm 系统字典查询表单
type SysDictQueryForm struct {
	DictName string `json:"dictName"`
	DictKey  string `json:"dictKey"`
	PageQueryForm
}

// SysDictValue 系统字典数据
type SysDictValue struct {
	ID      uint64 `json:"id" db:"id"`            // 数据ID
	DictId  uint64 `json:"dictId" db:"dict_Id"`   // 字典Id
	DictKey string `json:"dictKey" db:"dict_key"` // 字典Key
	Label   string `json:"label" db:"label"`      // 数据名称
	Value   string `json:"value" db:"value"`      // 数据值
	Remark  string `json:"remark" db:"remark"`    // 备注
	BaseEntity
}

// SysDictValueQueryForm 系统字典数据查询表单
type SysDictValueQueryForm struct {
	DictId  uint64 `json:"dictId"`
	DictKey string `json:"dictKey"`
	Label   string `json:"label"`
	PageQueryForm
}

type ISysDictRepo interface {
	SaveDict(dict *SysDict) error

	UpdateDict(dict *SysDict) error

	PageDict(query *SysDictQueryForm) ([]*SysDict, int64, error)

	DeleteDict(dictId int64) error

	SaveDictValue(value *SysDictValue) error

	UpdateDictValue(value *SysDictValue) error

	PageDictValue(query *SysDictValueQueryForm) ([]*SysDictValue, int64, error)

	DeleteDictValue(valueId int64) error
}

type ISysDictService interface {
}
