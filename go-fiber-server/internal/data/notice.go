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

type NoticeRepo struct {
	db *pgxpool.Pool
}

func NewNoticeRepo(data *Data) usercase.INoticeRepo {
	return &NoticeRepo{
		db: data.Db,
	}
}

func (self *NoticeRepo) Save(notice *usercase.Notice) error {
	builder := sqlbuild.NewInsertBuilder("t_blog_notice").
		Fields("title", "message", "level", "notice_type", "sort", "status").
		Values(notice.Title, notice.Message, notice.Level, notice.NoticeType, notice.Sort, notice.Status).
		Returning("notice_id")
	row := self.db.QueryRow(context.Background(), builder.Sql(), builder.Args()...)
	var noticeId uint64
	err := row.Scan(&noticeId)
	if err == nil {
		notice.NoticeId = noticeId
		slog.Info("保存系统通知完成", "noticeId", noticeId)
	}
	return err
}

func (self *NoticeRepo) Update(notice *usercase.Notice) error {
	builder := sqlbuild.NewUpdateBuilder("t_blog_notice").
		SetRaw("update_time", "now()").
		SetByMap(map[string]any{
			"title":       notice.Title,
			"message":     notice.Message,
			"level":       notice.Level,
			"notice_type": notice.NoticeType,
			"sort":        notice.Sort,
			"status":      notice.Status,
		}).Where("notice_id").Eq(notice.NoticeId).BuildAsUpdate()
	result, err := self.db.Exec(context.Background(), builder.Sql(), builder.Args()...)
	if err == nil {
		slog.Info("更新系统通知完成", "row", result.RowsAffected(), "noticeId", notice.NoticeId)
	}
	return err
}

func (self *NoticeRepo) ListByType(noticeType int) ([]usercase.Notice, error) {
	builder := sqlbuild.NewSelectBuilder("t_blog_notice").
		Select("notice_id", "title", "message", "level", "notice_type").
		Where("notice_type").Eq(noticeType).
		And("status").EqRaw("0").
		And("delete_at").EqRaw("0").BuildAsSelect().
		OrderBy("sort")
	rows, err := self.db.Query(context.Background(), builder.Sql(), noticeType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (usercase.Notice, error) {
		return pgx.RowToStructByNameLax[usercase.Notice](row)
	})
}

func (self *NoticeRepo) ManagePage(query *usercase.NoticeQueryForm) ([]*usercase.Notice, int64, error) {
	builder := sqlbuild.NewSelectBuilder("t_blog_notice").
		WhereByCondition(query.Title != "", "title").Like("%"+query.Title+"%").
		AndByCondition(query.Level != nil, "level").Eq(query.Level).
		AndByCondition(query.NoticeType != nil, "notice_type").Eq(query.NoticeType).
		AndByCondition(query.CreateTimeBegin != "", "create_time").Ge(query.CreateTimeBegin).
		AndByCondition(query.CreateTimeEnd != "", "create_time").Le(query.CreateTimeEnd).
		And("delete_at").EqRaw("0").BuildAsSelect().
		OrderBy("sort", "create_time desc")
	row := self.db.QueryRow(context.Background(), builder.CountSql(), builder.Args()...)
	var total int64
	if err := row.Scan(&total); err != nil {
		return nil, 0, err
	}
	records := make([]*usercase.Notice, 0)
	if total == 0 {
		return records, 0, nil

	}
	offset := tools.ComputeOffset(total, query.Page, query.Size, false)
	builder.Limit(int64(query.Size)).Offset(offset)
	rows, err := self.db.Query(context.Background(), builder.Sql(), builder.Args()...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	records, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (*usercase.Notice, error) {
		return pgx.RowToAddrOfStructByName[usercase.Notice](row)
	})
	return records, total, err
}

func (self *NoticeRepo) QueryNoticeTypeById(noticeId int64) int {
	builder := sqlbuild.NewSelectBuilder("t_blog_notice").
		Select("notice_type").
		Where("notice_id").Eq(noticeId).BuildAsSelect()
	row := self.db.QueryRow(context.Background(), builder.Sql(), noticeId)
	noticeType := -1
	_ = row.Scan(&noticeType)
	return noticeType
}

func (self *NoticeRepo) DeleteById(id int64) error {
	builder := sqlbuild.NewUpdateBuilder("t_blog_notice").
		Set("delete_at", time.Now().UnixMilli()).
		Where("notice_id").Eq(id).BuildAsUpdate()
	result, err := self.db.Exec(context.Background(), builder.Sql(), builder.Args()...)
	if err == nil {
		slog.Info("删除系统通知完成", "row", result.RowsAffected(), "noticeId", id)
	}
	return err
}
