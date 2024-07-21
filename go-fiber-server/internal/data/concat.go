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

type ConcatRepo struct {
	db *pgxpool.Pool
}

func NewConcatRepo(data *Data) usercase.IConcatRepo {
	return &ConcatRepo{
		db: data.Db,
	}
}

func (c *ConcatRepo) Save(concat *usercase.Concat) error {
	var insertId uint
	err := c.db.QueryRow(context.Background(), "insert into t_blog_concat (name, logo_url, target_url, is_main, sort, status) values ($1, $2, $3, $4, $5, $6) returning concat_id",
		concat.Name, concat.LogoUrl, concat.TargetUrl, concat.IsMain, concat.Sort, concat.Status).Scan(&insertId)
	if err == nil {
		slog.Info("联系方式添加完成，id：" + strconv.Itoa(int(insertId)))
	}
	return err
}

func (c *ConcatRepo) Update(concat *usercase.Concat) error {
	sql := "update t_blog_concat set update_time = now(), name = $1, logo_url = $2, target_url = $3, is_main = $4, sort = $5, status = $6 where concat_id = $7"
	result, err := c.db.Exec(context.Background(), sql, concat.Name, concat.LogoUrl, concat.TargetUrl, concat.IsMain, concat.Sort, concat.Status, concat.ConcatId)
	if err == nil {
		slog.Info(fmt.Sprintf("联系方式更新完成，id：%d，row：%d"), concat.ConcatId, result.RowsAffected())
	}
	return err
}

func (c *ConcatRepo) UpdateStatus(cid int, status uint) error {
	result, err := c.db.Exec(context.Background(), "update t_blog_concat set update_time = now(), status = $1 where concat_id = $2", status, cid)
	if err == nil {
		slog.Info(fmt.Sprintf("更新联系方式状态完成，row:%d,id:%d,status:%d", result.RowsAffected(), cid, status))
	}
	return err
}

func (c *ConcatRepo) List() ([]*usercase.Concat, error) {
	rows, err := c.db.Query(context.Background(), "select concat_id, name, logo_url, target_url, is_main, sort, status from t_blog_concat where delete_at = '0' and status = 0 order by sort, create_time desc")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (*usercase.Concat, error) {
		return pgx.RowToAddrOfStructByNameLax[usercase.Concat](row)
	})
}

func (c *ConcatRepo) ManageList(query *usercase.ConcatQueryForm) ([]*usercase.Concat, error) {
	var builder strings.Builder
	builder.WriteString("select * from t_blog_concat where delete_at = '0' ")
	args := make([]any, 0)
	if query.Name != "" {
		args = append(args, "%"+query.Name+"%")
		builder.WriteString(fmt.Sprintf("and name like $%d ", len(args)))
	}
	timeQueryConditionBuilder(query.CreateTimeBegin, query.CreateTimeEnd, &builder, &args)
	builder.WriteString("order by sort, create_time desc")
	rows, err := c.db.Query(context.Background(), builder.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (*usercase.Concat, error) {
		return pgx.RowToAddrOfStructByName[usercase.Concat](rows)
	})
}

func (c *ConcatRepo) CountByName(name string, cid uint) (uint8, error) {
	row := c.db.QueryRow(context.Background(), "select count(concat_id) from t_blog_concat where name = $1 and concat_id != $2", name, cid)
	var total uint8
	err := row.Scan(&total)
	return total, err
}

func (c *ConcatRepo) DeleteById(cid int) error {
	result, err := c.db.Exec(context.Background(), "update t_blog_concat set delete_at = $1 where concat_id = $2", time.Now().UnixMilli(), cid)
	if err == nil {
		slog.Info(fmt.Sprintf("删除联系方式完成，concatId:%d,row:%d", cid, result.RowsAffected()))
	}
	return err
}
