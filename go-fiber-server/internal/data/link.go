package data

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-fiber-ent-web-layout/internal/tools"
	"go-fiber-ent-web-layout/internal/usercase"
	"log/slog"
	"strconv"
	"strings"
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

func (l *LinkRepo) Save(link *usercase.Link) error {
	row := l.db.QueryRow(context.Background(), "insert into t_blog_link (name, summary, cover_url, target_url, sort, status) values ($1, $2, $3, $4, $5, $6) returning link_id",
		link.Name, link.Summary, link.CoverUrl, link.TargetUrl, link.Sort, link.Status)
	var linkId uint64
	err := row.Scan(&linkId)
	if err == nil {
		link.LinkId = linkId
		slog.Info(fmt.Sprintf("友情链接添加完成，id：%d", linkId))
	}
	return err
}

func (l *LinkRepo) Update(link *usercase.Link) error {
	result, err := l.db.Exec(context.Background(), "update t_blog_link set update_time = now(), name = $1, summary = $2, cover_url = $3, target_url = $4, sort = $5, status = $6 where link_id = $7",
		link.Name, link.Summary, link.CoverUrl, link.TargetUrl, link.Sort, link.Status, link.LinkId)
	if err == nil {
		slog.Info(fmt.Sprintf("更新友情链接完成，row:%d,id:%d", result.RowsAffected(), link.LinkId))
	}
	return err
}

func (l *LinkRepo) UpdateStatus(linkId int64, status uint8) error {
	result, err := l.db.Exec(context.Background(), "update t_blog_link set update_time = now(), status = $1 where link_id = $2", status, linkId)
	if err == nil {
		slog.Info(fmt.Sprintf("联系方式状态更新完成，row:%d,id:%d,status:%d", result.RowsAffected(), linkId, status))
	}
	return err
}

func (l *LinkRepo) Page(query *usercase.PageQueryForm) ([]*usercase.Link, int64, error) {
	condition := "where delete_at = '0' and status = 0"
	total, err := l.conditionTotal(condition)
	if err != nil {
		return nil, 0, err
	}
	links := make([]*usercase.Link, 0)
	if total == 0 {
		return links, 0, nil
	}
	offset := tools.ComputeOffset(total, query.Page, query.Size, false)
	rows, err := l.db.Query(context.Background(), "select link_id, name, summary, cover_url, target_url, click_num, create_time from t_blog_link "+condition+
		" order by sort, create_time desc limit $1 offset $2", query.Size, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	links, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (*usercase.Link, error) {
		return pgx.RowToAddrOfStructByNameLax[usercase.Link](row)
	})
	return links, total, err
}

func (l *LinkRepo) ManagePage(query *usercase.LinkQueryForm) ([]*usercase.Link, int64, error) {
	var condition strings.Builder
	condition.WriteString("where delete_at = '0'")
	if query.Name != "" {
		condition.WriteString(fmt.Sprintf(" and name like '%s'", "%"+query.Name+"%"))
	}
	if query.CreateTimeBegin != nil {
		condition.WriteString(fmt.Sprintf("create_time >= '%s'", query.CreateTimeBegin.Format("2006-04-02")))
	}
	if query.CreateTimeEnd != nil {
		condition.WriteString(fmt.Sprintf("create_time <= '%s'", query.CreateTimeEnd.Format("2006-01-02")))
	}
	total, err := l.conditionTotal(condition.String())
	if err != nil {
		return nil, 0, err
	}
	links := make([]*usercase.Link, 0)
	if total == 0 {
		return links, 0, nil
	}
	offset := tools.ComputeOffset(total, query.Page, query.Size, false)
	rows, err := l.db.Query(context.Background(), "select * from t_blog_link "+condition.String()+" order by create_time desc limit $1 offset $2", query.Size, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	links, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (*usercase.Link, error) {
		return pgx.RowToAddrOfStructByName[usercase.Link](row)
	})
	return links, total, err
}

func (l *LinkRepo) DeleteById(linkId int64) error {
	deleteAt := strconv.FormatInt(time.Now().UnixMilli(), 10)
	result, err := l.db.Exec(context.Background(), "update t_blog_link set delete_at = $1 where link_id = $2", deleteAt, linkId)
	if err == nil {
		slog.Info(fmt.Sprintf("友情链接删除完成，row：%d，linkId：%d", result.RowsAffected(), linkId))
	}
	return err
}

func (l *LinkRepo) BatchDelete(linkIds []int64) (int64, error) {
	if len(linkIds) == 0 {
		return 0, nil
	}
	var builder strings.Builder
	builder.WriteString("update t_blog_link set delete_at = $1 where link_id in (")
	for i, v := range linkIds {
		if i > 0 {
			builder.WriteByte(',')
		}
		builder.WriteRune(rune(v))
	}
	builder.WriteByte(')')
	deleteAt := strconv.FormatInt(time.Now().UnixMilli(), 10)
	result, err := l.db.Exec(context.Background(), builder.String(), deleteAt)
	return result.RowsAffected(), err
}

func (l *LinkRepo) conditionTotal(condition string) (int64, error) {
	row := l.db.QueryRow(context.Background(), "select count(link_id) from t_blog_link "+condition)
	var total int64
	err := row.Scan(&total)
	return total, err
}
