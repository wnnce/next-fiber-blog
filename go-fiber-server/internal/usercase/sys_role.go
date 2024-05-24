package usercase

import "time"

// SysRole 系统角色
type SysRole struct {
	RoleId   uint   `json:"roleId,omitempty" db:"role_id"`
	RoleName string `json:"roleName,omitempty" db:"role_name" validate:"required,max=64"`
	RoleKey  string `json:"roleKey,omitempty" db:"role_key" validate:"required,max=64"`
	BaseEntity
	Remark string `json:"remark,omitempty" db:"remark" validate:"max=255,omitempty"`
	Menus  []uint `json:"menus,omitempty" db:"menus"`
}

// SysRoleQueryForm 系统角色查询表单
type SysRoleQueryForm struct {
	Name            string     `json:"name"`
	Key             string     `json:"key"`
	CreateTimeBegin *time.Time `json:"createTimeBegin"`
	CreateTimeEnd   *time.Time `json:"createTimeEnd"`
	PageQueryForm
}

// ISysRoleRepo 系统角色Repo接口
type ISysRoleRepo interface {
	Save(role *SysRole) error

	Update(role *SysRole) error

	ListAll() ([]SysRole, error)

	Page(query *SysRoleQueryForm) ([]*SysRole, int64, error)

	DeleteById(roleId int) error

	CountByRoleKey(roleKey string, roleId uint) (uint8, error)

	ListRoleKeyByIds(ids []uint) ([]string, error)
}

// ISysRoleService 系统角色Service层接口
type ISysRoleService interface {
	SaveRole(role *SysRole) error

	UpdateRole(role *SysRole) error

	List() ([]SysRole, error)

	Page(query *SysRoleQueryForm) (*PageData[SysRole], error)

	Delete(roleId int) error
}
