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

type SysDictRepo struct {
	db *pgxpool.Pool
}

func NewSysDictRepo(data *Data) usercase.ISysDictRepo {
	return &SysDictRepo{
		db: data.Db,
	}
}

func (self *SysDictRepo) Transaction(ctx context.Context, fn func(tx pgx.Tx) error) error {
	err := pgx.BeginFunc(ctx, self.db, func(tx pgx.Tx) error {
		defer tx.Rollback(ctx)
		err := fn(tx)
		if err == nil {
			if err = tx.Commit(ctx); err != nil {
				slog.Error("pgx事务提交失败，开始回滚", "err", err.Error())
			}
		}
		return err
	})
	return err
}

func (self *SysDictRepo) SaveDict(dict *usercase.SysDict) error {
	sql := "insert into t_system_dict (dict_name, dict_key, sort, status, remark) values ($1, $2, $3, $4, $5) returning dict_id"
	row := self.db.QueryRow(context.Background(), sql, dict.DictName, dict.DictKey, dict.Sort, dict.Status, dict.Remark)
	var dictId uint64
	err := row.Scan(&dictId)
	if err == nil {
		dict.DictId = dictId
		slog.Info("保存系统字典完成", "dictId", dictId)
	}
	return err
}

func (self *SysDictRepo) UpdateDict(dict *usercase.SysDict, tx pgx.Tx) error {
	sql := "update t_system_dict set update_time = now(), dict_name = $1, dict_key = $2, sort = $3, status = $4, remark = $5 where dict_id = $6"
	result, err := smartExec(self.db, tx, context.Background(), sql, dict.DictName, dict.DictKey, dict.Sort, dict.Status, dict.Remark, dict.DictId)
	if err == nil {
		slog.Info("更新系统字典完成", "row", result.RowsAffected(), "dictId", dict.DictId)
	}
	return err
}

func (self *SysDictRepo) UpdateSelectiveDict(dict *usercase.SysDict, tx pgx.Tx) error {
	var sql strings.Builder
	args := make([]any, 0)
	sql.WriteString("update t_system_dict set update_time = now()")
	if dict.DictName != "" {
		args = append(args, dict.DictName)
		sql.WriteString(fmt.Sprintf(", dict_name = $%d", len(args)))
	}
	if dict.DictKey != "" {
		args = append(args, dict.DictKey)
		sql.WriteString(fmt.Sprintf(", dict_key = $%d", len(args)))
	}
	if dict.Remark != "" {
		args = append(args, dict.Remark)
		sql.WriteString(fmt.Sprintf(", remark = $%d", len(args)))
	}
	commonFieldUpdateBuilder(dict.Sort, dict.Status, &sql, &args)
	args = append(args, dict.DictId)
	sql.WriteString(fmt.Sprintf(" where dict_id = $%d", len(args)))
	result, err := smartExec(self.db, tx, context.Background(), sql.String(), args)
	if err == nil {
		slog.Info("更新系统字典完成", "row", result.RowsAffected(), "dictId", dict.DictId)
	}
	return err
}

func (self *SysDictRepo) CountByKey(key string, dictId uint64) (uint8, error) {
	sql := "select count(*) from t_system_dict where dict_key = $1 and dict_id != $2 and delete_at = 0"
	row := self.db.QueryRow(context.Background(), sql, key, dictId)
	var total uint8
	err := row.Scan(&total)
	return total, err
}

