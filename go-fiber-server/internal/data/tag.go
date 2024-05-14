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

type TagRepo struct {
	db *pgxpool.Pool
}

func NewTagRepo(data *Data) usercase.ITagRepo {
	return &TagRepo{
		db: data.Db,
	}
}

func (t *TagRepo) Save(form *usercase.TagForm) error {
	row := t.db.QueryRow(context.Background(), "Insert Into t_blog_tag (tag_name, cover_url, color, sort, status) values ($1,$2,$3,$4,$5) returning tag_id",
		form.TagName, form.CoverUrl, form.Color, *form.Sort, *form.Status)
	var insertId int
	err := row.Scan(&insertId)
	if err == nil {
		slog.Info("新增标签完成，id：" + strconv.Itoa(insertId))
	}
	return err
}

func (t *TagRepo) Update(form *usercase.TagForm) error {
	result, err := t.db.Exec(context.Background(), "update t_blog_tag set update_time = now(), tag_name = $1, cover_url = $2, color = $3, sort = $4, status = $5 where tag_id = $6",
		form.TagName, form.CoverUrl, form.Color, *form.Sort, *form.Status, form.TagId)
	if err == nil {
		slog.Info(fmt.Sprintf("标签更新完成，row:%d,id:%d", result.RowsAffected(), form.TagId))
	}
	return err
}

func (t *TagRepo) UpdateStatus(tagId int, status uint8) error {
	result, err := t.db.Exec(context.Background(), "update t_blog_tag set update_time = now(), status = $1 where tag_id = $2", status, tagId)
	if err == nil {
		slog.Info(fmt.Sprintf("标签状态更新完成，row:%d,id:%d,status:%d", result.RowsAffected(), tagId, status))
	}
	return err
}

func (t *TagRepo) UpdateViewNum(tagId int, addNum int) error {
	result, err := t.db.Exec(context.Background(), "update t_blog_tag set update_time = now(), view_num = view_num + $1 where tag_id = $2", addNum, tagId)
	if err == nil {
		slog.Info(fmt.Sprintf("更新标签查看次数完成，row:%d,id:%d,addnum:%d", result.RowsAffected(), tagId, addNum))
	}
	return err
}

func (t *TagRepo) SelectById(id int) (*usercase.Tag, error) {
	row, err := t.db.Query(context.Background(), "select * from t_blog_tag where tag_id = $1 and delete_at = '0' and status = 0", id)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	for row.Next() {
		tag, err := pgx.RowToStructByName[usercase.Tag](row)
		return &tag, err
	}
	return nil, nil
}

func (t *TagRepo) ManageList(form *usercase.TagQueryForm) ([]*usercase.Tag, error) {
	var builder strings.Builder
	builder.WriteString("select * from t_blog_tag where delete_at = '0'")
	if form.TagName != "" {
		builder.WriteString(fmt.Sprintf(" and tag_name like '%s'", "%"+form.TagName+"%"))
	}
	if form.CreateTimeBegin != nil {
		builder.WriteString(fmt.Sprintf(" and date(create_time) >= '%s'", form.CreateTimeBegin.Format("2006-04-02")))
	}
	if form.CreateTimeEnd != nil {
		builder.WriteString(fmt.Sprintf(" and date(create_time) <= '%s'", form.CreateTimeEnd.Format("2006-04-02")))
	}
	builder.WriteString(" order by sort asc, create_time desc")
	rows, err := t.db.Query(context.Background(), builder.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (*usercase.Tag, error) {
		return pgx.RowToAddrOfStructByName[usercase.Tag](row)
	})
}

func (t *TagRepo) List() ([]*usercase.Tag, error) {
	rows, err := t.db.Query(context.Background(), "select tag_id, tag_name, cover_url, view_num, color, create_time from t_blog_tag where delete_at = '0' and status = 0 order by sort asc, create_time desc")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (*usercase.Tag, error) {
		return pgx.RowToAddrOfStructByNameLax[usercase.Tag](row)
	})
}

func (t *TagRepo) ListByIds(ids []uint) ([]*usercase.Tag, error) {
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
	rows, err := t.db.Query(context.Background(), builder.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (*usercase.Tag, error) {
		return pgx.RowToAddrOfStructByNameLax[usercase.Tag](row)
	})
}

func (t *TagRepo) CountByTagName(name string, tagId uint) (uint8, error) {
	row := t.db.QueryRow(context.Background(), "select count(tag_id) from t_blog_tag where tag_name = $1 and delete_at = '0' and tag_id != $2", name, tagId)
	var total uint8
	err := row.Scan(&total)
	return total, err
}

func (t *TagRepo) DeleteById(id int) error {
	result, err := t.db.Exec(context.Background(), "update t_blog_tag set delete_at = $1 where tag_id = $2",
		strconv.FormatInt(time.Now().UnixMilli(), 10), id)
	if err == nil {
		slog.Info(fmt.Sprintf("删除标签完成，row：%d,id:%d", result.RowsAffected(), id))
	}
	return err
}

func (t *TagRepo) DeleteByIds(ids []int) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	var builder strings.Builder
	builder.WriteString("update t_blog_tag set delete_at = $1 where tag_id in (")
	for i, id := range ids {
		if i > 0 {
			builder.WriteByte(',')
		}
		builder.WriteRune(rune(id))
	}
	builder.WriteByte(')')
	result, err := t.db.Exec(context.Background(), builder.String(), strconv.FormatInt(time.Now().UnixMilli(), 10))
	return result.RowsAffected(), err
}
