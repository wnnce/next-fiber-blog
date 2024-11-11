package data

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-fiber-ent-web-layout/internal/usercase"
	sqlbuild "go-fiber-ent-web-layout/pkg/sql-build"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

type SysMenuRepo struct {
	db *pgxpool.Pool
}

func NewSysMenuRepo(data *Data) usercase.ISysMenuRepo {
	return &SysMenuRepo{
		db: data.Db,
	}
}

func (m *SysMenuRepo) Save(menu *usercase.SysMenu) error {
	builder := sqlbuild.NewInsertBuilder("t_system_menu").
		Fields("menu_name", "menu_type", "parent_id", "path", "component", "icon", "is_frame", "frame_url",
			"is_cache", "is_visible", "is_disable", "sort").
		Values(menu.MenuName, menu.MenuType, menu.ParentId, menu.Path, menu.Component, menu.Icon, menu.IsFrame,
			menu.FrameUrl, menu.IsCache, menu.IsVisible, menu.IsDisable, menu.Sort).
		Returning("menu_id")
	row := m.db.QueryRow(context.Background(), builder.Sql(), builder.Args()...)
	var menuId uint
	err := row.Scan(&menuId)
	if err == nil {
		slog.Info(fmt.Sprintf("菜单保存完成，id:%d", menuId))
		menu.MenuId = menuId
	}
	return err
}

func (m *SysMenuRepo) Update(menu *usercase.SysMenu) error {
	builder := sqlbuild.NewUpdateBuilder("t_system_menu").
		SetRaw("update_time", "now()").
		SetByMap(map[string]any{
			"menu_name":  menu.MenuName,
			"menu_type":  menu.MenuType,
			"parent_id":  menu.ParentId,
			"path":       menu.Path,
			"component":  menu.Component,
			"icon":       menu.Icon,
			"is_frame":   menu.IsFrame,
			"frame_url":  menu.FrameUrl,
			"is_cache":   menu.IsCache,
			"is_visible": menu.IsVisible,
			"is_disable": menu.IsDisable,
			"sort":       menu.Sort,
		}).Where("menu_id").Eq(menu.MenuId).BuildAsUpdate()
	result, err := m.db.Exec(context.Background(), builder.Sql(), builder.Args()...)
	if err == nil {
		slog.Info(fmt.Sprintf("菜单更新完成，row:%d,id:%d", result.RowsAffected(), menu.MenuId))
	}
	return err
}

func (m *SysMenuRepo) ListAll() ([]*usercase.SysMenu, error) {
	builder := sqlbuild.NewSelectBuilder("t_system_menu").
		Select("menu_id", "menu_name", "menu_type", "parent_id", "path", "component", "icon", "is_frame", "frame_url",
			"is_cache", "is_visible", "is_disable", "sort").
		Where("delete_at").EqRaw("0").BuildAsSelect().
		OrderBy("sort")
	rows, err := m.db.Query(context.Background(), builder.Sql())
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (*usercase.SysMenu, error) {
		return pgx.RowToAddrOfStructByNameLax[usercase.SysMenu](row)
	})
}

func (self *SysMenuRepo) RecursiveByRoleKeys(roleKeys []string) ([]*usercase.SysMenu, error) {
	sql := `WITH RECURSIVE tree_menu AS (
				SELECT menu_id, menu_name, menu_type, parent_id, path, component, icon, is_frame, frame_url, is_cache, is_visible, is_disable, sort
				FROM t_system_menu
				WHERE menu_id in (select distinct unnest(menus) from t_system_role where role_key in (%s)) and delete_at = 0
			  UNION ALL
				SELECT n.menu_id, n.menu_name, n.menu_type, n.parent_id, n.path, n.component, n.icon, n.is_frame, n.frame_url, n.is_cache, n.is_visible, n.is_disable, n.sort
				FROM t_system_menu n
				JOIN tree_menu np ON np.parent_id = n.menu_id
			)
			SELECT distinct menu_id, menu_name, menu_type, parent_id, path, component, icon, is_frame, frame_url, is_cache, is_visible, is_disable, sort
			FROM tree_menu order by sort`
	var builder strings.Builder
	args := make([]any, 0)
	builder.WriteByte('(')
	for index, key := range roleKeys {
		args = append(args, key)
		if index > 1 {
			builder.WriteByte(',')
		}
		builder.WriteString("$" + strconv.Itoa(len(args)))
	}
	builder.WriteByte(')')
	rows, err := self.db.Query(context.Background(), fmt.Sprintf(sql, builder.String()), args...)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (*usercase.SysMenu, error) {
		return pgx.RowToAddrOfStructByNameLax[usercase.SysMenu](row)
	})
}

func (m *SysMenuRepo) ManageListAll() ([]*usercase.SysMenu, error) {
	builder := sqlbuild.NewSelectBuilder("t_system_menu").
		Where("delete_at").EqRaw("0").BuildAsSelect().
		OrderBy("sort")
	rows, err := m.db.Query(context.Background(), builder.Sql())
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (*usercase.SysMenu, error) {
		return pgx.RowToAddrOfStructByNameLax[usercase.SysMenu](row)
	})
}

func (m *SysMenuRepo) DeleteById(menuId int) error {
	builder := sqlbuild.NewUpdateBuilder("t_system_menu").
		Set("delete_at", time.Now().UnixMilli()).
		Where("menu_id").Eq(menuId).BuildAsUpdate()
	result, err := m.db.Exec(context.Background(), builder.Sql(), builder.Args()...)
	if err == nil {
		slog.Info(fmt.Sprintf("菜单删除完成，row:%d,id:%d", result.RowsAffected(), menuId))
	}
	return err
}
