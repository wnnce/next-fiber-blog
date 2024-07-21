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

type TagRepo struct {
	db *pgxpool.Pool
}

func NewTagRepo(data *Data) usercase.ITagRepo {
	return &TagRepo{
		db: data.Db,
	}
}

func (self *TagRepo) Save(form *usercase.TagForm) error {
	row := self.db.QueryRow(context.Background(), "Insert Into t_blog_tag (tag_name, cover_url, color, sort, status) values ($1,$2,$3,$4,$5) returning tag_id",
		form.TagName, form.CoverUrl, form.Color, *form.Sort, *form.Status)
	var insertId int
	err := row.Scan(&insertId)
	if err == nil {
		slog.Info("新增标签完成，id：" + strconv.Itoa(insertId))
	}
	return err
}

func (self *TagRepo) Update(form *usercase.TagForm) error {
	result, err := self.db.Exec(context.Background(), "update t_blog_tag set update_time = now(), tag_name = $1, cover_url = $2, color = $3, sort = $4, status = $5 where tag_id = $6",
		form.TagName, form.CoverUrl, form.Color, *form.Sort, *form.Status, form.TagId)
	if err == nil {
		slog.Info(fmt.Sprintf("标签更新完成，row:%d,id:%d", result.RowsAffected(), form.TagId))
	}
	return err
}

func (self *TagRepo) UpdateStatus(tagId int, status uint8) error {
	result, err := self.db.Exec(context.Background(), "update t_blog_tag set update_time = now(), status = $1 where tag_id = $2", status, tagId)
	if err == nil {
		slog.Info(fmt.Sprintf("标签状态更新完成，row:%d,id:%d,status:%d", result.RowsAffected(), tagId, status))
	}
	return err
}

func (self *TagRepo) UpdateViewNum(tagId int, addNum int) error {
	result, err := self.db.Exec(context.Background(), "update t_blog_tag set update_time = now(), view_num = view_num + $1 where tag_id = $2", addNum, tagId)
	if err == nil {
		slog.Info(fmt.Sprintf("更新标签查看次数完成，row:%d,id:%d,addnum:%d", result.RowsAffected(), tagId, addNum))
	}
	return err
}

func (self *TagRepo) SelectById(id int) (*usercase.Tag, error) {
	rows, err := self.db.Query(context.Background(), "select * from t_blog_tag where tag_id = $1 and delete_at = '0' and status = 0", id)
	if err == nil && rows.Next() {
		defer rows.Close()
		return pgx.RowToAddrOfStructByName[usercase.Tag](rows)
	}
	return nil, err
}

func (self *TagRepo) Page(query *usercase.TagQueryForm) ([]*usercase.Tag, int64, error) {
	var condition strings.Builder
	condition.WriteString("where delete_at = '0'")
	args := make([]any, 0)
	if query.TagName != "" {
		args = append(args, "%"+query.TagName+"%")
		condition.WriteString(fmt.Sprintf(" and tag_name like $%d", len(args)))
	}
	timeQueryConditionBuilder(query.CreateTimeBegin, query.CreateTimeEnd, &condition, &args)
	total, err := self.conditionTotal(condition.String(), args...)
	if err != nil {
		return nil, 0, err
	}
	tags := make([]*usercase.Tag, 0)
	if total == 0 {
		return tags, total, nil
	}
	offset := tools.ComputeOffset(total, query.Page, query.Size, false)
	condition.WriteString(fmt.Sprintf(" order by sort asc, create_time desc limit $%d offset $%d", len(args)+1, len(args)+2))
	args = append(args, query.Size, offset)
	rows, err := self.db.Query(context.Background(), "select * from t_blog_tag "+condition.String(), args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	tags, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (*usercase.Tag, error) {
		return pgx.RowToAddrOfStructByName[usercase.Tag](row)
	})
	return tags, total, err
}

func (self *TagRepo) List() ([]*usercase.Tag, error) {
	rows, err := self.db.Query(context.Background(), "select tag_id, tag_name, cover_url, view_num, color, create_time from t_blog_tag where delete_at = '0' and status = 0 order by sort asc, create_time desc")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (*usercase.Tag, error) {
		return pgx.RowToAddrOfStructByNameLax[usercase.Tag](row)
	})
}

func (self *TagRepo) ListByIds(ids []uint) ([]*usercase.Tag, error) {
	if len(ids) == 0 {
		return make([]*usercase.Tag, 0), nil
	}
	var builder strings.Builder
	builder.WriteString("select tag_id, tag_name, color from t_blog_tag where delete_at = '0' and status = 0 and tag_id in (")
	for i, id := range ids {
		if i > 0 {
			builder.WriteByte(',')
		}
		builder.WriteRune(rune(id))
	}
	builder.WriteByte(')')
	rows, err := self.db.Query(context.Background(), builder.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (*usercase.Tag, error) {
		return pgx.RowToAddrOfStructByNameLax[usercase.Tag](row)
	})
}

func (self *TagRepo) CountByTagName(name string, tagId uint) (uint8, error) {
	row := self.db.QueryRow(context.Background(), "select count(tag_id) from t_blog_tag where tag_name = $1 and delete_at = '0' and tag_id != $2", name, tagId)
	var total uint8
	err := row.Scan(&total)
	return total, err
}

func (self *TagRepo) DeleteById(id int) error {
	result, err := self.db.Exec(context.Background(), "update t_blog_tag set delete_at = $1 where tag_id = $2",
		time.Now().UnixMilli(), id)
	if err == nil {
		slog.Info(fmt.Sprintf("删除标签完成，row：%d,id:%d", result.RowsAffected(), id))
	}
	return err
}

func (self *TagRepo) DeleteByIds(ids []int) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	var builder strings.Builder
	builder.WriteString("update t_blog_tag set delete_at = $1 where tag_id in (")
	for i, id := range ids {
		if i > 0 {
			builder.WriteByte(',')
		}
		builder.WriteString(strconv.Itoa(id))
	}
	builder.WriteByte(')')
	result, err := self.db.Exec(context.Background(), builder.String(), time.Now().UnixMilli())
	return result.RowsAffected(), err
}

func (self *TagRepo) conditionTotal(condition string, args ...any) (int64, error) {
	row := self.db.QueryRow(context.Background(), "select count(*) from t_blog_tag "+condition, args...)
	var total int64
	err := row.Scan(&total)
	return total, err
}
