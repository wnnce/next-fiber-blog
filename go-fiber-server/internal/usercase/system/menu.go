package usercase

import "time"

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
	Children   []*Menu    `json:"children"`                              // 子节点
}
