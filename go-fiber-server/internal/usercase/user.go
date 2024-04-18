package usercase

import (
	"github.com/lib/pq"
	"time"
)

// User 用户
type User struct {
	UserId     int64          `json:"userId" db:"user_id"`                   // 用户Id
	NickName   string         `json:"nickName,omitempty" db:"nick_name"`     // 用户昵称
	Summary    string         `json:"summary,omitempty" db:"summary"`        // 用户简介
	Avatar     string         `json:"avatar,omitempty" db:"avatar"`          // 用户头像
	Email      string         `json:"email,omitempty" db:"email"`            // 用户邮箱
	Link       string         `json:"link,omitempty" db:"link"`              // 用户的站点链接
	Username   string         `json:"username,omitempty" db:"username"`      // 用户名
	Password   string         `json:"-" db:"password"`                       // 密码 sha256加密
	Labels     pq.StringArray `json:"labels,omitempty" db:"labels"`          // 用户的标签列表
	UserType   int            `json:"userType" db:"user_type"`               // 用户类型 1：管理员（唯一） 2：普通用户
	CreateTime time.Time      `json:"createTime,omitempty" db:"create_time"` // 用户的创建时间
	UpdateTime time.Time      `json:"updateTime,omitempty" db:"update_time"` // 最后更新时间
	Status     int            `json:"status" db:"status"`                    // 状态 0：正常 1：禁用
}

// UserExtend 用户扩展数据
type UserExtend struct {
	ID               int64     `json:"id" db:"id"`
	UserId           int64     `json:"user_id,omitempty" db:"user_id"`
	Level            int       `json:"level,omitempty" db:"level"`
	Expertise        int64     `json:"expertise,omitempty" db:"expertise"`
	RegisterIp       string    `json:"registerIp,omitempty" db:"register_ip"`
	RegisterLocation string    `json:"registerLocation,omitempty" db:"register_location"`
	CreateTime       time.Time `json:"createTime,omitempty" db:"create_time"`
	UpdateTime       time.Time `json:"updateTime,omitempty" db:"update_time"`
}

// ExpertiseDetail 经验值明细
type ExpertiseDetail struct {
	ID         int64     `json:"id" db:"id"`
	UserId     int64     `json:"userId" db:"user_id"`
	Detail     int64     `json:"detail" db:"detail"`
	DetailType int       `json:"detailType" db:"detail_type"`
	Source     int       `json:"source" db:"source"`
	CreateTime time.Time `json:"createTime" db:"create_time"`
	Remark     string    `json:"remark,omitempty" db:"remark"`
}

func (u *User) GetUserId() int64 {
	return u.UserId
}
func (u *User) GetUserName() string {
	return u.Username
}
func (u *User) GetRoles() []string {
	return make([]string, 0)
}
func (u *User) GetPermissions() []string {
	return make([]string, 0)
}

type IUserService interface {
	Login(user *User) (string, error)
}
