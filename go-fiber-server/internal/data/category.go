package data

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-fiber-ent-web-layout/internal/usercase"
	sqlbuild "go-fiber-ent-web-layout/pkg/sql-build"
	"log/slog"
	"strconv"
	"time"
)

type CategoryRepo struct {
	db *pgxpool.Pool
}

func NewCategoryRepo(data *Data) usercase.ICategoryRepo {
	return &CategoryRepo{
		db: data.Db,
	}
}

func (c *CategoryRepo) Save(cat *usercase.Category) error {
	builder := sqlbuild.NewInsertBuilder("t_blog_category").
		Fields("category_name", "description", "cover_url", "parent_id", "is_top", "is_hot", "sort", "status").
		Values(cat.CategoryName, cat.Description, cat.CoverUrl, cat.ParentId, cat.IsTop, cat.IsHot, cat.Sort, cat.Status).
		Returning("category_id")
	row := c.db.QueryRow(context.Background(), builder.Sql(), builder.Args()...)
	var categoryId uint
	err := row.Scan(&categoryId)
	if err == nil {
		cat.CategoryId = categoryId
		slog.Info("分类保存完成，id：" + strconv.Itoa(int(categoryId)))
	}
	return err
}

func (c *CategoryRepo) Update(cat *usercase.Category) error {
	builder := sqlbuild.NewUpdateBuilder("t_blog_category").
		SetRaw("update_time", "now()").
		SetByMap(map[string]any{
			"category_name": cat.CategoryName,
			"description":   cat.Description,
			"cover_url":     cat.CoverUrl,
			"parent_id":     cat.ParentId,
			"is_hot":        cat.IsHot,
			"is_top":        cat.IsTop,
			"sort":          cat.Sort,
			"status":        cat.Status,
		}).Where("category_id").Eq(cat.CategoryId).BuildAsUpdate()
	result, err := c.db.Exec(context.Background(), builder.Sql(), builder.Args()...)
	if err == nil {
		slog.Info(fmt.Sprintf("分类更新完成，row:%d,id:%d", result.RowsAffected(), cat.CategoryId))
	}
	return err
}

func (c *CategoryRepo) UpdateSelective(form *usercase.CategoryUpdateForm) error {
	builder := sqlbuild.NewUpdateBuilder("t_blog_category").
		SetRaw("update_time", "now()").
		SetByCondition(form.Status != nil, "status", form.Status).
		SetByCondition(form.IsHot != nil, "is_hot", form.IsHot).
		SetByCondition(form.IsTop != nil, "is_top", form.IsTop).
		Where("category_id").Eq(form.CategoryId).BuildAsUpdate()
	result, err := c.db.Exec(context.Background(), builder.Sql(), builder.Args()...)
	if err == nil {
		slog.Info("分类快捷更新完成", "row", result.RowsAffected(), "categoryId", form.CategoryId)
	}
	return err
}

func (c *CategoryRepo) UpdateViewNum(catId uint, addNum int) error {
	builder := sqlbuild.NewUpdateBuilder("t_blog_category").
		SetRaw("update_time", "now()").
		SetRaw("view_num", "view_num + "+strconv.Itoa(addNum)).
		Where("category_id").Eq(catId).BuildAsUpdate()
	result, err := c.db.Exec(context.Background(), builder.Sql(), catId)
	if err == nil {
		slog.Info(fmt.Sprintf("更新分类查看次数完成，tagId:%d,addNum:%d,row:%d", catId, addNum, result.RowsAffected()))
	}
	return err
}

func (c *CategoryRepo) SelectById(catId int) (*usercase.Category, error) {
	builder := sqlbuild.NewSelectBuilder("t_blog_category").
		Where("category_id").Eq(catId).
		And("status").EqRaw("0").
		And("delete_at").EqRaw("0").BuildAsSelect()
	rows, err := c.db.Query(context.Background(), builder.Sql(), catId)
	defer rows.Close()
	if err == nil && rows.Next() {
		return pgx.RowToAddrOfStructByName[usercase.Category](rows)
	}
	return nil, err
}

