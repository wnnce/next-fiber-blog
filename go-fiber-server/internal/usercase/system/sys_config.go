package usercase

import "time"

// SysConfig 系统配置缓存表
type SysConfig struct {
	ConfigId    int64      `json:"configId" db:"config_id"`                 // 配置Id
	ConfigName  string     `json:"configName,omitempty" db:"config_name"`   // 配置名称
	ConfigKey   string     `json:"configKey,omitempty" db:"config_key"`     // 缓存Key
	ConfigValue string     `json:"configValue,omitempty" db:"config_value"` // 缓存Value
	CreateTime  *time.Time `json:"createTime,omitempty" db:"create_time"`   // 创建时间
	UpdateTime  *time.Time `json:"updateTime,omitempty" db:"update_time"`   // 最后更新时间
	Remark      string     `json:"remark,omitempty" db:"remark"`            // 备注
}
