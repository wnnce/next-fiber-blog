package data

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-fiber-ent-web-layout/internal/tools"
	"go-fiber-ent-web-layout/internal/usercase"
	"log/slog"
	"strings"
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
	sql := "insert into t_blog_notice (title, message, level, notice_type, sort, status) values ($1, $2, $3, $4, $5, $6) returning notice_id"
	row := self.db.QueryRow(context.Background(), sql, notice.Title, notice.Message, notice.Level, notice.NoticeType, notice.Sort, notice.Status)
	var noticeId uint64
	err := row.Scan(&noticeId)
	if err == nil {
		notice.NoticeId = noticeId
		slog.Info("保存系统通知完成", "noticeId", noticeId)
	}
	return err
}

func (self *NoticeRepo) Update(notice *usercase.Notice) error {
	sql := "update t_blog_notice set title = $1, message = $2, level = $3, notice_type = $4, sort = $5, status = $6 where notice_id = $7"
	result, err := self.db.Exec(context.Background(), sql, notice.Title, notice.Message, notice.Level, notice.NoticeType, notice.Sort, notice.Status, notice.NoticeId)
	if err == nil {
		slog.Info("更新系统通知完成", "row", result.RowsAffected(), "noticeId", notice.NoticeId)
	}
	return err
}

func (self *NoticeRepo) ListByType(noticeType int) ([]usercase.Notice, error) {
	sql := "select title, message, level, notice_type from t_blog_notice where notice_type = $1 and status = 0 and delete_at = 0 order by sort"
	rows, err := self.db.Query(context.Background(), sql, noticeType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (usercase.Notice, error) {
		return pgx.RowToStructByNameLax[usercase.Notice](row)
	})
}

func (self *NoticeRepo) ManagePage(query *usercase.NoticeQueryForm) ([]*usercase.Notice, int64, error) {
	var builder strings.Builder
	builder.WriteString(" where delete_at = 0")
	args := make([]any, 0)
	if query.Title != "" {
		args = append(args, "%"+query.Title+"%")
		builder.WriteString(fmt.Sprintf(" and title like $%d", len(args)))
	}
	if query.Level != nil {
		args = append(args, *query.Level)
		builder.WriteString(fmt.Sprintf(" and level = $%d", len(args)))
	}
	if query.NoticeType != nil {
		args = append(args, *query.NoticeType)
		builder.WriteString(fmt.Sprintf(" and notice_type = $%d", len(args)))
	}
	row := self.db.QueryRow(context.Background(), "select count(*) from t_blog_notice "+builder.String(), args...)
	var total int64
	if err := row.Scan(&total); err != nil {
		return nil, 0, err
	}
	records := make([]*usercase.Notice, 0)
	if total == 0 {
		return records, 0, nil

	}
	offset := tools.ComputeOffset(total, query.Page, query.Size, false)
	builder.WriteString(fmt.Sprintf(" order by sort, create_time desc limit $%d offset $%d", len(args)+1, len(args)+2))
	args = append(args, query.Size, offset)
	rows, err := self.db.Query(context.Background(), "select * from t_blog_notice "+builder.String(), args...)
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
	sql := "select notice_type from t_blog_notice where notice_id = $1"
	row := self.db.QueryRow(context.Background(), sql, noticeId)
	noticeType := -1
	_ = row.Scan(&noticeType)
	return noticeType
}

func (self *NoticeRepo) DeleteById(id int64) error {
	sql := "update t_blog_notice set delete_at = $1 where notice_id = $2"
	result, err := self.db.Exec(context.Background(), sql, time.Now().UnixMilli(), id)
	if err == nil {
		slog.Info("删除系统通知完成", "row", result.RowsAffected(), "noticeId", id)
	}
	return err
}
