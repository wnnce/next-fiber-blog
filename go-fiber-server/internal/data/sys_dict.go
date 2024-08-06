package data

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-fiber-ent-web-layout/internal/tools"
	"go-fiber-ent-web-layout/internal/usercase"
	sqlbuild "go-fiber-ent-web-layout/pkg/sql-build"
	"log/slog"
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
	return pgx.BeginFunc(ctx, self.db, func(tx pgx.Tx) error {
		return fn(tx)
	})
}

func (self *SysDictRepo) SaveDict(dict *usercase.SysDict) error {
	builder := sqlbuild.NewInsertBuilder("t_system_dict").
		Fields("dict_name", "dict_key", "sort", "status", "remark").
		Values(dict.DictName, dict.DictKey, dict.Sort, dict.Status, dict.Remark).
		Returning("dict_id")
	row := self.db.QueryRow(context.Background(), builder.Sql(), builder.Args()...)
	var dictId uint64
	err := row.Scan(&dictId)
	if err == nil {
		dict.DictId = dictId
		slog.Info("保存系统字典完成", "dictId", dictId)
	}
	return err
}

func (self *SysDictRepo) UpdateDict(dict *usercase.SysDict, tx pgx.Tx) error {
	builder := sqlbuild.NewUpdateBuilder("t_system_dict").
		SetRaw("update_time", "now()").
		SetByMap(map[string]any{
			"dict_name": dict.DictName,
			"dict_key":  dict.DictKey,
			"sort":      dict.Sort,
			"status":    dict.Status,
			"remark":    dict.Remark,
		}).Where("dict_id").Eq(dict.DictId).BuildAsUpdate()
	result, err := smartExec(self.db, tx, context.Background(), builder.Sql(), builder.Args()...)
	if err == nil {
		slog.Info("更新系统字典完成", "row", result.RowsAffected(), "dictId", dict.DictId)
	}
	return err
}

func (self *SysDictRepo) UpdateSelectiveDict(dict *usercase.SysDict, tx pgx.Tx) error {
	builder := sqlbuild.NewUpdateBuilder("t_system_dict").
		SetRaw("update_time", "now()").
		SetByCondition(dict.DictName != "", "dict_name", dict.DictName).
		SetByCondition(dict.DictKey != "", "dict_key", dict.DictKey).
		SetByCondition(dict.Remark != "", "remark", dict.Remark).
		SetByCondition(dict.Sort != nil, "sort", dict.Sort).
		SetByCondition(dict.Status != nil, "status", dict.Status).
		Where("dict_id").Eq(dict.DictId).BuildAsUpdate()
	result, err := smartExec(self.db, tx, context.Background(), builder.Sql(), builder.Args()...)
	if err == nil {
		slog.Info("更新系统字典完成", "row", result.RowsAffected(), "dictId", dict.DictId)
	}
	return err
}

func (self *SysDictRepo) CountByKey(key string, dictId uint64) (uint8, error) {
	builder := sqlbuild.NewSelectBuilder("t_system_dict").
		Where("dict_id").Ne(dictId).
		And("dict_key").Eq(key).
		And("delete_at").EqRaw("0").BuildAsSelect()
	row := self.db.QueryRow(context.Background(), builder.Sql(), builder.Args()...)
	var total uint8
	err := row.Scan(&total)
	return total, err
}

func (self *SysDictRepo) PageDict(query *usercase.SysDictQueryForm) ([]*usercase.SysDict, int64, error) {
	builder := sqlbuild.NewSelectBuilder("t_system_dict").
		WhereByCondition(query.DictName != "", "dict_name").Like("%"+query.DictName+"%").
		AndByCondition(query.DictKey != "", "dict_key").Like("%"+query.DictKey+"%").
		AndByCondition(query.CreateTimeBegin != "", "create_time").Ge(query.CreateTimeBegin).
		AndByCondition(query.CreateTimeEnd != "", "create_time").Le(query.CreateTimeEnd).
		And("delete_at").EqRaw("0").BuildAsSelect().
		OrderBy("sort", "create_time desc")
	row := self.db.QueryRow(context.Background(), builder.CountSql(), builder.Args()...)
	var total int64
	if err := row.Scan(&total); err != nil {
		return nil, 0, err
	}
	dicts := make([]*usercase.SysDict, 0)
	if total == 0 {
		return dicts, 0, nil
	}
	offset := tools.ComputeOffset(total, query.Page, query.Size, false)
	builder.Limit(int64(query.Size)).Offset(offset)
	rows, err := self.db.Query(context.Background(), builder.Sql(), builder.Args()...)
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
	builder := sqlbuild.NewSelectBuilder("t_system_dict").
		Where("dict_id").Eq(dictId).
		And("delete_at").EqRaw("0").BuildAsSelect()
	rows, err := self.db.Query(context.Background(), builder.Sql(), dictId)
	if err == nil && rows.Next() {
		return pgx.RowToAddrOfStructByNameLax[usercase.SysDict](rows)
	}
	return nil, err
}

