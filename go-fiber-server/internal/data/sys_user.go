package data

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-fiber-ent-web-layout/internal/tools"
	"go-fiber-ent-web-layout/internal/usercase"
	"log/slog"
	"strings"
	"time"
)

type SysUserRepo struct {
	db *pgxpool.Pool
}

func NewSysUserRepo(data *Data) usercase.ISysUserRepo {
	return &SysUserRepo{
		db: data.Db,
	}
}

func (sur *SysUserRepo) Save(user *usercase.SysUser) error {
	sql := `insert into t_system_user (username, nickname, password, email, phone, avatar, roles, sort, status, remark) 
			values ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10 ) returning user_id`
	row := sur.db.QueryRow(context.Background(), sql, user.Username, user.Nickname, user.Password, user.Email,
		user.Phone, user.Avatar, user.Roles, user.Sort, user.Status, user.Remark)
	var userId uint64
	err := row.Scan(&userId)
	if err == nil {
		user.UserId = userId
		slog.Info("保存系统用户完成", "userId", userId)
	}
	return err
}

func (sur *SysUserRepo) Update(user *usercase.SysUser) error {
	sql := `update t_system_user set update_time = now(), username = $1, nickname = $2, email = $3, phone = $4, avatar = $5,
                         roles = $6, sort = $7, status = $8, remark = $9 where user_id = $10`
	result, err := sur.db.Exec(context.Background(), sql, user.Username, user.Nickname, user.Email, user.Phone, user.Avatar, user.Roles,
		user.Sort, user.Status, user.Remark, user.UserId)
	if err == nil {
		slog.Info("更新系统用户完成", "row", result.RowsAffected(), "userId", user.UserId)
	}
	return err
}

func (sur *SysUserRepo) FindUserById(userId uint64) (*usercase.SysUser, error) {
	rows, err := sur.db.Query(context.Background(), `select user_id, username, nickname, email, phone, avatar, 
       roles, create_time, remark from t_system_user where user_id = $1 and delete_at = 0 and status = 0 `, userId)
	if err == nil && rows.Next() {
		defer rows.Close()
		return pgx.RowToAddrOfStructByNameLax[usercase.SysUser](rows)
	}
	return nil, err
}

func (sur *SysUserRepo) Page(query *usercase.SysUserQueryForm) ([]*usercase.SysUser, int64, error) {
	var condition strings.Builder
	condition.WriteString(" where delete_at = 0 ")
	if query.Username != "" {
		condition.WriteString(fmt.Sprintf("and username like '%s' ", "%"+query.Username+"%"))
	}
	if query.Nickname != "" {
		condition.WriteString(fmt.Sprintf("and nickname like '%s' ", "%"+query.Nickname+"%"))
	}
	if query.Phone != "" {
		condition.WriteString(fmt.Sprintf("and phone like '%s' ", "%"+query.Phone+"%"))
	}
	if query.Email != "" {
		condition.WriteString(fmt.Sprintf("and email like '%s' ", "%"+query.Email+"%"))
	}
	if query.RoleId > 0 {
		condition.WriteString(fmt.Sprintf("and %d = ANY (roles) ", query.RoleId))
	}
	if query.CreateTimeBegin != "" {
		condition.WriteString(fmt.Sprintf("and create_time >= '%s' ", query.CreateTimeBegin))
	}
	if query.CreateTimeEnd != "" {
		condition.WriteString(fmt.Sprintf("and create_time <= '%s' ", query.CreateTimeEnd))
	}
	total, err := sur.conditionTotal(condition.String())
	if err != nil {
		return nil, 0, err
	}
	users := make([]*usercase.SysUser, 0)
	if total == 0 {
		return users, total, nil
	}
	offset := tools.ComputeOffset(total, query.Page, query.Size, false)
	condition.WriteString("order by sort, create_time desc limit $1 offset $2 ")
	rows, err := sur.db.Query(context.Background(), `select user_id, username, nickname, email, phone, avatar, roles, 
       last_login_ip, last_login_time, create_time, update_time, sort, status, remark from t_system_user `+condition.String(),
		query.Size, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	users, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (*usercase.SysUser, error) {
		return pgx.RowToAddrOfStructByNameLax[usercase.SysUser](row)
	})
	return users, total, err
}

func (sur *SysUserRepo) DeleteById(userId int64) error {
	result, err := sur.db.Exec(context.Background(), "update t_system_user set delete_at = $1 where user_id = $2", time.Now().UnixMilli(), userId)
	if err == nil {
		slog.Info("删除系统用户完成", "row", result.RowsAffected(), "userId", userId)
	}
	return err
}

func (sur *SysUserRepo) CountByRoleId(roleId int) (int64, error) {
	row := sur.db.QueryRow(context.Background(), "select count(*) from t_system_user where $1 = ANY (roles) and delete_at = 0", roleId)
	var total int64
	err := row.Scan(&total)
	return total, err
}

func (sur *SysUserRepo) CountByUsername(username string, userId uint64) (uint8, error) {
	row := sur.db.QueryRow(context.Background(), "select count(*) from t_system_user where username = $1 and delete_at = 0 and user_id != $2",
		username, userId)
	var total uint8
	err := row.Scan(&total)
	return total, err
}

func (sur *SysUserRepo) QueryUserByUsernameAndPassword(username, password string) (*usercase.SysUser, error) {
	rows, err := sur.db.Query(context.Background(), `select user_id, username, nickname, password, email, phone, avatar, roles, status
		from t_system_user where username = $1 and password = $2 and delete_at = 0`, username, password)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			return pgx.RowToAddrOfStructByNameLax[usercase.SysUser](rows)
		}
	}
	return nil, err
}

func (sur *SysUserRepo) UpdatePassword(userId uint64, newPassword string) error {
	result, err := sur.db.Exec(context.Background(), "update t_system_user set password = $1 where user_id = $2", newPassword, userId)
	if err == nil {
		slog.Info("更新系统用户密码完成", "row", result.RowsAffected(), "userId", userId)
	}
	return nil
}

func (sur *SysUserRepo) UpdateLoginRecord(userId uint64, ip string) {
	_, err := sur.db.Exec(context.Background(), "update t_system_user set last_login_time = now(), last_login_ip = $1 where user_id = $2", ip, userId)
	if err != nil {
		slog.Error("更新系统用户登录记录失败", "err", err)
	}
}

func (sur *SysUserRepo) conditionTotal(condition string) (int64, error) {
	row := sur.db.QueryRow(context.Background(), "select count(*) from t_system_user "+condition)
	var total int64
	err := row.Scan(&total)
	return total, err
}
