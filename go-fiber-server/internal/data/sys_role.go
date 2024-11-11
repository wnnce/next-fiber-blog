package data

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-fiber-ent-web-layout/internal/tools"
	"go-fiber-ent-web-layout/internal/usercase"
	sqlbuild "go-fiber-ent-web-layout/pkg/sql-build"
	"log/slog"
	"time"
)

type SysRoleRepo struct {
	db *pgxpool.Pool
}

func NewSysRoleRepo(data *Data) usercase.ISysRoleRepo {
	return &SysRoleRepo{
		db: data.Db,
	}
}

func (self *SysRoleRepo) Save(role *usercase.SysRole) error {
	builder := sqlbuild.NewInsertBuilder("t_system_role").
		Fields("role_name", "role_key", "menus", "sort", "status", "remark").
		Values(role.RoleName, role.RoleKey, role.Menus, role.Sort, role.Status, role.Remark).
		Returning("role_id")
	row := self.db.QueryRow(context.Background(), builder.Sql(), builder.Args()...)
	var roleId uint
	err := row.Scan(&roleId)
	if err == nil {
		role.RoleId = roleId
		slog.Info("系统角色插入完成", "roleId", roleId)
	}
	return err
}

func (self *SysRoleRepo) Update(role *usercase.SysRole) error {
	builder := sqlbuild.NewUpdateBuilder("t_system_role").
		SetRaw("update_time", "now()").
		SetByMap(map[string]any{
			"role_name": role.RoleName,
			"role_key":  role.RoleKey,
			"menus":     role.Menus,
			"sort":      role.Sort,
			"status":    role.Status,
			"remark":    role.Remark,
		}).Where("role_id").Eq(role.RoleId).BuildAsUpdate()
	result, err := self.db.Exec(context.Background(), builder.Sql(), builder.Args()...)
	if err == nil {
		slog.Info(fmt.Sprintf("更新系统角色完成，row:%d,roleId:%d", result.RowsAffected(), role.RoleId))
	}
	return err
}

func (self *SysRoleRepo) UpdateSelective(form *usercase.SysRoleUpdateForm) error {
	builder := sqlbuild.NewUpdateBuilder("t_system_role").
		SetRaw("update_time", "now()").
		SetByCondition(form.Status != nil, "status", form.Status).
		Where("role_id").Eq(form.RoleId).BuildAsUpdate()
	result, err := self.db.Exec(context.Background(), builder.Sql(), builder.Args()...)
	if err == nil {
		slog.Info("快捷更新系统角色完成", "row", result.RowsAffected(), "roleId", form.RoleId)
	}
	return err
}

func (self *SysRoleRepo) ListAll() ([]usercase.SysRole, error) {
	builder := sqlbuild.NewSelectBuilder("t_system_role").
		Select("role_id", "role_key", "role_name", "sort").
		Where("status").EqRaw("0").
		And("delete_at").EqRaw("0").BuildAsSelect()
	rows, err := self.db.Query(context.Background(), builder.Sql())
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (usercase.SysRole, error) {
		return pgx.RowToStructByNameLax[usercase.SysRole](row)
	})
}

func (self *SysRoleRepo) Page(query *usercase.SysRoleQueryForm) ([]*usercase.SysRole, int64, error) {
	builder := sqlbuild.NewSelectBuilder("t_system_role").
		WhereByCondition(query.Name != "", "role_name").Like("%"+query.Name+"%").
		AndByCondition(query.Key != "", "role_key").Like("%"+query.Key+"%").
		AndByCondition(query.CreateTimeBegin != "", "create_time").Ge(query.CreateTimeBegin).
		AndByCondition(query.CreateTimeEnd != "", "create_time").Le(query.CreateTimeEnd).
		And("delete_at").EqRaw("0").
		BuildAsSelect().
		OrderBy("sort", "create_time desc")
	var total int64
	row := self.db.QueryRow(context.Background(), builder.CountSql(), builder.Args()...)
	if err := row.Scan(&total); err != nil {
		return nil, 0, err
	}
	roles := make([]*usercase.SysRole, 0)
	if total == 0 {
		return roles, total, nil
	}
	offset := tools.ComputeOffset(total, query.Page, query.Size, false)
	builder.Limit(int64(query.Size)).Offset(offset)
	rows, err := self.db.Query(context.Background(), builder.Sql(), builder.Args()...)
	defer rows.Close()
	if err != nil {
		return nil, 0, err
	}
	roles, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (*usercase.SysRole, error) {
		return pgx.RowToAddrOfStructByName[usercase.SysRole](row)
	})
	return roles, total, err
}

func (self *SysRoleRepo) DeleteById(roleId int) error {
	builder := sqlbuild.NewUpdateBuilder("t_system_role").
		Set("delete_at", time.Now().UnixMilli()).
		Where("role_id").Eq(roleId).BuildAsUpdate()
	result, err := self.db.Exec(context.Background(), builder.Sql(), builder.Args()...)
	if err == nil {
		slog.Info("删除系统角色完成", "row", result.RowsAffected(), "roleId", roleId)
	}
	return err
}

func (self *SysRoleRepo) CountByRoleKey(roleKey string, roleId uint) (uint8, error) {
	builder := sqlbuild.NewSelectBuilder("t_system_role").
		Select("count(*)").
		Where("role_id").Ne(roleId).
		And("role_key").Eq(roleKey).
		And("delete_at").EqRaw("0").BuildAsSelect()
	row := self.db.QueryRow(context.Background(), builder.Sql(), builder.Args()...)
	var total uint8
	err := row.Scan(&total)
	return total, err
}

func (self *SysRoleRepo) ListRoleKeyByIds(ids []uint) ([]string, error) {
	builder := sqlbuild.NewSelectBuilder("t_system_role").
		Select("role_key").
		Where("role_id").In(sqlbuild.SliceToAnySlice[uint](ids)...).
		And("status").EqRaw("0").
		And("delete_at").EqRaw("0").BuildAsSelect()
	rows, err := self.db.Query(context.Background(), builder.Sql(), builder.Args()...)
	defer rows.Close()
	if err != nil {
		slog.Error("查询角色Key列表失败", "err", err)
		return nil, err
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (string, error) {
		var key string
		err = row.Scan(&key)
		return key, err
	})
}
