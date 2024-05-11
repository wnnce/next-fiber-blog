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

type LinkRepo struct {
	db *sqlx.DB
}

func NewLinkRepo(data *Data) usercase.ILinkRepo {
	return &LinkRepo{
		db: data.Db,
	}
}

func (l *LinkRepo) Save(link *usercase.Link) error {
	result, err := l.db.Exec("insert into t_blog_link (name, summary, cover_url, target_url, sort, status) values ($1, $2, $3, $4, $5, $6)",
		link.Name, link.Summary, link.CoverUrl, link.TargetUrl, link.Sort, link.Status)
	if err == nil {
		id, _ := result.LastInsertId()
		slog.Info(fmt.Sprintf("友情链接新增完成，id:%d", id))
	}
	return err
}

func (l *LinkRepo) Update(link *usercase.Link) error {
	result, err := l.db.NamedExec("update t_blog_link set update_time = now(), name = :name, summary = :summary, cover_url = :cover, target_url = :target, sort = :sort, status = :status where link_id = :id", map[string]any{
		"name":    link.Name,
		"summary": link.Summary,
		"cover":   link.CoverUrl,
		"target":  link.TargetUrl,
		"sort":    link.Sort,
		"status":  link.Status,
		"id":      link.LinkId,
	})
	if err == nil {
		row, _ := result.RowsAffected()
		slog.Info(fmt.Sprintf("联系方式更新完成，row:%d,id:%d", row, link.LinkId))
	}
	return err
}

func (l *LinkRepo) UpdateStatus(linkId int64, status uint8) error {
	result, err := l.db.Exec("update t_blog_link set update_time = now(), status = $1 where link_id = $2", status, linkId)
	if err == nil {
		row, _ := result.RowsAffected()
		slog.Info(fmt.Sprintf("联系方式状态更新完成，row:%d,id:%d,status:%d", row, linkId, status))
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
	err = l.db.Select(&links, "select link_id, name, summary, cover_url, target_url, click_num, create_time, sort from t_blog_link "+condition+
		" order by sort, create_time desc limit $1 offset $2", query.Size, offset)
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
	condition.WriteString(fmt.Sprintf(" order by create_time desc limit %d offset %d", query.Size, offset))
	if err = l.db.Select(&links, "select * from t_blog_link "+condition.String()); err != nil {
		return nil, 0, err
	}
	return links, total, nil
}

func (l *LinkRepo) DeleteById(linkId int64) error {
	result, err := l.db.Exec("update t_blog_link set delete_at = $1 where link_id = $2", time.Now().UnixMilli(), linkId)
	if err == nil {
		row, _ := result.RowsAffected()
		slog.Info(fmt.Sprintf("友情链接删除完成，row：%d，linkId：%d", row, linkId))
	}
	return err
}

func (l *LinkRepo) BatchDelete(linkIds []int64) (int64, error) {
	var builder strings.Builder
	builder.WriteString("update t_blog_link set delete_at = $1 where link_id in (")
	for i, v := range linkIds {
		if i > 0 {
			builder.WriteByte(',')
		}
		builder.WriteRune(rune(v))
	}
	builder.WriteByte(')')
	result, err := l.db.Exec(builder.String())
	if err != nil {
		return 0, err
	}
	row, _ := result.RowsAffected()
	return row, err
}

func (l *LinkRepo) conditionTotal(condition string) (int64, error) {
	row := l.db.QueryRow("select count(link_id) from t_blog_link " + condition)
	var total int64
	if err := row.Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}
