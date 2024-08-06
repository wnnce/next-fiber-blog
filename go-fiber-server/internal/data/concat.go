package data

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-fiber-ent-web-layout/internal/usercase"
	sqlbuild "go-fiber-ent-web-layout/pkg/sql-build"
	"log/slog"
	"time"
)

type ConcatRepo struct {
	db *pgxpool.Pool
}

func NewConcatRepo(data *Data) usercase.IConcatRepo {
	return &ConcatRepo{
		db: data.Db,
	}
}

func (self *ConcatRepo) Save(concat *usercase.Concat) error {
	builder := sqlbuild.NewInsertBuilder("t_blog_concat").
		Fields("name", "logo_url", "target_url", "is_main", "sort", "status").
		Values(concat.Name, concat.LogoUrl, concat.TargetUrl, concat.IsMain, concat.Sort, concat.Status).
		Returning("concat_id")
	row := self.db.QueryRow(context.Background(), builder.Sql(), builder.Args()...)
	var concatId uint
	err := row.Scan(&concatId)
	if err == nil {
		slog.Info("联系方式添加完成", "concatId", concatId)
		concat.ConcatId = concatId
	}
	return err
}

func (self *ConcatRepo) Update(concat *usercase.Concat) error {
	builder := sqlbuild.NewUpdateBuilder("t_blog_concat").
		SetRaw("update_time", "now()").
		SetByMap(map[string]any{
			"name":       concat.Name,
			"logo_url":   concat.LogoUrl,
			"target_url": concat.TargetUrl,
			"is_main":    concat.IsMain,
			"sort":       concat.Sort,
			"status":     concat.Status,
		}).Where("concat_id").Eq(concat.ConcatId).BuildAsUpdate()
	result, err := self.db.Exec(context.Background(), builder.Sql(), builder.Args()...)
	if err == nil {
		slog.Info("联系方式更新完成", "row", result.RowsAffected(), "concatId", concat.ConcatId)
	}
	return err
}

func (self *ConcatRepo) UpdateSelective(form *usercase.ConcatUpdateForm) error {
	builder := sqlbuild.NewUpdateBuilder("t_blog_concat").
		SetRaw("update_time", "now()").
		SetByCondition(form.Status != nil, "status", form.Status).
		SetByCondition(form.IsMain != nil, "is_main", form.IsMain).
		Where("concat_id").Eq(form.ConcatId).BuildAsUpdate()
	result, err := self.db.Exec(context.Background(), builder.Sql(), builder.Args()...)
	if err == nil {
		slog.Info("联系方式快捷更新完成", "row", result.RowsAffected(), "concatId", form.ConcatId)
	}
	return err
}

func (self *ConcatRepo) List() ([]*usercase.Concat, error) {
	builder := sqlbuild.NewSelectBuilder("t_blog_concat").
		Select("concat_id", "name", "logo_url", "target_url", "is_main").
		Where("status").EqRaw("0").
		And("delete_at").EqRaw("0").BuildAsSelect().
		OrderBy("sort", "create_time desc")
	rows, err := self.db.Query(context.Background(), builder.Sql())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (*usercase.Concat, error) {
		return pgx.RowToAddrOfStructByNameLax[usercase.Concat](row)
	})
}

func (self *ConcatRepo) ManageList(query *usercase.ConcatQueryForm) ([]*usercase.Concat, error) {
	builder := sqlbuild.NewSelectBuilder("t_blog_concat").
		WhereByCondition(query.Name != "", "name").Like("%"+query.Name+"%").
		AndByCondition(query.CreateTimeBegin != "", "create_time").Ge(query.CreateTimeBegin).
		AndByCondition(query.CreateTimeEnd != "", "create_time").Le(query.CreateTimeEnd).
		And("delete_at").EqRaw("0").BuildAsSelect().
		OrderBy("sort", "create_time desc")
	rows, err := self.db.Query(context.Background(), builder.Sql(), builder.Args()...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (*usercase.Concat, error) {
		return pgx.RowToAddrOfStructByName[usercase.Concat](rows)
	})
}

func (self *ConcatRepo) CountByName(name string, cid uint) (uint8, error) {
	builder := sqlbuild.NewSelectBuilder("t_blog_concat").
		Select("count(*)").
		Where("concat_id").Ne(cid).
		And("name").Eq(name).
		And("delete_at").EqRaw("0").BuildAsSelect()
	row := self.db.QueryRow(context.Background(), builder.Sql(), builder.Args()...)
	var total uint8
	err := row.Scan(&total)
	return total, err
}

func (self *ConcatRepo) DeleteById(cid int) error {
	builder := sqlbuild.NewUpdateBuilder("t_blog_concat").
		Set("delete_at", time.Now().UnixMilli()).
		Where("concat_id").Eq(cid).BuildAsUpdate()
	result, err := self.db.Exec(context.Background(), builder.Sql(), builder.Args()...)
	if err == nil {
		slog.Info(fmt.Sprintf("删除联系方式完成，concatId:%d,row:%d", cid, result.RowsAffected()))
	}
	return err
}