func (c *CategoryRepo) List() ([]*usercase.CategoryVo, error) {
	builder := sqlbuild.NewSelectBuilder("t_blog_category as bc").
		Select("bc.category_id", "bc.category_name", "bc.parent_id", "bc.cover_url", "bc.is_top", "bc.is_hot", "bc.view_num").
		LeftJoin("t_blog_article as ba").On("bc.category_id").EqRaw("ANY(ba.category_ids)").And("ba.status").EqRaw("0").And("ba.delete_at").EqRaw("0").BuildAsSelect().
		Select("count(ba.*) as article_num").
		Where("bc.status").EqRaw("0").And("bc.delete_at").EqRaw("0").BuildAsSelect().
		GroupBy("bc.category_id").
		OrderBy("bc.is_top desc", "bc.sort")
	rows, err := c.db.Query(context.Background(), builder.Sql())
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (*usercase.CategoryVo, error) {
		return pgx.RowToAddrOfStructByNameLax[usercase.CategoryVo](row)
	})
}

func (c *CategoryRepo) ManageList() ([]*usercase.CategoryVo, error) {
	builder := sqlbuild.NewSelectBuilder("t_blog_category as bc").
		Select("bc.*").
		// 后台查看文章数量忽略status字段
		LeftJoin("t_blog_article as ba").On("bc.category_id").EqRaw("ANY(ba.category_ids)").And("ba.delete_at").EqRaw("0").BuildAsSelect().
		Select("count(ba.*) as article_num").
		Where("bc.delete_at").EqRaw("0").BuildAsSelect().
		GroupBy("bc.category_id").
		OrderBy("bc.is_top desc", "bc.sort")
	rows, err := c.db.Query(context.Background(), builder.Sql())
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (*usercase.CategoryVo, error) {
		return pgx.RowToAddrOfStructByNameLax[usercase.CategoryVo](row)
	})
}

func (c *CategoryRepo) ListByIds(ids []uint) ([]usercase.Category, error) {
	if len(ids) == 0 {
		return make([]usercase.Category, 0), nil
	}
	builder := sqlbuild.NewSelectBuilder("t_blog_category").
		Select("category_id", "category_name").
		Where("category_id").In(sqlbuild.SliceToAnySlice[uint](ids)...).
		And("status").EqRaw("0").
		And("delete_at").EqRaw("0").BuildAsSelect()
	rows, err := c.db.Query(context.Background(), builder.Sql(), builder.Args()...)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (usercase.Category, error) {
		return pgx.RowToStructByNameLax[usercase.Category](row)
	})
}

func (c *CategoryRepo) CountByName(name string, catId uint) (uint8, error) {
	builder := sqlbuild.NewSelectBuilder("t_blog_category").
		Select("count(*)").
		Where("category_id").Ne(catId).
		And("category_name").Eq(name).
		And("delete_at").EqRaw("0").BuildAsSelect()
	row := c.db.QueryRow(context.Background(), builder.Sql(), builder.Args()...)
	var total uint8
	err := row.Scan(&total)
	return total, err
}

func (c *CategoryRepo) DeleteById(catId int) error {
	builder := sqlbuild.NewUpdateBuilder("t_blog_category").
		Set("delete_at", time.Now().UnixMilli()).
		Where("category_id").Eq(catId).BuildAsUpdate()
	result, err := c.db.Exec(context.Background(), builder.Sql(), builder.Args()...)
	if err == nil {
		slog.Info(fmt.Sprintf("分类删除完成，row:%d,id:%d", result.RowsAffected(), catId))
	}
	return err
}

func (c *CategoryRepo) BatchDelete(ids []int) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	builder := sqlbuild.NewUpdateBuilder("t_blog_category").
		Set("delete_at", time.Now().UnixMilli()).
		Where("category_id").In(sqlbuild.SliceToAnySlice[int](ids)...).BuildAsUpdate()
	result, err := c.db.Exec(context.Background(), builder.Sql(), builder.Args()...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), err
}
