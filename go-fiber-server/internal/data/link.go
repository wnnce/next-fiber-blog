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
	"time"
)

type LinkRepo struct {
	db *pgxpool.Pool
}

func NewLinkRepo(data *Data) usercase.ILinkRepo {
	return &LinkRepo{
		db: data.Db,
	}
}

func (self *LinkRepo) Save(link *usercase.Link) error {
	builder := sqlbuild.NewInsertBuilder("t_blog_link").
		Fields("name", "summary", "cover_url", "target_url", "sort", "status").
		Values(link.Name, link.Summary, link.CoverUrl, link.TargetUrl, link.Sort, link.Status).
		Returning("link_id")
	row := self.db.QueryRow(context.Background(), builder.Sql(), builder.Args()...)
	var linkId uint64
	err := row.Scan(&linkId)
	if err == nil {
		link.LinkId = linkId
		slog.Info(fmt.Sprintf("友情链接添加完成，id：%d", linkId))
	}
	return err
}

func (self *LinkRepo) Update(link *usercase.Link) error {
	builder := sqlbuild.NewUpdateBuilder("t_blog_link").
		SetRaw("update_time", "now()").
		SetByMap(map[string]any{
			"name":       link.Name,
			"summary":    link.Summary,
			"cover_url":  link.CoverUrl,
			"target_url": link.TargetUrl,
			"sort":       link.Sort,
			"status":     link.Status,
		}).Where("link_id").Eq(link.LinkId).BuildAsUpdate()
	result, err := self.db.Exec(context.Background(), builder.Sql(), builder.Args()...)
	if err == nil {
		slog.Info(fmt.Sprintf("更新友情链接完成，row:%d,id:%d", result.RowsAffected(), link.LinkId))
	}
	return err
}

func (self *LinkRepo) UpdateSelective(form *usercase.LinkUpdateForm) error {
	builder := sqlbuild.NewUpdateBuilder("t_blog_link").
		SetRaw("update_time", "now()").
		SetByCondition(form.Status != nil, "status", form.Status).
		Where("link_id").Eq(form.LinkId).BuildAsUpdate()
	result, err := self.db.Exec(context.Background(), builder.Sql(), builder.Args()...)
	if err == nil {
		slog.Info("快捷更新友情链接完成", "row", result.RowsAffected(), "linkId", form.LinkId)
	}
	return err
}

func (self *LinkRepo) List() ([]*usercase.Link, error) {
	builder := sqlbuild.NewSelectBuilder("t_blog_link").
		Select("link_id", "name", "summary", "cover_url", "target_url").
		Where("status").EqRaw("0").
		And("delete_at").EqRaw("0").BuildAsSelect().
		OrderBy("sort", "create_time desc")
	rows, err := self.db.Query(context.Background(), builder.Sql())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (*usercase.Link, error) {
		return pgx.RowToAddrOfStructByNameLax[usercase.Link](row)
	})
}

func (self *LinkRepo) ManagePage(query *usercase.LinkQueryForm) ([]*usercase.Link, int64, error) {
	builder := sqlbuild.NewSelectBuilder("t_blog_link").
		WhereByCondition(query.Name != "", "name").Like("%"+query.Name+"%").
		AndByCondition(query.CreateTimeBegin != "", "create_time").Ge(query.CreateTimeBegin).
		AndByCondition(query.CreateTimeEnd != "", "create_time").Le(query.CreateTimeEnd).
		And("delete_at").EqRaw("0").BuildAsSelect().
		OrderBy("sort", "create_time desc")
	var total int64
	row := self.db.QueryRow(context.Background(), builder.CountSql(), builder.Args()...)
	if err := row.Scan(&total); err != nil {
		return nil, 0, err
	}
	links := make([]*usercase.Link, 0)
	if total == 0 {
		return links, 0, nil
	}
	offset := tools.ComputeOffset(total, query.Page, query.Size, false)
	builder.Limit(int64(query.Size)).Offset(offset)
	rows, err := self.db.Query(context.Background(), builder.Sql(), builder.Args()...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	links, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (*usercase.Link, error) {
		return pgx.RowToAddrOfStructByName[usercase.Link](row)
	})
	return links, total, err
}

func (self *LinkRepo) DeleteById(linkId int64) error {
	builder := sqlbuild.NewUpdateBuilder("t_blog_link").
		Set("delete_at", time.Now().UnixMilli()).
		Where("link_id").Eq(linkId).BuildAsUpdate()
	result, err := self.db.Exec(context.Background(), builder.Sql(), builder.Args()...)
	if err == nil {
		slog.Info(fmt.Sprintf("友情链接删除完成，row：%d，linkId：%d", result.RowsAffected(), linkId))
	}
	return err
}

func (self *LinkRepo) BatchDelete(linkIds []int64) (int64, error) {
	if len(linkIds) == 0 {
		return 0, nil
	}
	builder := sqlbuild.NewUpdateBuilder("t_blog_link").
		Set("delete_at", time.Now().UnixMilli()).
		Where("link_id").In(sqlbuild.SliceToAnySlice[int64](linkIds)).BuildAsUpdate()
	result, err := self.db.Exec(context.Background(), builder.Sql(), builder.Args()...)
	return result.RowsAffected(), err
}
