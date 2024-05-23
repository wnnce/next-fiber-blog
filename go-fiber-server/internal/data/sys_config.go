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

type SysConfigRepo struct {
	db *pgxpool.Pool
}

func NewSysConfigRepo(data *Data) usercase.ISysConfigRepo {
	return &SysConfigRepo{
		db: data.Db,
	}
}

func (sc *SysConfigRepo) Save(cfg *usercase.SysConfig) error {
	row := sc.db.QueryRow(context.Background(), "insert into t_system_config (config_name, config_key, config_value, remark) values ($1, $2, $3, $4) returning config_id",
		cfg.ConfigName, cfg.ConfigKey, cfg.ConfigValue, cfg.Remark)
	var configId uint64
	err := row.Scan(&configId)
	if err == nil {
		cfg.ConfigId = configId
		slog.Info(fmt.Sprintf("系统配置添加完成，id:%d", configId))
	}
	return err
}

func (sc *SysConfigRepo) Update(cfg *usercase.SysConfig) error {
	result, err := sc.db.Exec(context.Background(), "update t_system_config set update_time = now(), config_name = $1, config_key = $2, config_value = $3, remark = $4 where config_id = $5",
		cfg.ConfigName, cfg.ConfigKey, cfg.ConfigValue, cfg.Remark, cfg.ConfigId)
	if err == nil {
		slog.Info(fmt.Sprintf("系统配置更新完成，row:%d,id:%d", result.RowsAffected(), cfg.ConfigId))
	}
	return err
}

func (sc *SysConfigRepo) CountByKey(key string, cid uint64) (uint8, error) {
	row := sc.db.QueryRow(context.Background(), "select count(config_id) from t_system_config where delete_at = 0 and config_key = $1 and config_id != $2", key, cid)
	var total uint8
	err := row.Scan(&total)
	return total, err
}

func (sc *SysConfigRepo) ManagePage(query *usercase.SysConfigQueryForm) ([]usercase.SysConfig, int64, error) {
	var condition strings.Builder
	condition.WriteString("where delete_at = 0 ")
	if query.Name != "" {
		condition.WriteString(fmt.Sprintf("and config_name like '%s' ", "%"+query.Name+"%"))
	}
	if query.Key != "" {
		condition.WriteString(fmt.Sprintf("and config_key like '%s' ", "%"+query.Key+"%"))
	}
	if query.CreateTimeBegin != nil {
		condition.WriteString(fmt.Sprintf("and create_time >= '%s' ", query.CreateTimeBegin.Format("2006-04-02")))
	}
	if query.CreateTimeEnd != nil {
		condition.WriteString(fmt.Sprintf("and create_time <= '%s' ", query.CreateTimeEnd.Format("2006-04-02")))
	}
	total, err := sc.conditionTotal(condition.String())
	if err != nil {
		return nil, 0, err
	}
	configs := make([]usercase.SysConfig, 0)
	if total == 0 {
		return configs, total, nil
	}
	offset := tools.ComputeOffset(total, query.Page, query.Size, false)
	condition.WriteString("order by create_time desc limit $1 offset $2")
	rows, err := sc.db.Query(context.Background(), "select * from t_system_config "+condition.String(), query.Size, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	configs, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (usercase.SysConfig, error) {
		return pgx.RowToStructByName[usercase.SysConfig](row)
	})
	return configs, total, err
}

func (sc *SysConfigRepo) DeleteById(cid int64) error {
	result, err := sc.db.Exec(context.Background(), "update t_system_config set delete_at = $1 where config_id = $2", time.Now().UnixMilli(), cid)
	if err == nil {
		slog.Info(fmt.Sprintf("删除系统配置完成，row:%d,id:%d", result.RowsAffected(), cid))
	}
	return err
}

func (sc *SysConfigRepo) conditionTotal(condition string) (int64, error) {
	row := sc.db.QueryRow(context.Background(), "select count(*) from t_system_config "+condition)
	var total int64
	err := row.Scan(&total)
	return total, err
}
