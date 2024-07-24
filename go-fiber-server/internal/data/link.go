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

type LinkRepo struct {
	db *pgxpool.Pool
}

func NewLinkRepo(data *Data) usercase.ILinkRepo {
	return &LinkRepo{
		db: data.Db,
	}
}

func (self *LinkRepo) Save(link *usercase.Link) error {
	row := self.db.QueryRow(context.Background(), "insert into t_blog_link (name, summary, cover_url, target_url, sort, status) values ($1, $2, $3, $4, $5, $6) returning link_id",
		link.Name, link.Summary, link.CoverUrl, link.TargetUrl, link.Sort, link.Status)
	var linkId uint64
	err := row.Scan(&linkId)
	if err == nil {
		link.LinkId = linkId
		slog.Info(fmt.Sprintf("友情链接添加完成，id：%d", linkId))
	}
	return err
}

func (self *LinkRepo) Update(link *usercase.Link) error {
	result, err := self.db.Exec(context.Background(), "update t_blog_link set update_time = now(), name = $1, summary = $2, cover_url = $3, target_url = $4, sort = $5, status = $6 where link_id = $7",
		link.Name, link.Summary, link.CoverUrl, link.TargetUrl, link.Sort, link.Status, link.LinkId)
	if err == nil {
		slog.Info(fmt.Sprintf("更新友情链接完成，row:%d,id:%d", result.RowsAffected(), link.LinkId))
	}
	return err
}

func (self *LinkRepo) UpdateSelective(form *usercase.LinkUpdateForm) error {
	var builder strings.Builder
	builder.WriteString("update t_blog_link set update_time = now() ")
	args := make([]any, 0)
	if form.Status != nil {
		args = append(args, *form.Status)
		builder.WriteString(fmt.Sprintf(", status = $%d", len(args)))
	}
	builder.WriteString(fmt.Sprintf(" where link_id = $%d", len(args)+1))
	args = append(args, form.LinkId)
	result, err := self.db.Exec(context.Background(), builder.String(), args...)
	if err == nil {
		slog.Info("快捷更新友情链接完成", "row", result.RowsAffected(), "linkId", form.LinkId)
	}
	return err
}

func (self *LinkRepo) List() ([]*usercase.Link, error) {
	sql := "select link_id, name, summary, cover_url, target_url from t_blog_link where status = 0 and delete_at = 0 order by sort, create_time desc"
	rows, err := self.db.Query(context.Background(), sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (*usercase.Link, error) {
		return pgx.RowToAddrOfStructByNameLax[usercase.Link](row)
	})
}

func (self *LinkRepo) ManagePage(query *usercase.LinkQueryForm) ([]*usercase.Link, int64, error) {
	var condition strings.Builder
	condition.WriteString("where delete_at = 0")
	args := make([]any, 0)
	if query.Name != "" {
		args = append(args, "%"+query.Name+"%")
		condition.WriteString(fmt.Sprintf(" and name like $%d", len(args)))
	}
	timeQueryConditionBuilder(query.CreateTimeBegin, query.CreateTimeEnd, &condition, &args)
	total, err := self.conditionTotal(condition.String(), args...)
	if err != nil {
		return nil, 0, err
	}
	links := make([]*usercase.Link, 0)
	if total == 0 {
		return links, 0, nil
	}
	offset := tools.ComputeOffset(total, query.Page, query.Size, false)
	condition.WriteString(fmt.Sprintf(" order by sort, create_time desc limit $%d offset $%d", len(args)+1, len(args)+2))
	args = append(args, query.Size, offset)
	rows, err := self.db.Query(context.Background(), "select * from t_blog_link "+condition.String(), args...)
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
	result, err := self.db.Exec(context.Background(), "update t_blog_link set delete_at = $1 where link_id = $2", time.Now().UnixMilli(), linkId)
	if err == nil {
		slog.Info(fmt.Sprintf("友情链接删除完成，row：%d，linkId：%d", result.RowsAffected(), linkId))
	}
	return err
}

func (self *LinkRepo) BatchDelete(linkIds []int64) (int64, error) {
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
	result, err := self.db.Exec(context.Background(), builder.String(), time.Now().UnixMilli())
	return result.RowsAffected(), err
}

func (self *LinkRepo) conditionTotal(condition string, args ...any) (int64, error) {
	row := self.db.QueryRow(context.Background(), "select count(link_id) from t_blog_link "+condition, args...)
	var total int64
	err := row.Scan(&total)
	return total, err
}