func (self *SysDictRepo) PageDict(query *usercase.SysDictQueryForm) ([]*usercase.SysDict, int64, error) {
	var condition strings.Builder
	args := make([]any, 0)
	condition.WriteString(" where delete_at = 0 ")
	if query.DictName != "" {
		args = append(args, "%"+query.DictName+"%")
		condition.WriteString(fmt.Sprintf("and dict_name like $%d ", len(args)))
	}
	if query.DictKey != "" {
		args = append(args, "%"+query.DictKey+"%")
		condition.WriteString(fmt.Sprintf("and dict_key like $%d ", len(args)))
	}
	timeQueryConditionBuilder(query.CreateTimeBegin, query.CreateTimeEnd, &condition, &args)
	row := self.db.QueryRow(context.Background(), "select count(*) from t_system_dict "+condition.String(), args)
	var total int64
	if err := row.Scan(&total); err != nil {
		return nil, 0, err
	}
	dicts := make([]*usercase.SysDict, 0)
	if total == 0 {
		return dicts, 0, nil
	}
	offset := tools.ComputeOffset(total, query.Page, query.Size, false)
	condition.WriteString(fmt.Sprintf(" limit $%d offset $%d", len(args)+1, len(args)+2))
	args = append(args, query.Size, offset)
	rows, err := self.db.Query(context.Background(), "select * from t_system_dict "+condition.String(), args)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	dicts, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (*usercase.SysDict, error) {
		return pgx.RowToAddrOfStructByNameLax[usercase.SysDict](row)
	})
	return dicts, total, err
}

func (self *SysDictRepo) SelectDictById(dictId uint64) (*usercase.SysDict, error) {
	rows, err := self.db.Query(context.Background(), "select * from t_system_dict where dict_id = $1 and delete_at = 0", dictId)
	if err == nil && rows.Next() {
		return pgx.RowToAddrOfStructByNameLax[usercase.SysDict](rows)
	}
	return nil, err
}

func (self *SysDictRepo) DeleteDict(dictId int64, tx pgx.Tx) error {
	sql := "update t_system_dict set delete_at = $1 where dict_id = $2"
	result, err := smartExec(self.db, tx, context.Background(), sql, time.Now().UnixMilli(), dictId)
	if err == nil {
		slog.Info("删除系统字典完成", "row", result.RowsAffected(), "dictId", dictId)
	}
	return err
}

func (self *SysDictRepo) SaveDictValue(value *usercase.SysDictValue) error {
	sql := "insert into t_system_dict_value (dict_id, dict_key, label, value, sort, status, remark) values ($1, $2, $3, $4, $5, $6, $7) returning id"
	row := self.db.QueryRow(context.Background(), sql, value.DictId, value.DictKey, value.Label, value.Value, value.Sort, value.Status, value.Remark)
	var valueId uint64
	err := row.Scan(&valueId)
	if err == nil {
		value.ID = valueId
		slog.Info("保存系统字典数据完成", "valueId", valueId)
	}
	return err
}

func (self *SysDictRepo) UpdateDictValue(value *usercase.SysDictValue) error {
	sql := "update t_system_dict_value set update_time = now(), label = $1, value = $2, sort = $3, status = $4, remark = $5 where id = $6"
	result, err := self.db.Exec(context.Background(), sql, value.Label, value.Value, value.Sort, value.Status, value.Remark, value.ID)
	if err == nil {
		slog.Info("更新系统字典数据完成", "row", result.RowsAffected(), "valueId", value.ID)
	}
	return err
}

func (self *SysDictRepo) UpdateSelectiveDictValue(value *usercase.SysDictValue, tx pgx.Tx) error {
	var sql strings.Builder
	args := make([]any, 0)
	sql.WriteString("update t_system_dict_value set update_time = now()")
	if value.Label != "" {
		args = append(args, value.Label)
		sql.WriteString(fmt.Sprintf(", label = $%d", len(args)))
	}
	if value.Value != "" {
		args = append(args, value.Value)
		sql.WriteString(fmt.Sprintf(", value = $%d", len(args)))
	}
	if value.Remark != "" {
		args = append(args, value.Remark)
		sql.WriteString(fmt.Sprintf(", remakr = $%d", len(args)))
	}
	commonFieldUpdateBuilder(value.Sort, value.Status, &sql, &args)
	args = append(args, value.ID)
	sql.WriteString(fmt.Sprintf(" where id = $%d", len(args)))
	result, err := smartExec(self.db, tx, context.Background(), sql.String(), args)
	if err == nil {
		slog.Info("系统字典数据更新完成", "row", result.RowsAffected(), "valueId", value.ID)
	}
	return err
}

