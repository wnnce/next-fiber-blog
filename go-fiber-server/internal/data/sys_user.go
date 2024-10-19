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
	"strconv"
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

func (self *SysUserRepo) Save(user *usercase.SysUser) error {
	builder := sqlbuild.NewInsertBuilder("t_system_user").
		Fields("username", "nickname", "password", "email", "phone", "avatar", "roles", "sort", "status", "remark").
		Values(user.Username, user.Nickname, user.Password, user.Email, user.Phone, user.Avatar, user.Roles, user.Sort,
			user.Status, user.Remark).
		Returning("user_id")
	row := self.db.QueryRow(context.Background(), builder.Sql(), builder.Args()...)
	var userId uint64
	err := row.Scan(&userId)
	if err == nil {
		user.UserId = userId
		slog.Info("保存系统用户完成", "userId", userId)
	}
	return err
}

func (self *SysUserRepo) Update(user *usercase.SysUser) error {
	builder := sqlbuild.NewUpdateBuilder("t_system_user").
		SetRaw("update_time", "now()").
		SetByMap(map[string]any{
			"username": user.Username,
			"nickname": user.Nickname,
			"email":    user.Email,
			"phone":    user.Phone,
			"avatar":   user.Avatar,
			"roles":    user.Roles,
			"sort":     user.Sort,
			"status":   user.Status,
			"remark":   user.Remark,
		}).Where("user_id").Eq(user.UserId).BuildAsUpdate()
	result, err := self.db.Exec(context.Background(), builder.Sql(), builder.Args()...)
	if err == nil {
		slog.Info("更新系统用户完成", "row", result.RowsAffected(), "userId", user.UserId)
	}
	return err
}

func (self *SysUserRepo) UpdateSelective(form *usercase.SysUserUpdateForm) error {
	builder := sqlbuild.NewUpdateBuilder("t_system_user").
		SetRaw("update_time", "now()").
		SetByCondition(form.Status != nil, "status", form.Status).
		Where("user_id").Eq(form.UserId).BuildAsUpdate()
	result, err := self.db.Exec(context.Background(), builder.Sql(), builder.Args()...)
	if err == nil {
		slog.Info("快捷更新系统用户完成", "row", result.RowsAffected(), "userId", form.UserId)
	}
	return err
}

func (self *SysUserRepo) FindUserById(userId uint64) (*usercase.SysUser, error) {
	builder := sqlbuild.NewSelectBuilder("t_system_user").
		Select("user_id", "username", "password", "nickname", "email", "phone", "avatar", "roles", "create_time", "remark").
		Where("user_id").Eq(userId).
		And("status").EqRaw("0").
		And("delete_at").EqRaw("0").BuildAsSelect()
	rows, err := self.db.Query(context.Background(), builder.Sql(), userId)
	if err == nil && rows.Next() {
		defer rows.Close()
		return pgx.RowToAddrOfStructByNameLax[usercase.SysUser](rows)
	}
	return nil, err
}

func (self *SysUserRepo) Page(query *usercase.SysUserQueryForm) ([]*usercase.SysUser, int64, error) {
	builder := sqlbuild.NewSelectBuilder("t_system_user").
		Select("user_id", "username", "nickname", "email", "phone", "avatar", "roles", "last_login_ip",
			"last_login_time", "create_time", "update_time", "sort", "status", "remark").
		WhereByCondition(query.Username != "", "username").Like("%"+query.Username+"%").
		AndByCondition(query.Nickname != "", "nickname").Like("%"+query.Nickname+"%").
		AndByCondition(query.Phone != "", "phone").Eq(query.Phone).
		AndByCondition(query.Email != "", "email").Eq(query.Email).
		AndByCondition(query.CreateTimeBegin != "", "create_time").Ge(query.CreateTimeBegin).
		AndByCondition(query.CreateTimeEnd != "", "create_time").Le(query.CreateTimeEnd).
		AndByCondition(query.RoleId != nil, fmt.Sprintf("%d", query.RoleId)).EqRaw("ANY (roles)").
		And("delete_at").EqRaw("0").
		BuildAsSelect().
		OrderBy("sort", "create_time desc")
	var total int64
	row := self.db.QueryRow(context.Background(), builder.CountSql(), builder.Args()...)
	if err := row.Scan(&total); err != nil {
		return nil, 0, err
	}
	users := make([]*usercase.SysUser, 0)
	if total == 0 {
		return users, total, nil
	}
	offset := tools.ComputeOffset(total, query.Page, query.Size, false)
	builder.Limit(int64(query.Size)).Offset(offset)
	rows, err := self.db.Query(context.Background(), builder.Sql(), builder.Args()...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	users, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (*usercase.SysUser, error) {
		return pgx.RowToAddrOfStructByNameLax[usercase.SysUser](row)
	})
	return users, total, err
}

