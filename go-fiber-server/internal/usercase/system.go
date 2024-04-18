package usercase

import "time"

// LoginLog 登录日志
type LoginLog struct {
	ID         int64      `json:"id" db:"id"`                  // 主键ID
	UserId     int64      `json:"userId" db:"user_id"`         // 用户Id
	UserType   int        `json:"userType" db:"user_type"`     // 用户类型
	Username   string     `json:"username" db:"username"`      // 用户名
	LoginIP    string     `json:"loginIp" db:"login_ip"`       // 登录IP
	LoginUa    string     `json:"loginUa" db:"login_ua"`       // 登录UA
	CreateTime *time.Time `json:"createTime" db:"create_time"` // 登录时间
	Remark     string     `json:"remark" db:"remark"`          // 备注
	Result     int        `json:"result" db:"result"`          // 登录结果 0：成功 1：失败
}

// AccessLog 访问日志
type AccessLog struct {
	ID         int64      `json:"id" db:"id"`                  // 主键Id
	Location   string     `json:"location" db:"location"`      // 访问地址
	Referee    string     `json:"referee" db:"referee"`        // 访问来源
	AccessIp   string     `json:"accessIp" db:"access_ip"`     // 访问IP
	AccessUa   string     `json:"accessUa" db:"access_ua"`     // 访问UA
	CreateTime *time.Time `json:"createTime" db:"create_time"` // 访问时间
}

// Menu 菜单
type Menu struct {
	MenuId     int        `json:"menuId" db:"menu_id"`                   // 菜单Id
	MenuName   string     `json:"menuName" db:"menu_name"`               // 菜单名称
	MenuType   int        `json:"menuType" db:"menu_type"`               // 菜单类型 1：目录 2：菜单
	ParentId   int        `json:"parentId" db:"patent_id"`               // 菜单的上级菜单Id
	Path       string     `json:"path" db:"path"`                        // 菜单的路由地址
	Component  string     `json:"component" db:"component"`              // 菜单的组件地址
	Icon       string     `json:"icon" db:"icon"`                        // 菜单图标
	IsFrame    bool       `json:"isFrame" db:"is_frame"`                 // 菜单是否为Iframe窗口
	FrameUrl   bool       `json:"frameUrl" db:"frame_url"`               // Iframe窗口的地址
	IsCache    bool       `json:"isCache" db:"is_cache"`                 // 是否缓存
	IsVisible  bool       `json:"isVisible" db:"is_visible"`             // 是否可见
	IsDisable  bool       `json:"isDisable" db:"is_disable"`             // 是否关闭
	DeleteAt   string     `json:"deleteAt" db:"delete_at"`               // 是否删除
	CreateTime *time.Time `json:"createTime,omitempty" db:"create_time"` // 创建时间
	UpdateTime *time.Time `json:"updateTime,omitempty" db:"update_time"` // 更新时间
	Sort       int        `json:"sort" db:"sort"`                        // 排序
}

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
