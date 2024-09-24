package service

import (
	"go-fiber-ent-web-layout/internal/tools"
	"go-fiber-ent-web-layout/internal/usercase"
	"log/slog"
	"slices"
)

type SysMenuService struct {
	repo usercase.ISysMenuRepo
}

func NewMenuService(repo usercase.ISysMenuRepo) usercase.ISysMenuService {
	return &SysMenuService{
		repo: repo,
	}
}

func (self *SysMenuService) CreateMenu(menu *usercase.SysMenu) error {
	err := self.repo.Save(menu)
	if err != nil {
		slog.Error("菜单添加失败", "error-message", err)
		return tools.FiberServerError("添加失败")
	}
	return nil
}

func (self *SysMenuService) UpdateMenu(menu *usercase.SysMenu) error {
	err := self.repo.Update(menu)
	if err != nil {
		slog.Error("菜单更新失败", "error-message", err)
		return tools.FiberServerError("更新失败")
	}
	return nil
}

func (self *SysMenuService) TreeMenu(roleKeys []string) (menus []*usercase.SysMenu, err error) {
	if slices.Contains(roleKeys, "admin") {
		menus, err = self.repo.ListAll()
	} else {
		menus, err = self.repo.RecursiveByRoleKeys(roleKeys)
	}
	if err != nil {
		slog.Error("获取全部菜单失败", "error-message", err)
		return nil, tools.FiberServerError("获取菜单失败")
	}
	return tools.BuilderTree[uint](menus), nil
}

func (self *SysMenuService) ManageTreeMenu() ([]*usercase.SysMenu, error) {
	menus, err := self.repo.ManageListAll()
	if err != nil {
		slog.Error("管理端获取全部菜单失败", "error-message", err)
		return nil, tools.FiberServerError("获取菜单失败")
	}
	return tools.BuilderTree[uint](menus), nil
}

func (self *SysMenuService) Delete(menuId int) error {
	err := self.repo.DeleteById(menuId)
	if err != nil {
		slog.Error("删除菜单失败", "error-message", err)
		return tools.FiberServerError("删除失败")
	}
	return nil
}
