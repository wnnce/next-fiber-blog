package data

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-fiber-ent-web-layout/internal/tools"
	"go-fiber-ent-web-layout/internal/usercase"
	sqlbuild "go-fiber-ent-web-layout/pkg/sql-build"
	"log/slog"
	"strconv"
	"time"
)

type TopicRepo struct {
	db *pgxpool.Pool
}

func NewTopicRepo(data *Data) usercase.ITopicRepo {
	return &TopicRepo{
		db: data.Db,
	}
}

func (self *TopicRepo) Save(topic *usercase.Topic) error {
	builder := sqlbuild.NewInsertBuilder("t_blog_topic").
		Fields("content", "image_urls", "location", "is_hot", "is_top", "mode", "sort", "status").
		Values(topic.Content, topic.ImageUrls, topic.Location, topic.IsHot, topic.IsTop, topic.Mode, *topic.Sort, *topic.Status).
		Returning("topic_id")
	row := self.db.QueryRow(context.Background(), builder.Sql(), builder.Args()...)
	var topicId uint64
	err := row.Scan(&topicId)
	if err == nil {
		topic.TopicId = topicId
		slog.Info("保存博客动态完成", "topicId", topicId)
	}
	return err
}

func (self *TopicRepo) Update(topic *usercase.Topic) error {
	builder := sqlbuild.NewUpdateBuilder("t_blog_topic").
		SetRaw("update_time", "now()").
		SetByMap(map[string]any{
			"content":    topic.Content,
			"image_urls": topic.ImageUrls,
			"location":   topic.Location,
			"is_hot":     topic.IsHot,
			"is_top":     topic.IsTop,
			"mode":       topic.Mode,
			"sort":       *topic.Sort,
			"status":     *topic.Status,
		}).
		Where("topic_id").Eq(topic.TopicId).BuildAsUpdate()
	result, err := self.db.Exec(context.Background(), builder.Sql(), builder.Args()...)
	if err == nil {
		slog.Info("更新博客动态完成", "row", result.RowsAffected(), "topicId", topic.TopicId)
	}
	return err
}

func (self *TopicRepo) UpdateSelective(form *usercase.TopicUpdateForm) error {
	builder := sqlbuild.NewUpdateBuilder("t_blog_topic").
		SetRaw("update_time", "now()").
		SetByCondition(form.IsTop != nil, "is_top", form.IsTop).
		SetByCondition(form.IsHot != nil, "is_hot", form.IsHot).
		SetByCondition(form.Status != nil, "status", form.Status).
		Where("topic_id").Eq(form.TopicId).BuildAsUpdate()
	result, err := self.db.Exec(context.Background(), builder.Sql(), builder.Args()...)
	if err == nil {
		slog.Info("快捷更新博客动态完成", "row", result.RowsAffected(), "topicId", form.TopicId)
	}
	return err
}

func (self *TopicRepo) Page(query *usercase.TopicQueryForm) ([]*usercase.Topic, int64, error) {
	var selectFields []string
	if query.IsAdmin {
		selectFields = []string{"*"}
	} else {
		selectFields = []string{"topic_id", "content", "image_urls", "location", "is_hot", "is_top", "vote_up", "mode", "create_time"}
	}
	builder := sqlbuild.NewSelectBuilder("t_blog_topic").
		Select(selectFields...).
		WhereByCondition(query.Location != "", "location").Eq(query.Location).
		AndByCondition(query.Status != nil, "status").Eq(query.Status).
		AndByCondition(query.CreateTimeBegin != "", "create_time").Ge(query.CreateTimeBegin).
		AndByCondition(query.CreateTimeEnd != "", "create_time").Le(query.CreateTimeEnd).
		And("delete_at").EqRaw("0").BuildAsSelect().
		OrderBy("is_top desc", "sort", "create_time desc")
	var total int64
	row := self.db.QueryRow(context.Background(), builder.CountSql(), builder.Args()...)
	if err := row.Scan(&total); err != nil {
		return nil, 0, err
	}
	topics := make([]*usercase.Topic, 0)
	if total == 0 {
		return topics, 0, nil
	}
	offset := tools.ComputeOffset(total, query.Page, query.Size, true)
	builder.Limit(int64(query.Size)).Offset(offset)
	rows, err := self.db.Query(context.Background(), builder.Sql(), builder.Args()...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	topics, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (*usercase.Topic, error) {
		return pgx.RowToAddrOfStructByNameLax[usercase.Topic](row)
	})
	return topics, total, err
}

func (self *TopicRepo) VoteUp(topicId uint64, num int) error {
	builder := sqlbuild.NewUpdateBuilder("t_blog_topic").
		SetRaw("update_time", "now()").
		SetRaw("vote_up", "vote_up + "+strconv.Itoa(num)).
		Where("topic_id").Eq(topicId).BuildAsUpdate()
	result, err := self.db.Exec(context.Background(), builder.Sql(), builder.Args()...)
	if err == nil {
		slog.Info("更新动态点赞数完成", "rows", result.RowsAffected(), "topicId", topicId)
	}
	return err
}

func (self *TopicRepo) DeleteById(topicId int64) error {
	builder := sqlbuild.NewUpdateBuilder("t_blog_topic").
		Set("delete_at", time.Now().UnixMilli()).
		Where("topic_id").Eq(topicId).BuildAsUpdate()
	result, err := self.db.Exec(context.Background(), builder.Sql(), builder.Args()...)
	if err == nil {
		slog.Info("删除博客动态完成", "row", result.RowsAffected(), "topicId", topicId)
	}
	return err
}
