package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go-fiber-ent-web-layout/internal/middleware/auth"
	"go-fiber-ent-web-layout/internal/tools"
	"go-fiber-ent-web-layout/internal/tools/github"
	"go-fiber-ent-web-layout/internal/tools/region"
	"go-fiber-ent-web-layout/internal/usercase"
	"go-fiber-ent-web-layout/pkg/pool"
	"log/slog"
)

const levelExpertiseNumber = 100

type UserService struct {
	repo usercase.IUserRepo
}

func NewUserService(repo usercase.IUserRepo) usercase.IUserService {
	return &UserService{
		repo: repo,
	}
}

func (self *UserService) LoginWithGithub(code, ip string) (string, error) {
	accessToken, err := github.AccessToken(code)
	if err != nil || "" == accessToken {
		slog.Error("获取Github AccessToken失败", "err", err.Error())
		return "", tools.FiberServerError("获取Github Token失败")
	}
	ch1 := make(chan *github.Profile)
	ch2 := make(chan []*github.Email)
	pool.Go(func() {
		defer close(ch1)
		profile, profileErr := github.UserProfile(accessToken)
		if profileErr != nil {
			slog.Error("获取Github用户信息失败", "err", err.Error())
			return
		}
		ch1 <- profile
	})
	pool.Go(func() {
		defer close(ch2)
		emails, emailErr := github.UserEmails(accessToken)
		if emailErr != nil {
			slog.Error("获取Github用户邮箱失败", "err", err.Error())
			return
		}
		ch2 <- emails
	})
	profile := <-ch1
	emails := <-ch2
	if profile == nil || emails == nil {
		return "", tools.FiberServerError("获取三方登录用户信息失败")
	}
	user, err := self.repo.QueryUserByUsername(profile.Login)
	if err != nil {
		return "", tools.FiberServerError("登录失败")
	}
	// 如果用户不存在 那么就创建用户
	if user == nil {
		var email string
		if len(emails) == 1 {
			email = emails[0].Email
		} else {
			for _, value := range emails {
				if value.Primary {
					email = value.Email
				}
			}
		}
		location := region.SearchLocation(ip)
		userVo := &usercase.UserVo{
			User: usercase.User{
				Username: profile.Login,
				Nickname: profile.Name,
				Email:    email,
				Avatar:   profile.AvatarUrl,
				Summary:  profile.Company,
				Link:     profile.HtmlUrl,
			},
			Level:            1,
			RegisterIp:       ip,
			RegisterLocation: location,
		}
		if err = self.repo.Save(userVo); err != nil {
			slog.Error("添加用户失败", "err", err.Error())
			return "", tools.FiberServerError("保存用户登录信息失败")
		}
		user, _ = self.repo.QueryUserByUsername(profile.Login)
	}
	if user.Status != 0 {
		return "", tools.FiberServerError("用户被禁用或状态异常")
	}
	subject := uuid.New().String()
	token, err := tools.GenerateToken(subject, true)
	if err != nil {
		slog.Error("生成用户登录Token失败", "err", err.Error())
		return "", tools.FiberServerError("登录失败")
	}
	if err = auth.AddClassicLoginUser(subject, user); err != nil {
		slog.Error("保存登录用户信息到Redis失败", "err", err.Error())
		return "", tools.FiberServerError("登录失败")
	}
	return token, nil
}

func (self *UserService) UserInfo(user *usercase.User) (*usercase.UserVo, error) {
	extend, err := self.repo.QueryUserExtendById(user.GetUserId())
	if err != nil {
		slog.Error("获取登录用户扩展数据失败", "err", err.Error())
		return nil, err
	}
	return &usercase.UserVo{
		User:             *user,
		Level:            extend.Level,
		Expertise:        extend.Expertise,
		RegisterIp:       extend.RegisterIp,
		RegisterLocation: extend.RegisterLocation,
	}, nil
}

func (self *UserService) Logout(userId uint64) error {
	if err := auth.RemoveClassicLoginUserById(userId); err != nil {
		slog.Error("注销博客端登录用户失败", "err", err.Error())
		return tools.FiberServerError("注销登录失败")
	}
	return nil
}

// UpdateUserExpertise 更新用户经验值
// 在更新用户经验值时 判断当前经验值是否达到了这个等级的经验上限
// 如果达到 那么就将用户升级
func (self *UserService) UpdateUserExpertise(count int64, userId uint64) error {
	return self.repo.Transaction(context.Background(), func(tx pgx.Tx) error {
		if err := self.repo.SaveExpertiseDetail(&usercase.ExpertiseDetail{
			UserId:     userId,
			Detail:     count,
			DetailType: 1,
			Source:     2,
		}, tx); err != nil {
			slog.Error("保存经验值变更明细失败", "err", err.Error())
			return err
		}
		expertise, level, err := self.repo.UpdateUserExpertise(count, userId, tx)
		if err != nil {
			slog.Error("更新用户经验值失败", "err", err.Error(), "userId", userId)
			return err
		}
		// 计算用户当前等级所需要的升级经验值
		upgradeExpertise := levelExpertiseNumber << level
		// 如果当前经验值大于当前等级的总经验值 那么就将用户升级
		if expertise >= uint64(upgradeExpertise) {
			if upgradeErr := self.repo.UpdateUserLevel(level+1, userId, tx); upgradeErr != nil {
				slog.Error("更新用户等级失败", "err", err.Error(), "userId", userId, "level", level)
				return upgradeErr
			}
		}
		return nil
	})
}

func (self *UserService) PageUser(query *usercase.UserQueryForm) (*usercase.PageData[usercase.UserVo], error) {
	page, err := self.repo.Page(query)
	if err != nil {
		slog.Error("分页查询博客用户信息失败", "error", err.Error())
		return nil, tools.FiberServerError("查询失败")
	}
	return page, nil
}

func (self *UserService) UpdateUser(user *usercase.User) error {
	err := self.repo.Update(user)
	if err != nil {
		slog.Error("更新用户信息失败", "error", err.Error())
		return tools.FiberServerError("更新失败")
	}
	// 当禁用用户时 判断当前用户是否已经登录
	// 如果登录那么需要强制退出
	pool.Go(func() {
		if user.Status == 1 {
			if logoutErr := auth.RemoveClassicLoginUserById(user.UserId); err != nil {
				slog.Error("删除博客端用户登录信息失败", "error", logoutErr)
			}
		}
	})
	return nil
}

func (self *UserService) PageExpertise(query *usercase.ExpertiseQueryForm) (*usercase.PageData[usercase.ExpertiseDetailVo], error) {
	page, err := self.repo.PageExpertise(query)
	if err != nil {
		slog.Error("查询用户经验值明细失败", "error", err.Error())
		return nil, tools.FiberServerError("查询失败")
	}
	return page, nil
}