func (self *SysDictRepo) DeleteDict(dictId int64, tx pgx.Tx) error {
	builder := sqlbuild.NewUpdateBuilder("t_system_dict").
		Set("delete_at", time.Now().UnixMilli()).
		Where("dict_id").Eq(dictId).BuildAsUpdate()
	result, err := smartExec(self.db, tx, context.Background(), builder.Sql(), builder.Args()...)
	if err == nil {
		slog.Info("删除系统字典完成", "row", result.RowsAffected(), "dictId", dictId)
	}
	return err
}

func (self *SysDictRepo) SaveDictValue(value *usercase.SysDictValue) error {
	builder := sqlbuild.NewInsertBuilder("t_system_dict_value").
		Fields("dict_id", "dict_key", "label", "value", "sort", "status", "remark").
		Values(value.DictId, value.DictKey, value.Label, value.Value, value.Sort, value.Status, value.Remark).
		Returning("id")
	row := self.db.QueryRow(context.Background(), builder.Sql(), builder.Args()...)
	var valueId uint64
	err := row.Scan(&valueId)
	if err == nil {
		value.ID = valueId
		slog.Info("保存系统字典数据完成", "valueId", valueId)
	}
	return err
}

func (self *SysDictRepo) UpdateDictValue(value *usercase.SysDictValue) error {
	builder := sqlbuild.NewUpdateBuilder("t_system_dict_value").
		SetRaw("update_time", "now()").
		SetByMap(map[string]any{
			"label":  value.Label,
			"value":  value.Value,
			"sort":   value.Sort,
			"status": value.Status,
			"remark": value.Remark,
		}).Where("id").Eq(value.ID).BuildAsUpdate()
	result, err := self.db.Exec(context.Background(), builder.Sql(), builder.Args()...)
	if err == nil {
		slog.Info("更新系统字典数据完成", "row", result.RowsAffected(), "valueId", value.ID)
	}
	return err
}

func (self *SysDictRepo) UpdateSelectiveDictValue(value *usercase.SysDictValue, tx pgx.Tx) error {
	builder := sqlbuild.NewUpdateBuilder("t_system_dict_value").
		SetRaw("update_time", "now()").
		SetByCondition(value.Label != "", "label", value.Label).
		SetByCondition(value.Value != "", "value", value.Value).
		SetByCondition(value.Remark != "", "remark", value.Remark).
		SetByCondition(value.Sort != nil, "sort", value.Sort).
		SetByCondition(value.Status != nil, "status", value.Status).
		Where("id").Eq(value.ID).BuildAsUpdate()
	result, err := smartExec(self.db, tx, context.Background(), builder.Sql(), builder.Args()...)
	if err == nil {
		slog.Info("系统字典数据更新完成", "row", result.RowsAffected(), "valueId", value.ID)
	}
	return err
}

func (self *SysDictRepo) CountValueById(value string, dictId uint64, valueId uint64) (uint8, error) {
	builder := sqlbuild.NewSelectBuilder("t_system_dict_value").
		Select("count(*)").
		Where("id").Ne(valueId).And("dict_id").Ne(dictId).
		And("value").Eq(value).And("delete_at").EqRaw("0").BuildAsSelect()
	row := self.db.QueryRow(context.Background(), builder.Sql(), builder.Args()...)
	var total uint8
	err := row.Scan(&total)
	return total, err
}

func (self *SysDictRepo) SelectDictKeyById(dictId int64) string {
	builder := sqlbuild.NewSelectBuilder("t_system_dict").
		Select("dict_key").
		Where("dict_id").Eq(dictId).BuildAsSelect()
	row := self.db.QueryRow(context.Background(), builder.Sql(), dictId)
	var dictKey string
	_ = row.Scan(&dictKey)
	return dictKey
}

