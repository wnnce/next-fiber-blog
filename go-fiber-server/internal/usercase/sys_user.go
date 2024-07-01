package usercase

import "time"

type SysUser struct {
	UserId        uint64     `json:"userId,omitempty" db:"user_id"`                               // 用户Id
	Username      string     `json:"username,omitempty" db:"username" validate:"required,max=64"` // 用户名
	Nickname      string     `json:"nickname,omitempty" db:"nickname"`                            // 昵称
	Password      string     `json:"password,omitempty" db:"password" validate:"required"`        // 密码 sha256
	Email         string     `json:"email,omitempty" db:"email"`                                  // 邮箱
	Phone         string     `json:"phone,omitempty" db:"phone"`                                  // 电话
	Avatar        string     `json:"avatar,omitempty" db:"avatar"`                                // 头像
	LastLoginIp   *string    `json:"lastLoginIp,omitempty" db:"last_login_ip"`                    // 最后登录IP
	LastLoginTime *time.Time `json:"lastLoginTime,omitempty" db:"last_login_time"`                // 最后登录时间
	Roles         []uint     `json:"roles,omitempty" db:"roles" validate:"required"`              // 拥有的角色
	Remark        string     `json:"remark,omitempty" db:"remark"`                                // 备注
	RoleNames     []string   `json:"roleNames,omitempty" db:"-"`
	BaseEntity
}

func (su *SysUser) GetUserId() uint64 {
	return su.UserId
}

func (su *SysUser) GetUsername() string {
	return su.Username
}

func (su *SysUser) GetRoles() []string {
	return su.RoleNames
}

func (su *SysUser) GetPermissions() []string {
	return make([]string, 0)
}
func (su *SysUser) SetUsername(username string) {
	su.Username = username
}
func (su *SysUser) SetRoles(roles []string) {
	su.RoleNames = roles
}
func (su *SysUser) SetPermissions(permissions []string) {
	// TODO 待实现权限
}

// LoginForm 用户登录表单
type LoginForm struct {
	Username string `json:"username" validate:"required"` // 用户名
	Password string `json:"password" validate:"required"` // 密码
	Code     string `json:"code" validate:"required"`     // 验证码
}

// UpdatePasswordForm 更改密码表单
type UpdatePasswordForm struct {
	OldPassword string `json:"oldPassword" validate:"required;len=32"` // 旧密码
	NewPassword string `json:"newPassword" validate:"required;len=32"` // 新密码
	UserId      uint64 `json:"-"`
}

// SysUserQueryForm 系统用户查询表单
type SysUserQueryForm struct {
	Username        string `json:"username"`
	Nickname        string `json:"nickname"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	RoleId          uint   `json:"roleId"`
	CreateTimeBegin string `json:"createTimeBegin"`
	CreateTimeEnd   string `json:"createTimeEnd"`
	PageQueryForm
}

// ISysUserRepo 系统用户Repo接口
type ISysUserRepo interface {
	// Save 保存用户
	Save(user *SysUser) error
	// Update 更新用户
	Update(user *SysUser) error
	// FindUserById 查询用户详情
	FindUserById(userId uint64) (*SysUser, error)
	// Page 分页查询用户
	Page(query *SysUserQueryForm) ([]*SysUser, int64, error)
	// DeleteById 删除用户
	DeleteById(userId int64) error
	// CountByRoleId 通过roleId获取用户数量
	CountByRoleId(roleId int) (int64, error)
	// CountByUsername 通过username获取用户数量
	CountByUsername(username string, userId uint64) (uint8, error)
	// QueryUserByUsernameAndPassword 通过用户名和密码查询用户
	QueryUserByUsernameAndPassword(username, password string) (*SysUser, error)
	// UpdatePassword 更改用户密码
	UpdatePassword(userId uint64, newPassword string) error
	// UpdateLoginRecord 记录登录日志
	UpdateLoginRecord(userId uint64, ip string)
}

// ISysUserService 系统用户Service接口
type ISysUserService interface {
	SaveUser(user *SysUser) error

	UpdateUser(user *SysUser) error

	QueryUserInfo(userId uint64) (*SysUser, error)

	Page(query *SysUserQueryForm) (*PageData[SysUser], error)

	Login(form *LoginForm, ip, ua string) (string, error)

	UpdatePassword(form *UpdatePasswordForm) error

	Delete(userId int64) error
}
