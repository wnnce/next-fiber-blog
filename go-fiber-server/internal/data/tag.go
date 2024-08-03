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

type TagRepo struct {
	db *pgxpool.Pool
}

func NewTagRepo(data *Data) usercase.ITagRepo {
	return &TagRepo{
		db: data.Db,
	}
}

func (self *TagRepo) Save(form *usercase.TagForm) error {
	builder := sqlbuild.NewInsertBuilder("t_blog_tag").
		Fields("tag_name", "cover_url", "color", "sort", "status").
		Values(form.TagName, form.CoverUrl, form.Color, *form.Sort, *form.Status).
		Returning("tag_id")
	row := self.db.QueryRow(context.Background(), builder.Sql(), builder.Args())
	var insertId int
	err := row.Scan(&insertId)
	if err == nil {
		slog.Info("新增标签完成，id：" + strconv.Itoa(insertId))
	}
	return err
}

func (self *TagRepo) Update(form *usercase.TagForm) error {
	builder := sqlbuild.NewUpdateBuilder("t_blog_tag").
		SetRaw("update_time", "now()").
		SetByMap(map[string]any{
			"tag_name":  form.TagName,
			"cover_url": form.CoverUrl,
			"color":     form.Color,
			"sort":      *form.Sort,
			"status":    *form.Status,
		}).
		Where("tag_id").Eq(form.TagId).BuildAsUpdate()
	result, err := self.db.Exec(context.Background(), builder.Sql(), builder.Args()...)
	if err == nil {
		slog.Info(fmt.Sprintf("标签更新完成，row:%d,id:%d", result.RowsAffected(), form.TagId))
	}
	return err
}

func (self *TagRepo) UpdateSelective(form *usercase.TagUpdateForm) error {
	builder := sqlbuild.NewUpdateBuilder("t_blog_tag").
		SetRaw("update_time", "now()").
		SetByCondition(form.Status != nil, "status", form.Status).
		Where("tag_id").Eq(form.TagId).BuildAsUpdate()
	result, err := self.db.Exec(context.Background(), builder.Sql(), builder.Args()...)
	if err == nil {
		slog.Info("标签快捷更新完成", "row", result.RowsAffected(), "tagId", form.TagId)
	}
	return err
}

func (self *TagRepo) UpdateViewNum(tagId int, addNum int) error {
	builder := sqlbuild.NewUpdateBuilder("t_blog_tag").
		SetRaw("update_time", "now()").
		SetRaw("view_num", "view_num"+strconv.Itoa(addNum)).
		Where("tag_id").Eq(tagId).BuildAsUpdate()
	result, err := self.db.Exec(context.Background(), builder.Sql(), tagId)
	if err == nil {
		slog.Info(fmt.Sprintf("更新标签查看次数完成，row:%d,id:%d,addnum:%d", result.RowsAffected(), tagId, addNum))
	}
	return err
}

func (self *TagRepo) SelectById(id int) (*usercase.Tag, error) {
	builder := sqlbuild.NewSelectBuilder("t_blog_tag").
		Where("tag_id").Eq(id).
		And("status").EqRaw("0").
		And("delete_at").EqRaw("0").BuildAsSelect()
	rows, err := self.db.Query(context.Background(), builder.Sql(), id)
	if err == nil && rows.Next() {
		defer rows.Close()
		return pgx.RowToAddrOfStructByName[usercase.Tag](rows)
	}
	return nil, err
}

func (self *TagRepo) Page(query *usercase.TagQueryForm) ([]*usercase.Tag, int64, error) {
	builder := sqlbuild.NewSelectBuilder("t_blog_tag").
		WhereByCondition(query.TagName != "", "tag_name").Like("%"+query.TagName+"%").
		AndByCondition(query.CreateTimeBegin != "", "create_time").Ge(query.CreateTimeBegin).
		AndByCondition(query.CreateTimeEnd != "", "create_time").Le(query.CreateTimeEnd).
		And("delete_at").EqRaw("0").BuildAsSelect().
		OrderBy("sort", "create_time desc")
	var total int64
	row := self.db.QueryRow(context.Background(), builder.CountSql(), builder.Args()...)
	if err := row.Scan(&total); err != nil {
		return nil, 0, err
	}
	tags := make([]*usercase.Tag, 0)
	if total == 0 {
		return tags, total, nil
	}
	offset := tools.ComputeOffset(total, query.Page, query.Size, false)
	builder.Limit(int64(query.Size)).Offset(offset)
	rows, err := self.db.Query(context.Background(), builder.Sql(), builder.Args()...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	tags, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (*usercase.Tag, error) {
		return pgx.RowToAddrOfStructByName[usercase.Tag](row)
	})
	return tags, total, err
}

func (self *TagRepo) List() ([]*usercase.Tag, error) {
	builder := sqlbuild.NewSelectBuilder("t_blog_tag").
		Select("tag_id", "tag_name", "cover_url", "view_num", "color").
		Where("status").EqRaw("0").
		And("delete_at").EqRaw("0").BuildAsSelect().
		OrderBy("sort", "create_time desc")
	rows, err := self.db.Query(context.Background(), builder.Sql())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (*usercase.Tag, error) {
		return pgx.RowToAddrOfStructByNameLax[usercase.Tag](row)
	})
}

func (self *TagRepo) ListByIds(ids []uint) ([]*usercase.Tag, error) {
	if len(ids) == 0 {
		return make([]*usercase.Tag, 0), nil
	}
	builder := sqlbuild.NewSelectBuilder("t_blog_tag").
		Select("tag_id", "tag_name", "color").
		Where("tag_id").In(sqlbuild.SliceToAnySlice[uint](ids)...).
		And("status").EqRaw("0").
		And("delete_at").EqRaw("0").BuildAsSelect()
	rows, err := self.db.Query(context.Background(), builder.Sql(), builder.Args()...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (*usercase.Tag, error) {
		return pgx.RowToAddrOfStructByNameLax[usercase.Tag](row)
	})
}

func (self *TagRepo) CountByTagName(name string, tagId uint) (uint8, error) {
	builder := sqlbuild.NewSelectBuilder("t_blog_tag").
		Select("count(*)").
		Where("tag_id").Ne(tagId).
		And("tag_name").Eq(name).
		And("delete_at").EqRaw("0").BuildAsSelect()
	row := self.db.QueryRow(context.Background(), builder.Sql(), builder.Args()...)
	var total uint8
	err := row.Scan(&total)
	return total, err
}

func (self *TagRepo) DeleteById(id int) error {
	builder := sqlbuild.NewUpdateBuilder("t_blog_tag").
		Set("delete_at", time.Now().UnixMilli()).
		Where("tag_id").Eq(id).BuildAsUpdate()
	result, err := self.db.Exec(context.Background(), builder.Sql(), builder.Args()...)
	if err == nil {
		slog.Info(fmt.Sprintf("删除标签完成，row：%d,id:%d", result.RowsAffected(), id))
	}
	return err
}

func (self *TagRepo) DeleteByIds(ids []int) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	builder := sqlbuild.NewUpdateBuilder("t_blog_tag").
		Set("delete_at", time.Now().UnixMilli()).
		Where("tag_id").In(sqlbuild.SliceToAnySlice[int](ids)...).BuildAsUpdate()
	result, err := self.db.Exec(context.Background(), builder.Sql(), builder.Args()...)
	return result.RowsAffected(), err
}