func (self *SysDictRepo) SelectDictKeyByValueId(valueId int64) string {
	builder := sqlbuild.NewSelectBuilder("t_system_dict_value").
		Select("dict_key").
		Where("id").Eq(valueId).BuildAsSelect()
	row := self.db.QueryRow(context.Background(), builder.Sql(), valueId)
	var dictKey string
	_ = row.Scan(&dictKey)
	return dictKey
}

func (self *SysDictRepo) PageDictValue(query *usercase.SysDictValueQueryForm) ([]*usercase.SysDictValue, int64, error) {
	builder := sqlbuild.NewSelectBuilder("t_system_dict_value").
		WhereByCondition(query.DictId != 0, "dict_id").Eq(query.DictId).
		AndByCondition(query.DictKey != "", "dict_key").Eq(query.DictKey).
		AndByCondition(query.Label != "", "label").Like("%"+query.Label+"%").
		AndByCondition(query.CreateTimeBegin != "", "create_time").Ge(query.CreateTimeBegin).
		AndByCondition(query.CreateTimeEnd != "", "create_time").Le(query.CreateTimeEnd).
		And("delete_at").EqRaw("0").BuildAsSelect().
		OrderBy("sort", "create_time desc")
	row := self.db.QueryRow(context.Background(), builder.CountSql(), builder.Args()...)
	var total int64
	if err := row.Scan(&total); err != nil {
		return nil, 0, err
	}
	values := make([]*usercase.SysDictValue, 0)
	if total == 0 {
		return values, 0, nil
	}
	offset := tools.ComputeOffset(total, query.Page, query.Size, false)
	builder.Limit(int64(query.Size)).Offset(offset)
	rows, err := self.db.Query(context.Background(), builder.Sql(), builder.Args()...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	values, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (*usercase.SysDictValue, error) {
		return pgx.RowToAddrOfStructByName[usercase.SysDictValue](row)
	})
	return values, total, err
}

func (self *SysDictRepo) ListDictValueByDictKey(dictKey string) ([]usercase.SysDictValue, error) {
	builder := sqlbuild.NewSelectBuilder("t_system_dict_value").
		Select("id", "dict_key", "label", "value").
		Where("dict_key").Eq(dictKey).
		And("delete_at").EqRaw("0").BuildAsSelect().
		OrderBy("sort")
	rows, err := self.db.Query(context.Background(), builder.Sql(), dictKey)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (usercase.SysDictValue, error) {
		return pgx.RowToStructByNameLax[usercase.SysDictValue](row)
	})
}

func (self *SysDictRepo) DeleteDictValue(valueId int64) error {
	builder := sqlbuild.NewUpdateBuilder("t_system_dict_value").
		Set("delete_at", time.Now().UnixMilli()).
		Where("id").Eq(valueId).BuildAsUpdate()
	result, err := self.db.Exec(context.Background(), builder.Sql(), builder.Args()...)
	if err == nil {
		slog.Info("删除系统字典数据完成", "row", result.RowsAffected(), "valueId", valueId)
	}
	return err
}

func (self *SysDictRepo) UpdateDictValueByDickId(value *usercase.SysDictValue, tx pgx.Tx) error {
	var newStatus int
	if value.Status != nil && *value.Status == 0 {
		newStatus = 2
	}
	builder := sqlbuild.NewUpdateBuilder("t_system_dict_value").
		Set("update_time", "now()").
		SetByCondition(value.DictKey != "", "dict_key", value.DictKey).
		SetByCondition(value.Status != nil, "status", value.Status).
		Where("dict_id").Eq(value.DictId).
		And("delete_at").EqRaw("0").
		And("status").Eq(newStatus).BuildAsUpdate()
	result, err := smartExec(self.db, tx, context.Background(), builder.Sql(), builder.Args()...)
	if err == nil {
		slog.Info("通过字典Id更新字典数据完成", "row", result.RowsAffected(), "dictId", value.DictId)
	}
	return err
}

func (self *SysDictRepo) DeleteDictValueByDictId(dictId int64, tx pgx.Tx) error {
	builder := sqlbuild.NewUpdateBuilder("t_system_dict_value").
		Set("delete_at", time.Now().UnixMilli()).
		Where("dict_id").Eq(dictId).
		And("delete_at").EqRaw("0").BuildAsUpdate()
	result, err := smartExec(self.db, tx, context.Background(), builder.Sql(), builder.Args()...)
	if err == nil {
		slog.Info("通过字典Id删除字典数据完成", "row", result.RowsAffected(), "dictId", dictId)
	}
	return err
}
