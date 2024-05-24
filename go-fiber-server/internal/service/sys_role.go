package service

import (
	"go-fiber-ent-web-layout/internal/tools"
	"go-fiber-ent-web-layout/internal/usercase"
	"log/slog"
	"math"
)

type SysRoleService struct {
	repo     usercase.ISysRoleRepo
	userRepo usercase.ISysUserRepo
}

func NewSysRoleService(repo usercase.ISysRoleRepo, userRepo usercase.ISysUserRepo) usercase.ISysRoleService {
	return &SysRoleService{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (rs *SysRoleService) SaveRole(role *usercase.SysRole) error {
	total, _ := rs.repo.CountByRoleKey(role.RoleKey, 0)
	if total > 0 {
		return tools.FiberRequestError("roleKey已经存在")
	}
	if err := rs.repo.Save(role); err != nil {
		slog.Error("保存系统角色失败", "err", err)
		return tools.FiberServerError("保存失败")
	}
	return nil
}

func (rs *SysRoleService) UpdateRole(role *usercase.SysRole) error {
	total, _ := rs.repo.CountByRoleKey(role.RoleKey, role.RoleId)
	if total > 0 {
		return tools.FiberRequestError("roleKey已经存在")
	}
	if err := rs.repo.Update(role); err != nil {
		slog.Error("更新系统角色失败", "err", err)
		return tools.FiberServerError("更新失败")
	}
	return nil
}

func (rs *SysRoleService) List() ([]usercase.SysRole, error) {
	roles, err := rs.repo.ListAll()
	if err != nil {
		slog.Error("获取全部系统角色失败", "err", err)
		return nil, tools.FiberServerError("查询失败")
	}
	return roles, nil
}

func (rs *SysRoleService) Page(query *usercase.SysRoleQueryForm) (*usercase.PageData[usercase.SysRole], error) {
	roles, total, err := rs.repo.Page(query)
	if err != nil {
		slog.Error("分页查询系统角色失败", "err", err)
		return nil, err
	}
	pages := int(math.Ceil(float64(total) / float64(query.Size)))
	return &usercase.PageData[usercase.SysRole]{
		Current: query.Page,
		Size:    query.Size,
		Pages:   pages,
		Total:   total,
		Records: roles,
	}, nil
}

func (rs *SysRoleService) Delete(roleId int) error {
	total, err := rs.userRepo.CountByRoleId(roleId)
	if err != nil {
		slog.Error("查询角色关联用户数量失败", "err", err)
		return tools.FiberServerError("删除失败")
	}
	if total > 0 {
		return tools.FiberRequestError("当前角色还有用户未删除")
	}
	if err = rs.repo.DeleteById(roleId); err != nil {
		slog.Error("删除系统角色失败", "err", err)
		return tools.FiberServerError("删除失败")
	}
	return nil
}
