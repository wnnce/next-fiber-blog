package data

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"go-fiber-ent-web-layout/internal/tools"
	"go-fiber-ent-web-layout/internal/usercase"
	"log/slog"
	"strings"
	"time"
)

type ConcatRepo struct {
	db *sqlx.DB
}

func NewConcatRepo(data *Data) usercase.IConcatRepo {
	return &ConcatRepo{
		db: data.Db,
	}
}

func (c *ConcatRepo) Save(concat *usercase.Concat) error {
	result, err := c.db.Exec("insert into t_blog_concat (name, logo_url, target_url, is_main, sort, status) values ($1, $2, $3, $4, $5, $6)", concat.Name, concat.LogoUrl, concat.TargetUrl, concat.IsMain, concat.Sort, concat.Status)
	if err == nil {
		id, _ := result.LastInsertId()
		slog.Info(fmt.Sprintf("联系方式添加完成，id：%d", id))
	}
	return err
}

func (c *ConcatRepo) Update(concat *usercase.Concat) error {
	result, err := c.db.NamedExec("update t_blog_concat set update_time = now(), name = :name, logo_url = :logoUrl, target_url = :targetUrl, is_main = :isMain, sort = :sort, status = :status where concat_id = :id", map[string]any{
		"name":      concat.Name,
		"logoUrl":   concat.LogoUrl,
		"targetUrl": concat.TargetUrl,
		"isMain":    concat.IsMain,
		"sort":      concat.Sort,
		"status":    concat.Status,
		"id":        concat.ConcatId,
	})
	if err == nil {
		row, _ := result.RowsAffected()
		slog.Info(fmt.Sprintf("联系方式更新完成，row:%d, id:%d", row, concat.ConcatId))
	}
	return err
}

func (c *ConcatRepo) UpdateStatus(cid int, status uint) error {
	result, err := c.db.Exec("update t_blog_concat set update_time = now(), status = $1 where concat_id = $2", status, cid)
	if err == nil {
		row, _ := result.RowsAffected()
		slog.Info(fmt.Sprintf("更新联系方式状态完成，row:%d, id:%d, status:%d", row, cid, status))
	}
	return err
}

func (c *ConcatRepo) List() ([]*usercase.Concat, error) {
	concats := make([]*usercase.Concat, 0)
	err := c.db.Select(&concats, "select concat_id, name, logo_url, target_url, is_main from t_blog_concat where delete_at = '0' and status = 0 order by sort")
	return concats, err
}

func (c *ConcatRepo) ManageList(query *usercase.ConcatQueryForm) ([]*usercase.Concat, error) {
	var builder strings.Builder
	builder.WriteString("select * from t_blog_concat where delete_at = '0'")
	args := make(map[string]any)
	if query.Name != "" {
		builder.WriteString(" and name like :name")
		args["name"] = "%" + query.Name + "%"
	}
	if query.CreateTimeBegin != nil {
		builder.WriteString(" and create_time >= :begin")
		args["begin"] = query.CreateTimeBegin.Format("2006-04-02")
	}
	if query.CreateTimeEnd != nil {
		builder.WriteString("and create_time <= :end")
		args["end"] = query.CreateTimeEnd.Format("2006-04-02")
	}
	builder.WriteString(" order by sort")
	rows, err := c.db.NamedQuery(builder.String(), args)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
	concats := make([]*usercase.Concat, 0)
	concats = tools.SqlxRowsScan(rows, concats)
	return concats, err
}

func (c *ConcatRepo) CountByName(name string, cid uint) (uint8, error) {
	row := c.db.QueryRow("select count(concat_id) from t_blog_concat where name = $1 and concat_id != $2", name, cid)
	var total uint8
	err := row.Scan(&total)
	return total, err
}

func (c *ConcatRepo) DeleteById(cid int) error {
	result, err := c.db.Exec("update t_blog_concat set delete_at = $1 where concat_id = $2", time.Now().UnixMilli(), cid)
	if err == nil {
		row, _ := result.RowsAffected()
		slog.Info(fmt.Sprintf("删除联系方式完成，row:%d, cid:%d", row, cid))
	}
	return err
}
