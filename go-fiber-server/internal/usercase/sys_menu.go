package usercase

import (
	"time"
)

// SysMenu 菜单 实现Tree接口
type SysMenu struct {
	MenuId     uint       `json:"menuId" db:"menu_id"`                         // 菜单Id
	MenuName   string     `json:"menuName" db:"menu_name" validate:"required"` // 菜单名称
	MenuType   int        `json:"menuType" db:"menu_type" validate:"required"` // 菜单类型 1：目录 2：菜单
	ParentId   uint       `json:"parentId" db:"parent_id"`                     // 菜单的上级菜单Id
	Path       string     `json:"path" db:"path" validate:"required"`          // 菜单的路由地址
	Component  string     `json:"component,omitempty" db:"component"`          // 菜单的组件地址
	Icon       string     `json:"icon" db:"icon" validate:"required"`          // 菜单图标
	IsFrame    bool       `json:"isFrame" db:"is_frame"`                       // 菜单是否为Iframe窗口
	FrameUrl   string     `json:"frameUrl,omitempty" db:"frame_url"`           // Iframe窗口的地址
	IsCache    bool       `json:"isCache" db:"is_cache"`                       // 是否缓存
	IsVisible  bool       `json:"isVisible" db:"is_visible"`                   // 是否可见
	IsDisable  bool       `json:"isDisable" db:"is_disable"`                   // 是否关闭
	DeleteAt   int64      `json:"deleteAt,omitempty" db:"delete_at"`           // 是否删除
	CreateTime *time.Time `json:"createTime,omitempty" db:"create_time"`       // 创建时间
	UpdateTime *time.Time `json:"updateTime,omitempty" db:"update_time"`       // 更新时间
	Sort       int        `json:"sort" db:"sort"`                              // 排序
	Children   []*SysMenu `json:"children"`                                    // 子节点
}

func (m *SysMenu) GetId() uint {
	return m.MenuId
}

func (m *SysMenu) GetParentId() uint {
	return m.ParentId
}

func (m *SysMenu) AppendChild(t Tree[uint]) {
	if menu, ok := t.(*SysMenu); ok {
		if m.Children == nil {
			m.Children = make([]*SysMenu, 0)
		}
		m.Children = append(m.Children, menu)
	}
}

// ISysMenuRepo 菜单持久层接口
type ISysMenuRepo interface {
	// Save 保存
	Save(menu *SysMenu) error
	// Update 更新
	Update(menu *SysMenu) error
	// ListAll 查询所有
	ListAll() ([]*SysMenu, error)
	// RecursiveByMenuIds 递归查询指定的菜单id列表 会查询出父节点
	RecursiveByMenuIds(menuIds []uint) ([]*SysMenu, error)
	// ManageListAll 管理端查询所有
	ManageListAll() ([]*SysMenu, error)
	// DeleteById 通过id删除
	DeleteById(menuId int) error
}

type ISysMenuService interface {
	CreateMenu(menu *SysMenu) error

	UpdateMenu(menu *SysMenu) error

	TreeMenu() ([]*SysMenu, error)

	ManageTreeMenu() ([]*SysMenu, error)

	Delete(menuId int) error
}