func (self *SysUserRepo) DeleteById(userId int64) error {
	builder := sqlbuild.NewUpdateBuilder("t_system_user").
		Set("delete_at", time.Now().UnixMilli()).
		Where("user_id").Eq(userId).BuildAsUpdate()
	result, err := self.db.Exec(context.Background(), builder.Sql(), builder.Args()...)
	if err == nil {
		slog.Info("删除系统用户完成", "row", result.RowsAffected(), "userId", userId)
	}
	return err
}

func (self *SysUserRepo) CountByRoleId(roleId int) (int64, error) {
	builder := sqlbuild.NewSelectBuilder("t_system_user").
		Select("count(*)").
		Where(strconv.Itoa(roleId)).EqRaw("ANY (roles)").
		And("delete_at").EqRaw("0").BuildAsSelect()
	row := self.db.QueryRow(context.Background(), builder.Sql(), roleId)
	var total int64
	err := row.Scan(&total)
	return total, err
}

func (self *SysUserRepo) CountByUsername(username string, userId uint64) (uint8, error) {
	builder := sqlbuild.NewSelectBuilder("t_system_user").
		Select("count(*)").
		Where("user_id").Ne(userId).
		And("username").Eq(username).
		And("delete_at").EqRaw("0").BuildAsSelect()
	row := self.db.QueryRow(context.Background(), builder.Sql(), builder.Args()...)
	var total uint8
	err := row.Scan(&total)
	return total, err
}

func (self *SysUserRepo) QueryUserByUsernameAndPassword(username, password string) (*usercase.SysUser, error) {
	builder := sqlbuild.NewSelectBuilder("t_system_user").
		Select("user_id", "username", "nickname", "password", "email", "phone", "avatar", "roles", "status").
		Where("username").Eq(username).
		And("password").Eq(password).
		And("delete_at").EqRaw("0").BuildAsSelect()
	rows, err := self.db.Query(context.Background(), builder.Sql(), builder.Args()...)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			return pgx.RowToAddrOfStructByNameLax[usercase.SysUser](rows)
		}
	}
	return nil, err
}

func (self *SysUserRepo) UpdatePassword(userId uint64, newPassword string) error {
	builder := sqlbuild.NewUpdateBuilder("t_system_user").
		Set("password", newPassword).
		Where("user_id").Eq(userId).BuildAsUpdate()
	result, err := self.db.Exec(context.Background(), builder.Sql(), builder.Args()...)
	if err == nil {
		slog.Info("更新系统用户密码完成", "row", result.RowsAffected(), "userId", userId)
	}
	return nil
}

func (self *SysUserRepo) UpdateLoginRecord(userId uint64, ip string) {
	builder := sqlbuild.NewUpdateBuilder("t_system_user").
		SetRaw("last_login_time", "now()").
		Set("last_login_ip", ip).
		Where("user_id").Eq(userId).BuildAsUpdate()
	_, err := self.db.Exec(context.Background(), builder.Sql(), builder.Args()...)
	if err != nil {
		slog.Error("更新系统用户登录记录失败", "err", err)
	}
}
