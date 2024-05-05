package data

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"go-fiber-ent-web-layout/internal/usercase"
	"log/slog"
	"strings"
	"time"
)

type CategoryRepo struct {
	db *sqlx.DB // Sqlx数据库连接
}

func NewCategoryRepo(data *Data) usercase.ICategoryRepo {
	return &CategoryRepo{
		db: data.Db,
	}
}

func (c *CategoryRepo) Save(cat *usercase.Category) error {
	result, err := c.db.Exec("insert into t_blog_category (category_name, description, cover_url, parent_id, is_top, is_hot, sort, status) values ($1, $2, $3, $4, $5, $6, $7, $8)",
		cat.CategoryName, cat.Description, cat.CoverUrl, cat.ParentId, cat.IsTop, cat.IsHot, cat.Sort, cat.Status)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	slog.Info(fmt.Sprintf("分类新增成功，分类ID:%d", id))
	return nil
}

func (c *CategoryRepo) Update(cat *usercase.Category) error {
	result, err := c.db.Exec(`update t_blog_category
					set update_time = now(),
						category_name = $1,description = $2,cover_url = $3,parent_id = $4,
						is_hot = $5,is_top  = $6,sort  = $7,status  = $8
					where category_id = $9`, cat.CategoryName, cat.Description, cat.CoverUrl, cat.ParentId,
		cat.IsHot, cat.IsTop, cat.Sort, cat.Status, cat.CategoryId)
	if err != nil {
		return err
	}
	row, _ := result.RowsAffected()
	slog.Info(fmt.Sprintf("分类更新成功，tagId:%d,row:%d", cat.CategoryId, row))
	return nil
}

func (c *CategoryRepo) UpdateStatus(catId int, status uint8) error {
	result, err := c.db.Exec("update t_blog_category set status = $s where category_id = $2", status, catId)
	if err != nil {
		return err
	}
	row, _ := result.RowsAffected()
	slog.Info(fmt.Sprintf("更新分类状态完成，tagId：%d，status：%d，row：%d", catId, status, row))
	return nil
}

func (c *CategoryRepo) UpdateViewNum(catId uint, addNum uint) error {
	result, err := c.db.Exec("update t_blog_category set view_num = view_num + $1 where category_id = $2", addNum, catId)
	if err != nil {
		return err
	}
	row, _ := result.RowsAffected()
	slog.Info(fmt.Sprintf("更新分类查看次数完成，tagId:%d,addNum:%d,row:%d", catId, addNum, row))
	return nil
}

func (c *CategoryRepo) SelectById(catId int) (*usercase.Category, error) {
	category := &usercase.Category{}
	if err := c.db.Get(category, "select * from t_blog_category where category_id = $1 and delete_at = '0' and status = 0", catId); err != nil {
		return nil, err
	}
	return category, nil
}

func (c *CategoryRepo) List() []*usercase.Category {
	categorys := make([]*usercase.Category, 0)
	if err := c.db.Select(&categorys, `select category_id, category_name, parent_id, cover_url, is_top, is_hot, view_num
									from t_blog_category
									where delete_at = '0'
									  and status = 0
									order by is_top desc, sort`); err != nil {
		slog.Error(fmt.Sprintf("获取分类列表失败，错误信息：%s", err))
	}
	return categorys
}

func (c *CategoryRepo) ManageList() ([]*usercase.Category, error) {
	categorys := make([]*usercase.Category, 0)
	err := c.db.Select(&categorys, `select *
									from t_blog_category
									where delete_at = '0'
									order by is_top desc, sort`)
	return categorys, err
}

func (c *CategoryRepo) ListByIds(ids []uint) ([]*usercase.Category, error) {
	categorys := make([]*usercase.Category, 0)
	if len(ids) == 0 {
		return categorys, nil
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
	err := c.db.Select(&categorys, builder.String())
	return categorys, err
}

func (c *CategoryRepo) CountByName(name string, catId uint) (uint8, error) {
	row := c.db.QueryRow("select count(category_id) from t_blog_category where delete_at = '0' and category_name = $1 and category_id != $2", name, catId)
	var total uint8
	err := row.Scan(&total)
	return total, err
}

func (c *CategoryRepo) DeleteById(catId int) error {
	result, err := c.db.Exec("update t_blog_category set delete_at = $1 where category_id = $2", time.Now().UnixMilli(), catId)
	if err != nil {
		return err
	}
	row, _ := result.RowsAffected()
	slog.Info(fmt.Sprintf("分类删除完成，catId:%d, row:%d", catId, row))
	return nil
}

func (c *CategoryRepo) BatchDelete(ids []int) (int64, error) {
	var builder strings.Builder
	builder.WriteString("update t_blog_category set delete_at = $1 where category_id in (")
	for i, id := range ids {
		if i > 0 {
			builder.WriteByte(',')
		}
		builder.WriteRune(rune(id))
	}
	builder.WriteByte(')')
	result, err := c.db.Exec(builder.String())
	if err != nil {
		return 0, err
	}
	row, _ := result.RowsAffected()
	return row, nil
}
