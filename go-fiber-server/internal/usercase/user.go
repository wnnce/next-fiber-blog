package usercase

import (
	"context"
	"github.com/jackc/pgx/v5"
	"time"
)

// User 博客端用户 需要实现 LoginUser 接口
// 博客端和管理端是两套登录逻辑
type User struct {
	UserId     uint64     `json:"userId" db:"user_id"`                   // 用户Id
	Nickname   string     `json:"nickname,omitempty" db:"nick_name"`     // 用户昵称
	Summary    string     `json:"summary,omitempty" db:"summary"`        // 用户简介
	Avatar     string     `json:"avatar,omitempty" db:"avatar"`          // 用户头像
	Email      string     `json:"email,omitempty" db:"email"`            // 用户邮箱
	Link       string     `json:"link,omitempty" db:"link"`              // 用户的站点链接
	Username   string     `json:"username,omitempty" db:"username"`      // 用户名
	Password   string     `json:"-" db:"password"`                       // 密码 sha256加密
	Labels     []string   `json:"labels,omitempty" db:"labels"`          // 用户的标签列表
	UserType   uint8      `json:"userType,omitempty" db:"user_type"`     // 用户类型 1：管理员 2：普通用户
	CreateTime *time.Time `json:"createTime,omitempty" db:"create_time"` // 用户的创建时间
	UpdateTime *time.Time `json:"updateTime,omitempty" db:"update_time"` // 最后更新时间
	Status     uint8      `json:"status,omitempty" db:"status"`          // 状态 0：正常 1：禁用
}

func (self *User) GetUserId() uint64 {
	return self.UserId
}

func (self *User) GetUsername() string {
	return self.Username
}

func (self *User) SetLabels(labels []string) error {
	return nil
}

type UserVo struct {
	User
	Level            uint8  `json:"level,omitempty" db:"level"`                        // 用户等级
	Expertise        uint64 `json:"expertise,omitempty" db:"expertise"`                // 用户经验值
	RegisterIp       string `json:"registerIp,omitempty" db:"register_ip"`             // 用户注册IP
	RegisterLocation string `json:"registerLocation,omitempty" db:"register_location"` // 用户注册地址
}

// UserExtend 用户扩展数据
type UserExtend struct {
	ID               uint64     `json:"id" db:"id"`                                        // 扩展数据ID
	UserId           uint64     `json:"user_id,omitempty" db:"user_id"`                    // 用户Id
	Level            uint8      `json:"level,omitempty" db:"level"`                        // 用户等级
	Expertise        uint64     `json:"expertise,omitempty" db:"expertise"`                // 用户经验值
	RegisterIp       string     `json:"registerIp,omitempty" db:"register_ip"`             // 用户注册IP
	RegisterLocation string     `json:"registerLocation,omitempty" db:"register_location"` // 用户注册地址
	CreateTime       *time.Time `json:"createTime,omitempty" db:"create_time"`             // 创建时间
	UpdateTime       *time.Time `json:"updateTime,omitempty" db:"update_time"`             // 最后更新时间
}

// UserQueryForm 用户查询表单
type UserQueryForm struct {
	Nickname        string     `json:"nickname"`
	Email           string     `json:"email"`
	Username        string     `json:"username"`
	Level           uint8      `json:"level"`
	CreateTimeBegin *time.Time `json:"createTimeBegin"`
	CreateTimeEnd   *time.Time `json:"createTimeEnd"`
	PageQueryForm
}

// ExpertiseDetail 经验值明细
type ExpertiseDetail struct {
	ID         uint64     `json:"id" db:"id"`                   // 主键ID
	UserId     uint64     `json:"userId" db:"user_id"`          // 用户ID
	Detail     int64      `json:"detail" db:"detail"`           // 经验值明细
	DetailType uint8      `json:"detailType" db:"detail_type"`  // 明细类型 1：收入 2：支出
	Source     uint8      `json:"source" db:"source"`           // 来源类型 1：点赞 2：评论
	CreateTime *time.Time `json:"createTime" db:"create_time"`  // 创建时间
	Remark     *string    `json:"remark,omitempty" db:"remark"` // 备注
}

// ExpertiseDetailVo 经验值明细 Vo类
type ExpertiseDetailVo struct {
	ExpertiseDetail
	Username string `json:"username" db:"username"`
	Nickname string `json:"nickname" db:"nickname"`
}

// ExpertiseQueryForm 经验值明细查询参数
type ExpertiseQueryForm struct {
	Username        string     `json:"username"`
	DetailType      uint8      `json:"detailType"`
	Source          uint8      `json:"source"`
	CreateTimeBegin *time.Time `json:"createTimeBegin"`
	CreateTimeEnd   *time.Time `json:"createTimeEnd"`
	PageQueryForm
}

type IUserRepo interface {
	// Transaction 事务
	Transaction(ctx context.Context, fn func(tx pgx.Tx) error) error
	// QueryUserByUsername 通过用户名查询完整用户信息
	QueryUserByUsername(username string) (*User, error)
	// Save 保存用户 会同时初始化userExtend表
	Save(user *UserVo) error
	// QueryUserExtendById 通过userId查询用户扩展信息
	QueryUserExtendById(userId uint64) (*UserExtend, error)
	// SaveExpertiseDetail 保存用户经验值明细
	SaveExpertiseDetail(detail *ExpertiseDetail, tx pgx.Tx) error
	// UpdateUserExpertise 更新用户经验
	UpdateUserExpertise(count int64, userId uint64, tx pgx.Tx) (uint64, uint8, error)
	// UpdateUserLevel 更新用户等级
	UpdateUserLevel(level uint8, userId uint64, tx pgx.Tx) error

	// Page 分页查询用户信息
	Page(query *UserQueryForm) (*PageData[UserVo], error)

	// Update 更新用户信息
	Update(user *User) error

	// PageExpertise 分页查询用户经验值明细
	PageExpertise(query *ExpertiseQueryForm) (*PageData[ExpertiseDetailVo], error)
}

type IUserService interface {
	// LoginWithGithub 通过Github登录
	LoginWithGithub(code, ip string) (string, error)

	// UserInfo 获取用户详细信息包含扩展数据
	UserInfo(user *User) (*UserVo, error)

	// Logout 注销登录
	Logout(userId uint64) error

	// UpdateUserExpertise 更新用户经验值
	// 更新经验值的同时还会保存经验值明细 如果总经验值达到了升级阈值 则会提升用户等级
	UpdateUserExpertise(count int64, userId uint64, source uint8) error

	// PageUser 分页查询用户信息
	PageUser(query *UserQueryForm) (*PageData[UserVo], error)

	// UpdateUser 更新用户信息
	UpdateUser(user *User) error

	// PageExpertise 分页查询用户经验值明细
	PageExpertise(query *ExpertiseQueryForm) (*PageData[ExpertiseDetailVo], error)
}
