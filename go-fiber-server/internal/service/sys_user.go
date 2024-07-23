package service

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go-fiber-ent-web-layout/internal/middleware/auth"
	"go-fiber-ent-web-layout/internal/tools"
	"go-fiber-ent-web-layout/internal/usercase"
	"go-fiber-ent-web-layout/pkg/pool"
	"hash"
	"log/slog"
	"math"
	"sync"
)

type SysUserService struct {
	repo         usercase.ISysUserRepo
	roleRepo     usercase.ISysRoleRepo
	otherService usercase.IOtherService
	hashPool     *sync.Pool
}

func NewSysUserService(repo usercase.ISysUserRepo, roleRepo usercase.ISysRoleRepo,
	otherService usercase.IOtherService) usercase.ISysUserService {
	return &SysUserService{
		repo:         repo,
		roleRepo:     roleRepo,
		otherService: otherService,
		hashPool: &sync.Pool{
			New: func() any {
				return sha256.New()
			},
		},
	}
}

func (sus *SysUserService) SaveUser(user *usercase.SysUser) error {
	cryptoPassword, err := sus.cryptoPassword(user.Password)
	if err != nil {
		slog.Error("密码解析失败", "err", err)
		return tools.FiberServerError("保存失败")
	}
	user.Password = cryptoPassword
	if err = sus.repo.Save(user); err != nil {
		slog.Error("保存用户失败", "err", err)
		return tools.FiberServerError("保存失败")
	}
	return nil
}

func (sus *SysUserService) UpdateUser(user *usercase.SysUser) error {
	if err := sus.repo.Update(user); err != nil {
		slog.Error("更新用户失败", "err", err)
		return tools.FiberServerError("更新失败")
	}
	// 更新用户信息后 需要同步同步更新管理端登录用户信息
	if logUser := auth.GetManageLoginUserById(user.UserId); logUser != nil {
		if *user.Status == 1 {
			auth.RemoveManageLoginUserById(user.UserId)
		} else {
			if logUser.GetUsername() != user.Username {
				logUser.SetUsername(user.Username)
			}
			if keys, err := sus.roleRepo.ListRoleKeyByIds(user.Roles); err == nil {
				logUser.SetRoles(keys)
			} else {
				slog.Error("用户更新后获取角色列表失败，清空登录用户角色", "userId", user.UserId, "roles", user.Roles)
				logUser.SetRoles(make([]string, 0))
			}
		}
	}
	return nil
}

func (sus *SysUserService) QueryUserInfo(userId uint64) (*usercase.SysUser, error) {
	user, err := sus.repo.FindUserById(userId)
	if err != nil {
		slog.Error("查询用户失败", "err", err)
		return nil, tools.FiberServerError("查询失败")
	}
	if user == nil {
		slog.Warn("查询用户不存在", "userId", userId)
		return nil, fiber.NewError(fiber.StatusNotFound, "用户不存在或被禁用")
	}
	return user, nil
}

func (sus *SysUserService) Page(query *usercase.SysUserQueryForm) (*usercase.PageData[usercase.SysUser], error) {
	users, total, err := sus.repo.Page(query)
	if err != nil {
		slog.Error("分页查询系统用户失败", "err", err)
		return nil, tools.FiberServerError("查询失败")
	}
	pages := int(math.Ceil(float64(total) / float64(query.Size)))
	return &usercase.PageData[usercase.SysUser]{
		Current: query.Page,
		Size:    query.Size,
		Pages:   pages,
		Total:   total,
		Records: users,
	}, nil
}

func (sus *SysUserService) Login(form *usercase.LoginForm, ip, ua string) (string, error) {
	password, _ := sus.cryptoPassword(form.Password)
	user, err := sus.repo.QueryUserByUsernameAndPassword(form.Username, password)
	if err != nil {
		slog.Error("查询登录用户失败", "form", form, "err", err)
		return "", tools.FiberServerError("登录失败")
	}
	if user == nil {
		errMessage := "用户不存在或密码错误"
		pool.Go(func() {
			sus.recordLogin(0, form.Username, ip, ua, errMessage, false)
		})
		return "", tools.FiberRequestError(errMessage)
	}
	if *user.Status == 1 {
		errMessage := "用户被禁用"
		pool.Go(func() {
			sus.recordLogin(user.UserId, user.Username, ip, ua, errMessage, false)
		})
		slog.Warn("登录用户被禁用", "userId", user.UserId)
		return "", tools.FiberServerError(errMessage)
	}
	keys, err := sus.roleRepo.ListRoleKeyByIds(user.Roles)
	if err != nil {
		slog.Error("查询用户角色列表失败", "err", err)
		return "", tools.FiberServerError("登录失败")
	}
	user.RoleNames = keys
	subject := uuid.New().String()
	token, err := tools.GenerateToken(subject)
	if err != nil {
		slog.Error("登录Token生成失败", "err", err)
		return "", tools.FiberServerError("登录失败")
	}
	// 存入的不是完整的token 而是token中的uuid
	auth.AddManageLoginUser(subject, user, auth.ManageUserCacheExpireTime)
	pool.Go(func() {
		sus.recordLogin(user.UserId, user.Username, ip, ua, "登录成功", true)
	})
	return token, nil
}

func (sus *SysUserService) UpdatePassword(form *usercase.UpdatePasswordForm) error {
	oldPassword, _ := sus.cryptoPassword(form.OldPassword)
	newPassword, _ := sus.cryptoPassword(form.NewPassword)
	user, err := sus.repo.FindUserById(form.UserId)
	if err != nil {
		slog.Error("查询用户失败", "err", err)
		return tools.FiberServerError("更新密码失败")
	}
	if user.Password != oldPassword {
		slog.Error("用户旧密码错误", "oldPasswd", oldPassword, "userId", form.UserId)
		return tools.FiberServerError("旧密码错误")
	}
	if err = sus.repo.UpdatePassword(user.UserId, newPassword); err != nil {
		slog.Error("更新用户密码失败", "err", err)
		return tools.FiberServerError("密码更新失败")
	}
	return nil
}

func (sus *SysUserService) Delete(userId int64) error {
	if err := sus.repo.DeleteById(userId); err != nil {
		slog.Error("系统用户删除失败", "err", err)
		return tools.FiberServerError("删除失败")
	}
	// 同时删除管理端登录用户缓存
	auth.RemoveManageLoginUserById(uint64(userId))
	return nil
}

func (self *SysUserService) Logout(userId uint64) {
	auth.RemoveManageLoginUserById(userId)
}

// 密码摘要 将前面传入的密码通过base64解码后再通过sha256摘要计算
func (sus *SysUserService) cryptoPassword(password string) (string, error) {
	originPassword, err := base64.URLEncoding.DecodeString(password)
	if err != nil {
		slog.Error("登录密码解析失败", "password", password, "err", err)
		return "", tools.FiberServerError("密码解析失败")
	}
	hasher := sus.hashPool.Get().(hash.Hash)
	defer func() {
		hasher.Reset()
		sus.hashPool.Put(hasher)
	}()
	hasher.Write(originPassword)
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// 记录登录日志
func (sus *SysUserService) recordLogin(userId uint64, username, ip, ua, remark string, success bool) {
	var result int
	if success {
		result = 0
	} else {
		result = 1
	}
	sus.otherService.TraceLogin(&usercase.LoginLog{
		UserId:    userId,
		Username:  username,
		UserType:  3,
		LoginIP:   ip,
		LoginUa:   ua,
		Remark:    remark,
		LoginType: 2,
		Result:    result,
	})
	if success {
		sus.repo.UpdateLoginRecord(userId, ip)
	}
}
