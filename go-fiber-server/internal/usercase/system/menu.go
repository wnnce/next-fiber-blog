package usercase

import (
	"go-fiber-ent-web-layout/internal/usercase"
	"time"
)

// Menu 菜单 实现Tree接口
type Menu struct {
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
	DeleteAt   string     `json:"deleteAt,omitempty" db:"delete_at"`           // 是否删除
	CreateTime *time.Time `json:"createTime,omitempty" db:"create_time"`       // 创建时间
	UpdateTime *time.Time `json:"updateTime,omitempty" db:"update_time"`       // 更新时间
	Sort       int        `json:"sort" db:"sort"`                              // 排序
	Children   []*Menu    `json:"children"`                                    // 子节点
}

func (m *Menu) GetId() uint {
	return m.MenuId
}

func (m *Menu) GetParentId() uint {
	return m.ParentId
}

func (m *Menu) AppendChild(t usercase.Tree[uint]) {
	if menu, ok := t.(*Menu); ok {
		if m.Children == nil {
			m.Children = make([]*Menu, 0)
		}
		m.Children = append(m.Children, menu)
	}
}

// IMenuRepo 菜单持久层接口
type IMenuRepo interface {
	// Save 保存
	Save(menu *Menu) error
	// Update 更新
	Update(menu *Menu) error
	// ListAll 查询所有
	ListAll() ([]*Menu, error)
	// ManageListAll 管理端查询所有
	ManageListAll() ([]*Menu, error)
	// DeleteById 通过id删除
	DeleteById(menuId int) error
}

type IMenuService interface {
	CreateMenu(menu *Menu) error

	UpdateMenu(menu *Menu) error

	TreeMenu() ([]*Menu, error)

	ManageTreeMenu() ([]*Menu, error)

	Delete(menuId int) error
}
