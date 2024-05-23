package service

import (
	"go-fiber-ent-web-layout/internal/tools"
	"go-fiber-ent-web-layout/internal/usercase"
	"log/slog"
)

type SysMenuService struct {
	repo usercase.ISysMenuRepo
}

func NewMenuService(repo usercase.ISysMenuRepo) usercase.ISysMenuService {
	return &SysMenuService{
		repo: repo,
	}
}

func (ms *SysMenuService) CreateMenu(menu *usercase.SysMenu) error {
	err := ms.repo.Save(menu)
	if err != nil {
		slog.Error("菜单添加失败", "error-message", err)
		return tools.FiberServerError("添加失败")
	}
	return nil
}

func (ms *SysMenuService) UpdateMenu(menu *usercase.SysMenu) error {
	err := ms.repo.Update(menu)
	if err != nil {
		slog.Error("菜单更新失败", "error-message", err)
		return tools.FiberServerError("更新失败")
	}
	return nil
}

func (ms *SysMenuService) TreeMenu() ([]*usercase.SysMenu, error) {
	menus, err := ms.repo.ListAll()
	if err != nil {
		slog.Error("获取全部菜单失败", "error-message", err)
		return nil, tools.FiberServerError("获取菜单失败")
	}
	return tools.BuilderTree[uint](menus), nil
}

func (ms *SysMenuService) ManageTreeMenu() ([]*usercase.SysMenu, error) {
	menus, err := ms.repo.ManageListAll()
	if err != nil {
		slog.Error("管理端获取全部菜单失败", "error-message", err)
		return nil, tools.FiberServerError("获取菜单失败")
	}
	return tools.BuilderTree[uint](menus), nil
}

func (ms *SysMenuService) Delete(menuId int) error {
	err := ms.repo.DeleteById(menuId)
	if err != nil {
		slog.Error("删除菜单失败", "error-message", err)
		return tools.FiberServerError("删除失败")
	}
	return nil
}
