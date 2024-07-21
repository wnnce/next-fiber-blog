package data

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-fiber-ent-web-layout/internal/tools"
	"go-fiber-ent-web-layout/internal/usercase"
	"log/slog"
	"strconv"
	"strings"
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

func (sr *SysRoleRepo) Save(role *usercase.SysRole) error {
	sql := "insert into t_system_role (role_name, role_key, menus, sort, status, remark) values ($1, $2, $3, $4, $5, $6) returning role_id"
	row := sr.db.QueryRow(context.Background(), sql, role.RoleName, role.RoleKey, role.Menus, role.Sort, role.Status, role.Remark)
	var roleId uint
	err := row.Scan(&roleId)
	if err == nil {
		role.RoleId = roleId
		slog.Info("系统角色插入完成", "roleId", roleId)
	}
	return err
}

func (sr *SysRoleRepo) Update(role *usercase.SysRole) error {
	sql := `update t_system_role set update_time = now(), role_name = $1, role_key = $2, menus = $3, sort = $4, status = $5, remark = $6 
                     where role_id = $7`
	result, err := sr.db.Exec(context.Background(), sql, role.RoleName, role.RoleKey, role.Menus, role.Sort, role.Status, role.Remark, role.RoleId)
	if err == nil {
		slog.Info(fmt.Sprintf("更新系统角色完成，row:%d,roleId:%d", result.RowsAffected(), role.RoleId))
	}
	return err
}

func (sr *SysRoleRepo) ListAll() ([]usercase.SysRole, error) {
	rows, err := sr.db.Query(context.Background(), "select role_id, role_key, role_name, sort from t_system_role where delete_at = 0 and status = 0")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (usercase.SysRole, error) {
		return pgx.RowToStructByNameLax[usercase.SysRole](row)
	})
}

func (sr *SysRoleRepo) Page(query *usercase.SysRoleQueryForm) ([]*usercase.SysRole, int64, error) {
	var condition strings.Builder
	condition.WriteString("where delete_at = 0 ")
	args := make([]any, 0)
	if query.Name != "" {
		args = append(args, "%"+query.Name+"%")
		condition.WriteString(fmt.Sprintf("and role_name like $%d ", len(args)))
	}
	if query.Key != "" {
		args = append(args, "%"+query.Key+"%")
		condition.WriteString(fmt.Sprintf("and role_key like $%d ", len(args)))
	}
	timeQueryConditionBuilder(query.CreateTimeBegin, query.CreateTimeEnd, &condition, &args)
	total, err := sr.conditionToal(condition.String(), args...)
	if err != nil {
		return nil, 0, err
	}
	roles := make([]*usercase.SysRole, 0)
	if total == 0 {
		return roles, total, nil
	}
	offset := tools.ComputeOffset(total, query.Page, query.Size, false)
	condition.WriteString(fmt.Sprintf(" order by sort, create_time desc limit $%d offset $%d", len(args)+1, len(args)+2))
	args = append(args, query.Size, offset)
	rows, err := sr.db.Query(context.Background(), "select * from t_system_role "+condition.String(), args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	roles, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (*usercase.SysRole, error) {
		return pgx.RowToAddrOfStructByName[usercase.SysRole](row)
	})
	return roles, total, err
}

func (sr *SysRoleRepo) DeleteById(roleId int) error {
	result, err := sr.db.Exec(context.Background(), "update t_system_role set delete_at = $1 where role_id = $2", time.Now().UnixMilli(), roleId)
	if err == nil {
		slog.Info("删除系统角色完成", "row", result.RowsAffected(), "roleId", roleId)
	}
	return err
}

func (sr *SysRoleRepo) CountByRoleKey(roleKey string, roleId uint) (uint8, error) {
	row := sr.db.QueryRow(context.Background(), "select count(*) from t_system_role where role_key = $1 and delete_at = 0 and role_id != $2",
		roleKey, roleId)
	var total uint8
	err := row.Scan(&total)
	return total, err
}

func (sr *SysRoleRepo) ListRoleKeyByIds(ids []uint) ([]string, error) {
	var builder strings.Builder
	builder.WriteString("select role_key from t_system_role where role_id in (")
	for i, id := range ids {
		if i > 0 {
			builder.WriteByte(',')
		}
		builder.WriteString(strconv.FormatUint(uint64(id), 10))
	}
	builder.WriteString(") and delete_at = 0 and status = 0")
	fmt.Println(builder.String())
	rows, err := sr.db.Query(context.Background(), builder.String())
	if err != nil {
		slog.Error("查询角色Key列表失败", "err", err)
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (string, error) {
		var key string
		err = row.Scan(&key)
		return key, err
	})
}

func (sr *SysRoleRepo) conditionToal(condition string, args ...any) (int64, error) {
	row := sr.db.QueryRow(context.Background(), "select count(*) from t_system_role "+condition, args...)
	var total int64
	err := row.Scan(&total)
	return total, err
}
