package service

import (
	"github.com/google/uuid"
	"go-fiber-ent-web-layout/internal/middleware/auth"
	"go-fiber-ent-web-layout/internal/tools"
	"go-fiber-ent-web-layout/internal/tools/github"
	"go-fiber-ent-web-layout/internal/tools/region"
	"go-fiber-ent-web-layout/internal/usercase"
	"go-fiber-ent-web-layout/pkg/pool"
	"log/slog"
)

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
				NickName: profile.Name,
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
	token, err := tools.GenerateToken(subject)
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
