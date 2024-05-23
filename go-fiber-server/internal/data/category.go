package data

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-fiber-ent-web-layout/internal/usercase"
	"log/slog"
	"strconv"
	"strings"
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
	row := c.db.QueryRow(context.Background(), "insert into t_blog_category (category_name, description, cover_url, parent_id, is_top, is_hot, sort, status) values ($1, $2, $3, $4, $5, $6, $7, $8) returning category_id",
		cat.CategoryName, cat.Description, cat.CoverUrl, cat.ParentId, cat.IsTop, cat.IsHot, cat.Sort, cat.Status)
	var categoryId uint
	err := row.Scan(&categoryId)
	if err == nil {
		cat.CategoryId = categoryId
		slog.Info("分类保存完成，id：" + strconv.Itoa(int(categoryId)))
	}
	return err
}

func (c *CategoryRepo) Update(cat *usercase.Category) error {
	result, err := c.db.Exec(context.Background(), `update t_blog_category
					set update_time = now(),
						category_name = $1,description = $2,cover_url = $3,parent_id = $4,
						is_hot = $5,is_top  = $6,sort  = $7,status  = $8
					where category_id = $9`, cat.CategoryName, cat.Description, cat.CoverUrl, cat.ParentId,
		cat.IsHot, cat.IsTop, cat.Sort, cat.Status, cat.CategoryId)
	if err == nil {
		slog.Info(fmt.Sprintf("分类更新完成，row:%d,id:%d", result.RowsAffected(), cat.CategoryId))
	}
	return err
}

func (c *CategoryRepo) UpdateStatus(catId int, status uint8) error {
	result, err := c.db.Exec(context.Background(), "update t_blog_category set status = $s where category_id = $2", status, catId)
	if err == nil {
		slog.Info(fmt.Sprintf("分类状态更新完成，row:%d,id:%d,status:%d", result.RowsAffected(), catId, status))
	}
	return err
}

func (c *CategoryRepo) UpdateViewNum(catId uint, addNum uint) error {
	result, err := c.db.Exec(context.Background(), "update t_blog_category set view_num = view_num + $1 where category_id = $2", addNum, catId)
	if err == nil {
		slog.Info(fmt.Sprintf("更新分类查看次数完成，tagId:%d,addNum:%d,row:%d", catId, addNum, result.RowsAffected()))
	}
	return err
}

func (c *CategoryRepo) SelectById(catId int) (*usercase.Category, error) {
	rows, err := c.db.Query(context.Background(), "select * from t_blog_category where category_id = $1 and delete_at = '0' and status = 0", catId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		category, err := pgx.RowToStructByNameLax[usercase.Category](rows)
		return &category, err
	}
	return nil, nil
}

func (c *CategoryRepo) List() ([]*usercase.Category, error) {
	rows, err := c.db.Query(context.Background(), `select category_id, category_name, parent_id, cover_url, is_top, is_hot, view_num
									from t_blog_category
									where delete_at = '0'
									  and status = 0
									order by is_top desc, sort`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (*usercase.Category, error) {
		return pgx.RowToAddrOfStructByNameLax[usercase.Category](row)
	})
}

func (c *CategoryRepo) ManageList() ([]*usercase.Category, error) {
	rows, err := c.db.Query(context.Background(), `select *
									from t_blog_category
									where delete_at = '0'
									order by is_top desc, sort`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (*usercase.Category, error) {
		return pgx.RowToAddrOfStructByNameLax[usercase.Category](row)
	})
}

func (c *CategoryRepo) ListByIds(ids []uint) ([]usercase.Category, error) {
	if len(ids) == 0 {
		return make([]usercase.Category, 0), nil
	}
	var builder strings.Builder
	builder.WriteString("select category_id, category_name from t_blog_category where delete_at = '0' and status = 0 and category_id in (")
	for i, id := range ids {
		if i > 0 {
			builder.WriteByte(',')
		}
		builder.WriteRune(rune(id))
	}
	builder.WriteByte(')')
	rows, err := c.db.Query(context.Background(), builder.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (usercase.Category, error) {
		return pgx.RowToStructByNameLax[usercase.Category](row)
	})
}

func (c *CategoryRepo) CountByName(name string, catId uint) (uint8, error) {
	row := c.db.QueryRow(context.Background(), "select count(category_id) from t_blog_category where delete_at = '0' and category_name = $1 and category_id != $2", name, catId)
	var total uint8
	err := row.Scan(&total)
	return total, err
}

func (c *CategoryRepo) DeleteById(catId int) error {
	result, err := c.db.Exec(context.Background(), "update t_blog_category set delete_at = $1 where category_id = $2",
		time.Now().UnixMilli(), catId)
	if err == nil {
		slog.Info(fmt.Sprintf("分类删除完成，row:%d,id:%d", result.RowsAffected(), catId))
	}
	return err
}

func (c *CategoryRepo) BatchDelete(ids []int) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	var builder strings.Builder
	builder.WriteString("update t_blog_category set delete_at = $1 where category_id in (")
	for i, id := range ids {
		if i > 0 {
			builder.WriteByte(',')
		}
		builder.WriteRune(rune(id))
	}
	builder.WriteByte(')')
	result, err := c.db.Exec(context.Background(), builder.String(), time.Now().UnixMilli())
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), err
}
