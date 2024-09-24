package usercase

import "time"

// SysConfig 系统配置缓存表
type SysConfig struct {
	ConfigId    uint64     `json:"configId" db:"config_id"`                                          // 配置Id
	ConfigName  string     `json:"configName,omitempty" db:"config_name" validate:"required,max=64"` // 配置名称
	ConfigKey   string     `json:"configKey,omitempty" db:"config_key" validate:"required,max=64"`   // 缓存Key
	ConfigValue string     `json:"configValue,omitempty" db:"config_value" validate:"required"`      // 缓存Value
	DeleteAt    string     `json:"deleteAt,omitempty" db:"delete_at"`                                // 删除标志
	CreateTime  *time.Time `json:"createTime,omitempty" db:"create_time"`                            // 创建时间
	UpdateTime  *time.Time `json:"updateTime,omitempty" db:"update_time"`                            // 最后更新时间
	Remark      string     `json:"remark,omitempty" db:"remark"`                                     // 备注
}

// SysConfigQueryForm 查询表单
type SysConfigQueryForm struct {
	Name            string `json:"name"`
	Key             string `json:"key"`
	CreateTimeBegin string `json:"createTimeBegin"`
	CreateTimeEnd   string `json:"createTimeEnd"`
	PageQueryForm
}

// ISysConfigRepo 系统配置Repo层接口
type ISysConfigRepo interface {
	Save(cfg *SysConfig) error

	Update(cfg *SysConfig) error

	CountByKey(key string, cid uint64) (uint8, error)

	ManagePage(query *SysConfigQueryForm) ([]SysConfig, int64, error)

	DeleteById(cid int64) error
}

// ISysConfigService 系统配置Service层接口
type ISysConfigService interface {
	CreateConfig(cfg *SysConfig) error

	UpdateConfig(cfg *SysConfig) error

	ManageList(query *SysConfigQueryForm) (*PageData[SysConfig], error)

	Delete(cid int64) error
}
