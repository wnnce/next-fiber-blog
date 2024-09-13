package usercase

import (
	"time"
)

// User 用户
type User struct {
	UserId     uint64     `json:"userId" db:"user_id"`                   // 用户Id
	NickName   string     `json:"nickName,omitempty" db:"nick_name"`     // 用户昵称
	Summary    string     `json:"summary,omitempty" db:"summary"`        // 用户简介
	Avatar     string     `json:"avatar,omitempty" db:"avatar"`          // 用户头像
	Email      string     `json:"email,omitempty" db:"email"`            // 用户邮箱
	Link       string     `json:"link,omitempty" db:"link"`              // 用户的站点链接
	Username   string     `json:"username,omitempty" db:"username"`      // 用户名
	Password   string     `json:"-" db:"password"`                       // 密码 sha256加密
	Labels     []string   `json:"labels,omitempty" db:"labels"`          // 用户的标签列表
	UserType   uint8      `json:"userType" db:"user_type"`               // 用户类型 1：管理员 2：普通用户
	CreateTime *time.Time `json:"createTime,omitempty" db:"create_time"` // 用户的创建时间
	UpdateTime *time.Time `json:"updateTime,omitempty" db:"update_time"` // 最后更新时间
	Status     uint8      `json:"status" db:"status"`                    // 状态 0：正常 1：禁用
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

// ExpertiseDetail 经验值明细
type ExpertiseDetail struct {
	ID         uint64     `json:"id" db:"id"`
	UserId     uint64     `json:"userId" db:"user_id"`
	Detail     int64      `json:"detail" db:"detail"`
	DetailType uint8      `json:"detailType" db:"detail_type"`
	Source     uint8      `json:"source" db:"source"`
	CreateTime *time.Time `json:"createTime" db:"create_time"`
	Remark     string     `json:"remark,omitempty" db:"remark"`
}

type IUserRepo interface {
	// QueryUserByUsername 通过用户名查询完整用户信息
	QueryUserByUsername(username string) (*User, error)
	// Save 保存用户 会同时初始化userExtend表
	Save(user *UserVo) error
	// QueryUserExtendById 通过userId查询用户扩展信息
	QueryUserExtendById(userId uint64) (*UserExtend, error)
}

type IUserService interface {
	// LoginWithGithub 通过Github登录
	LoginWithGithub(code, ip string) (string, error)

	// UserInfo 获取用户详细信息包含扩展数据
	UserInfo(user *User) (*UserVo, error)

	// Logout 注销登录
	Logout(userId uint64) error
}
