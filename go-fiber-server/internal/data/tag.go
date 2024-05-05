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

type TagRepo struct {
	db *sqlx.DB
}

func NewTagRepo(data *Data) usercase.ITagRepo {
	return &TagRepo{
		db: data.Db,
	}
}

// Save 保存标签
// form 保存标签的表单参数
func (t *TagRepo) Save(form *usercase.TagForm) error {
	result, err := t.db.Exec("Insert Into t_blog_tag (tag_name, cover_url, color, sort, status) values ($1,$2,$3,$4,$5)",
		form.TagName, form.CoverUrl, form.Color, *form.Sort, *form.Status)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	slog.Info(fmt.Sprintf("标签新增成功，Id:%d", id))
	return nil
}

// Update 更新标签
func (t *TagRepo) Update(form *usercase.TagForm) error {
	result, err := t.db.Exec("update t_blog_tag set update_time = now(), tag_name = $1, cover_url = $2, color = $3, sort = $4, status = $5 where tag_id = $6",
		form.TagName, form.CoverUrl, form.Color, *form.Sort, *form.Status, form.TagId)
	if err != nil {
		return err
	}
	row, _ := result.RowsAffected()
	slog.Info(fmt.Sprintf("标签更新成功，tagId:%d,row:%d", form.TagId, row))
	return nil
}

// UpdateStatus 更新标签状态
func (t *TagRepo) UpdateStatus(tagId int, status uint8) error {
	result, err := t.db.Exec("update t_blog_tag set update_time = now(), status = $1 where tag_id = $2", status, tagId)
	if err != nil {
		return err
	}
	row, _ := result.RowsAffected()
	slog.Info(fmt.Sprintf("更新标签状态完成，tagId：%d，status：%d，row：%d", tagId, status, row))
	return nil
}

// UpdateViewNum 更新标签查看次数
func (t *TagRepo) UpdateViewNum(tagId int, addNum int) error {
	result, err := t.db.Exec("update t_blog_tag set update_time = now(), view_num = view_num + $1 where tag_id = $2", addNum, tagId)
	if err != nil {
		return err
	}
	row, _ := result.RowsAffected()
	slog.Info(fmt.Sprintf("更新标签查看次数完成，tagId:%d,addNum:%d,row:%d", tagId, addNum, row))
	return nil
}

// SelectById 通过Id查找标签
func (t *TagRepo) SelectById(tagId int) (*usercase.Tag, error) {
	tag := &usercase.Tag{}
	if err := t.db.Get(tag, "select * from t_blog_tag where tag_id = $1 and delete_at = '0' and status = 0", tagId); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return nil, nil
		}
		return nil, err
	}
	return tag, nil
}

// ManageList 管理端获取标签列表
func (t *TagRepo) ManageList(form *usercase.TagQueryForm) ([]*usercase.Tag, error) {
	var builder strings.Builder
	builder.WriteString("select * from t_blog_tag where tag_id > 0")
	args := make(map[string]interface{})
	if form.TagName != "" {
		builder.WriteString(" and tag_name like :name")
		args["name"] = "%" + form.TagName + "%"
	}
	if form.CreateTimeBegin != nil {
		builder.WriteString(" and date(create_time) >= :begin")
		args["begin"] = form.CreateTimeBegin.Format("2006-04-02")
	}
	if form.CreateTimeEnd != nil {
		builder.WriteString(fmt.Sprintf(" and date(create_time) <= :end"))
		args["end"] = form.CreateTimeEnd.Format("2006-04-02")
	}
	builder.WriteString(" order by sort asc, create_time desc")
	tags := make([]*usercase.Tag, 0)
	rows, err := t.db.NamedQuery(builder.String(), args)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
	tags = tools.SqlxRowsScan(rows, tags)
	return tags, nil
}

// List 用户端获取所有标签
func (t *TagRepo) List() []*usercase.Tag {
	tags := make([]*usercase.Tag, 0)
	if err := t.db.Select(&tags, "select tag_id, tag_name, cover_url, view_num, color, create_time from t_blog_tag where delete_at = '0' and status = 0 order by sort asc, create_time desc"); err != nil {
		slog.Error(fmt.Sprintf("获取标签列表失败，错误信息：%s", err))
	}
	return tags
}

// ListByIds 通过标签Id列表获取标签
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
	var tags []*usercase.Tag
	if err := t.db.Select(tags, builder.String()); err != nil {
		return nil, err
	}
	return tags, nil
}

// CountByTagName 通过标签名称和标签Id查询数量 判断标签名称是否重复
func (t *TagRepo) CountByTagName(name string, tagId uint) (uint8, error) {
	row := t.db.QueryRow("select count(tag_id) from t_blog_tag where tag_name = $1 and delete_at = '0' and tag_id != $2", name, tagId)
	var total uint8
	err := row.Scan(&total)
	return total, err
}

// DeleteById 通过标签Id删除标签 逻辑删除
func (t *TagRepo) DeleteById(tagId int) error {
	result, err := t.db.Exec("update t_blog_tag set delete_at = $1 where tag_id = $2", time.Now().UnixMilli(), tagId)
	if err != nil {
		return err
	}
	row, _ := result.RowsAffected()
	slog.Info(fmt.Sprintf("删除标签完成，tagId:%d,row:%d", tagId, row))
	return nil
}

// DeleteByIds 通过标签Id列表批量删除标签 逻辑删除
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
	result, err := t.db.Exec(builder.String(), time.Now().UnixMilli())
	if err != nil {
		return 0, err
	}
	row, _ := result.RowsAffected()
	return row, nil
}
