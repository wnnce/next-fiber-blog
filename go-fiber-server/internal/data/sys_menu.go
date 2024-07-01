package data

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-fiber-ent-web-layout/internal/usercase"
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
	sql := `insert into t_system_menu (menu_name, menu_type, parent_id, path, component, icon, is_frame, frame_url, is_cache, 
                    is_visible, is_disable, sort) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) returning menu_id`
	row := m.db.QueryRow(context.Background(), sql, menu.MenuName, menu.MenuType, menu.ParentId, menu.Path, menu.Component,
		menu.Icon, menu.IsFrame, menu.FrameUrl, menu.IsCache, menu.IsVisible, menu.IsDisable, menu.Sort)
	var menuId uint
	err := row.Scan(&menuId)
	if err == nil {
		slog.Info(fmt.Sprintf("菜单保存完成，id:%d", menuId))
		menu.MenuId = menuId
	}
	return err
}

func (m *SysMenuRepo) Update(menu *usercase.SysMenu) error {
	sql := `
		update t_system_menu set update_time = now(), menu_name = $1, menu_type = $2, parent_id = $3, path = $4, component = $5, 
		                  icon = $6, is_frame = $7, frame_url = $8, is_cache = $9, is_visible = $10, is_disable = $11, sort = $12
		where menu_id = $13`
	result, err := m.db.Exec(context.Background(), sql, menu.MenuName, menu.MenuType, menu.ParentId, menu.Path, menu.Component, menu.Icon,
		menu.IsFrame, menu.FrameUrl, menu.IsCache, menu.IsVisible, menu.IsDisable, menu.Sort, menu.MenuId)
	if err == nil {
		slog.Info(fmt.Sprintf("菜单更新完成，row:%d,id:%d", result.RowsAffected(), menu.MenuId))
	}
	return err
}

func (m *SysMenuRepo) ListAll() ([]*usercase.SysMenu, error) {
	rows, err := m.db.Query(context.Background(), `select menu_id, menu_name, menu_type, parent_id, path, component, 
       icon, is_frame, frame_url, is_cache, is_visible, is_disable, sort from t_system_menu where delete_at = 0 order by sort`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (*usercase.SysMenu, error) {
		return pgx.RowToAddrOfStructByNameLax[usercase.SysMenu](row)
	})
}

func (self *SysMenuRepo) RecursiveByMenuIds(menuIds []uint) ([]*usercase.SysMenu, error) {
	sql := `WITH RECURSIVE tree_menu AS (
				SELECT menu_id, menu_name, menu_type, parent_id, path, component, icon, is_frame, frame_url, is_cache, is_visible, is_disable, sort
				FROM t_system_menu
				WHERE menu_id in (%s) and delete_at = 0
			  UNION ALL
				SELECT n.menu_id, n.menu_name, n.menu_type, n.parent_id, n.path, n.component, n.icon, n.is_frame, n.frame_url, n.is_cache, n.is_visible, n.is_disable, n.sort
				FROM t_system_menu n
				JOIN tree_menu np ON np.parent_id = n.menu_id
			)
			SELECT distinct menu_id, menu_name, menu_type, parent_id, path, component, icon, is_frame, frame_url, is_cache, is_visible, is_disable, sort
			FROM tree_menu`
	var builder strings.Builder
	builder.WriteByte('(')
	for index, id := range menuIds {
		if index > 1 {
			builder.WriteByte(',')
		}
		builder.WriteString(strconv.FormatUint(uint64(id), 10))
	}
	builder.WriteByte(')')
	rows, err := self.db.Query(context.Background(), fmt.Sprintf(sql, builder.String()))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (*usercase.SysMenu, error) {
		return pgx.RowToAddrOfStructByNameLax[usercase.SysMenu](row)
	})
}

func (m *SysMenuRepo) ManageListAll() ([]*usercase.SysMenu, error) {
	rows, err := m.db.Query(context.Background(), "select * from t_system_menu where delete_at = 0 order by sort")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (*usercase.SysMenu, error) {
		return pgx.RowToAddrOfStructByNameLax[usercase.SysMenu](row)
	})
}

func (m *SysMenuRepo) DeleteById(menuId int) error {
	result, err := m.db.Exec(context.Background(), "update t_system_menu set delete_at = $1 where menu_id = $2", time.Now().UnixMilli(), menuId)
	if err == nil {
		slog.Info(fmt.Sprintf("菜单删除完成，row:%d,id:%d", result.RowsAffected(), menuId))
	}
	return err
}