func (self *SysDictRepo) CountValueById(value string, dictId uint64, valueId uint64) (uint8, error) {
	sql := "select count(*) from t_system_dict_value where value = $1 and dict_id = $2 and id != $3 and delete_at = 0"
	row := self.db.QueryRow(context.Background(), sql, value, dictId, valueId)
	var total uint8
	err := row.Scan(&total)
	return total, err
}

func (self *SysDictRepo) PageDictValue(query *usercase.SysDictValueQueryForm) ([]*usercase.SysDictValue, int64, error) {
	var condition strings.Builder
	args := make([]any, 0)
	condition.WriteString(" where delete_at = 0")
	if query.DictId != 0 {
		args = append(args, query.DictId)
		condition.WriteString(fmt.Sprintf(" and dict_id = $%d", len(args)))
	}
	if query.DictKey != "" {
		args = append(args, query.DictKey)
		condition.WriteString(fmt.Sprintf(" and dict_key = $%d", len(args)))
	}
	if query.Label != "" {
		args = append(args, "%"+query.Label+"%")
		condition.WriteString(fmt.Sprintf(" and label link $%d", len(args)))
	}
	timeQueryConditionBuilder(query.CreateTimeBegin, query.CreateTimeEnd, &condition, &args)
	row := self.db.QueryRow(context.Background(), "select count(*) from t_system_dict_value "+condition.String())
	var total int64
	if err := row.Scan(&total); err != nil {
		return nil, 0, err
	}
	values := make([]*usercase.SysDictValue, 0)
	if total == 0 {
		return values, 0, nil
	}
	offset := tools.ComputeOffset(total, query.Page, query.Size, false)
	condition.WriteString(fmt.Sprintf(" limit $%d offset $%d", len(args)+1, len(args)+2))
	args = append(args, query.Size, offset)
	rows, err := self.db.Query(context.Background(), "select * from t_system_dict_value "+condition.String())
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	values, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (*usercase.SysDictValue, error) {
		return pgx.RowToAddrOfStructByName[usercase.SysDictValue](row)
	})
	return values, total, err
}

func (self *SysDictRepo) DeleteDictValue(valueId int64) error {
	sql := "update t_system_dict_value set delete_at = $1 where id = $2"
	result, err := self.db.Exec(context.Background(), sql, time.Now().UnixMilli(), valueId)
	if err == nil {
		slog.Info("删除系统字典数据完成", "row", result.RowsAffected(), "valueId", valueId)
	}
	return err
}

func (self *SysDictRepo) UpdateDictValueByDickId(value *usercase.SysDictValue, tx pgx.Tx) error {
	var sql strings.Builder
	args := make([]any, 0)
	sql.WriteString("update t_system_dict_value set update_time = now() ")
	if value.DictKey != "" {
		args = append(args, value.DictKey)
		sql.WriteString(fmt.Sprintf(", set dict_key = $%d", len(args)))
	}
	if value.Status != nil {
		args = append(args, *value.Status)
		sql.WriteString(fmt.Sprintf(", set status = $%d", len(args)))
	}
	args = append(args, value.DictId)
	sql.WriteString(fmt.Sprintf(" where dict_id = $%d and delete_at = 0", len(args)))
	if *value.Status == 0 {
		sql.WriteString(" and status = 2")
	} else if *value.Status == 2 {
		sql.WriteString(" and status = 0")
	}
	result, err := smartExec(self.db, tx, context.Background(), sql.String(), args)
	if err == nil {
		slog.Info("通过字典Id更新字典数据完成", "row", result.RowsAffected(), "dictId", value.DictId)
	}
	return err
}

func (self *SysDictRepo) DeleteDictValueByDictId(dictId uint64, tx pgx.Tx) error {
	sql := "update t_system_dict_value set delete_at = $1 where dict_id = $2 and delete_at = 0"
	result, err := smartExec(self.db, tx, context.Background(), sql, time.Now().UnixMilli(), dictId)
	if err == nil {
		slog.Info("通过字典Id删除字典数据完成", "row", result.RowsAffected(), "dictId", dictId)
	}
	return err
}
