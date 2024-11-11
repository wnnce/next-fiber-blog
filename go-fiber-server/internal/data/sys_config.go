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

type SysConfigRepo struct {
	db *pgxpool.Pool
}

func NewSysConfigRepo(data *Data) usercase.ISysConfigRepo {
	return &SysConfigRepo{
		db: data.Db,
	}
}

func (self *SysConfigRepo) Save(cfg *usercase.SysConfig) error {
	builder := sqlbuild.NewInsertBuilder("t_system_config").
		Fields("config_name", "config_key", "config_value", "remark").
		Values(cfg.ConfigName, cfg.ConfigKey, cfg.ConfigValue, cfg.Remark).
		Returning("config_id")
	row := self.db.QueryRow(context.Background(), builder.Sql(), builder.Args()...)
	var configId uint64
	err := row.Scan(&configId)
	if err == nil {
		cfg.ConfigId = configId
		slog.Info(fmt.Sprintf("系统配置添加完成，id:%d", configId))
	}
	return err
}

func (self *SysConfigRepo) Update(cfg *usercase.SysConfig) error {
	builder := sqlbuild.NewUpdateBuilder("t_system_config").
		SetRaw("update_time", "now()").
		SetByMap(map[string]any{
			"config_name":  cfg.ConfigName,
			"config_key":   cfg.ConfigKey,
			"config_value": cfg.ConfigValue,
			"remark":       cfg.Remark,
		}).Where("config_id").Eq(cfg.ConfigId).BuildAsUpdate()
	result, err := self.db.Exec(context.Background(), builder.Sql(), builder.Args()...)
	if err == nil {
		slog.Info(fmt.Sprintf("系统配置更新完成，row:%d,id:%d", result.RowsAffected(), cfg.ConfigId))
	}
	return err
}

func (self *SysConfigRepo) CountByKey(key string, cid uint64) (uint8, error) {
	builder := sqlbuild.NewSelectBuilder("t_system_config").
		Select("count(*)").
		Where("config_id").Ne(cid).
		And("config_key").Eq(key).And("delete_at").EqRaw("0").BuildAsSelect()
	row := self.db.QueryRow(context.Background(), builder.Sql(), builder.Args()...)
	var total uint8
	err := row.Scan(&total)
	return total, err
}

func (self *SysConfigRepo) ManagePage(query *usercase.SysConfigQueryForm) ([]usercase.SysConfig, int64, error) {
	builder := sqlbuild.NewSelectBuilder("t_system_config").
		WhereByCondition(query.Name != "", "config_name").Like("%"+query.Name+"%").
		AndByCondition(query.Key != "", "config_key").Like("%"+query.Key+"%").
		AndByCondition(query.CreateTimeBegin != "", "create_time").Ge(query.CreateTimeBegin).
		AndByCondition(query.CreateTimeEnd != "", "create_time").Le(query.CreateTimeEnd).BuildAsSelect().
		OrderByDesc("create_time")
	var total int64
	row := self.db.QueryRow(context.Background(), builder.CountSql(), builder.Args()...)
	if err := row.Scan(&total); err != nil {
		return nil, 0, err
	}
	configs := make([]usercase.SysConfig, 0)
	if total == 0 {
		return configs, total, nil
	}
	offset := tools.ComputeOffset(total, query.Page, query.Size, false)
	builder.Limit(int64(query.Size)).Offset(offset)
	rows, err := self.db.Query(context.Background(), builder.Sql(), builder.Args()...)
	defer rows.Close()
	if err != nil {
		return nil, 0, err
	}
	configs, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (usercase.SysConfig, error) {
		return pgx.RowToStructByName[usercase.SysConfig](row)
	})
	return configs, total, err
}

func (self *SysConfigRepo) DeleteById(cid int64) error {
	builder := sqlbuild.NewUpdateBuilder("t_system_config").
		Set("delete_at", time.Now().UnixMilli()).
		Where("config_id").Eq(cid).BuildAsUpdate()
	result, err := self.db.Exec(context.Background(), builder.Sql(), builder.Args()...)
	if err == nil {
		slog.Info(fmt.Sprintf("删除系统配置完成，row:%d,id:%d", result.RowsAffected(), cid))
	}
	return err
}
