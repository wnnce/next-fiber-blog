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
	"strconv"
	"strings"
	"time"
)

type ArticleRepo struct {
	db *pgxpool.Pool
}

func NewArticleRepo(data *Data) usercase.IArticleRepo {
	return &ArticleRepo{
		db: data.Db,
	}
}

func (self *ArticleRepo) Save(article *usercase.Article) error {
	builder := sqlbuild.NewInsertBuilder("t_blog_article").
		Fields("title", "summary", "cover_url", "category_ids", "tag_ids", "content", "protocol", "tips",
			"password", "is_hot", "is_top", "is_comment", "is_private", "sort", "status").
		Values(article.Title, article.Summary, article.CoverUrl, article.CategoryIds, article.TagIds, article.Content,
			article.Protocol, article.Tips, article.Password, article.IsHot, article.IsTop, article.IsComment,
			article.IsPrivate, *article.Sort, *article.Status).
		Returning("article_id")
	row := self.db.QueryRow(context.Background(), builder.Sql(), builder.Args()...)
	var articleId uint64
	err := row.Scan(&articleId)
	if err == nil {
		slog.Info("保存博客文章完成", "articleId", articleId)
		article.ArticleId = articleId
	}
	return err
}

func (self *ArticleRepo) Update(article *usercase.Article) error {
	builder := sqlbuild.NewUpdateBuilder("t_blog_article").
		SetRaw("update_time", "now()").
		SetByMap(map[string]any{
			"title":        article.Title,
			"summary":      article.Summary,
			"cover_url":    article.CoverUrl,
			"category_ids": article.CategoryIds,
			"tag_ids":      article.TagIds,
			"protocol":     article.Protocol,
			"tips":         article.Tips,
			"password":     article.Password,
			"is_hot":       article.IsHot,
			"is_top":       article.IsTop,
			"is_comment":   article.IsComment,
			"is_private":   article.IsPrivate,
			"sort":         article.Sort,
			"status":       article.Status,
		})
	if strings.TrimSpace(article.Content) != "" {
		builder.Set("content", article.Content)
	}
	builder.Where("article_id").Eq(article.ArticleId).BuildAsUpdate()
	result, err := self.db.Exec(context.Background(), builder.Sql(), builder.Args()...)
	if err == nil {
		slog.Info("更新博客文章完成", "row", result.RowsAffected(), "articleId", article.ArticleId)
	}
	return err
}

func (self *ArticleRepo) UpdateSelective(form *usercase.ArticleUpdateForm) error {
	builder := sqlbuild.NewUpdateBuilder("t_blog_article").
		SetRaw("update_time", "now()").
		SetByCondition(form.IsHot != nil, "is_hot", form.IsHot).
		SetByCondition(form.IsTop != nil, "is_top", form.IsTop).
		SetByCondition(form.IsComment != nil, "is_comment", form.IsComment).
		SetByCondition(form.Status != nil, "status", form.Status).
		Where("article_id").Eq(form.ArticleId).BuildAsUpdate()
	result, err := self.db.Exec(context.Background(), builder.Sql(), builder.Args()...)
	if err == nil {
		slog.Info("快捷更新博客文章完成", "row", result.RowsAffected(), "articleId", form.ArticleId)
	}
	return err
}

func (self *ArticleRepo) Page(query *usercase.ArticleQueryForm) ([]*usercase.ArticleVo, int64, error) {
	builder := sqlbuild.NewSelectBuilder("t_blog_article as ba").
		Select("ba.article_id", "ba.title", "ba.summary", "ba.cover_url", "ba.category_ids", "ba.tag_ids",
			"ba.view_num", "ba.share_num", "ba.vote_up", "ba.protocol", "ba.tips", "ba.password", "ba.is_hot", "ba.is_top",
			"ba.is_comment", "ba.is_private", "ba.create_time", "ba.sort", "ba.status").
		LeftJoin("t_blog_comment as bc").On("bc.article_id").EqRaw("ba.article_id").And("bc.status").EqRaw("0").And("bc.delete_at").EqRaw("0").BuildAsSelect().
		Select("count(DISTINCT bc.comment_id) as comment_num").
		LeftJoin("t_blog_category as ct").On("ct.category_id").EqRaw("ANY(ba.category_ids)").And("ct.status").EqRaw("0").And("ct.delete_at").EqRaw("0").BuildAsSelect().
		// 直接使用 jsonb_agg 将分类字段聚合为分类列表
		// json字段需要使用驼峰命名 因为pgx在处理jsonb类型字段时 会调用json库直接反序列化
		Select("jsonb_agg(DISTINCT jsonb_build_object('categoryId', ct.category_id, 'categoryName', ct.category_name)) AS categories").
		LeftJoin("t_blog_tag as bt").On("bt.tag_id").EqRaw("ANY(ba.tag_ids)").And("bt.status").EqRaw("0").And("bt.delete_at").EqRaw("0").BuildAsSelect().
		// 聚合为标签列表
		Select("jsonb_agg(DISTINCT jsonb_build_object('tagId', bt.tag_id, 'tagName', bt.tag_name, 'color', bt.color)) AS tags").
		WhereByCondition(query.Title != "", "ba.title").Like("%"+query.Title+"%").
		AndByCondition(query.TagId != nil, fmt.Sprintf("%d", query.TagId)).EqRaw("ANY(ba.tag_ids)").
		AndByCondition(query.CategoryId != nil, fmt.Sprintf("%d", query.CategoryId)).EqRaw("ANY(ba.category_ids)").
		AndByCondition(query.Status != nil, "ba.status").Eq(query.Status).
		AndByCondition(query.CreateTimeBegin != "", "ba.create_time").Ge(query.CreateTimeBegin).
		AndByCondition(query.CreateTimeEnd != "", "ba.create_time").Le(query.CreateTimeEnd).
		And("ba.delete_at").EqRaw("0").BuildAsSelect().
		GroupBy("ba.article_id")
	row := self.db.QueryRow(context.Background(), builder.CountSql(), builder.Args()...)
	var total int64
	if err := row.Scan(&total); err != nil {
		return nil, 0, err
	}
	articles := make([]*usercase.ArticleVo, 0)
	if total == 0 {
		return articles, 0, nil
	}
	offset := tools.ComputeOffset(total, query.Page, query.Size, false)
	builder.Limit(int64(query.Size)).Offset(offset)
	rows, err := self.db.Query(context.Background(), builder.Sql(), builder.Args()...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	articles, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (*usercase.ArticleVo, error) {
		return pgx.RowToAddrOfStructByNameLax[usercase.ArticleVo](row)
	})
	return articles, total, err
}

func (self *ArticleRepo) SelectById(articleId uint64, checkStatus bool) (*usercase.ArticleVo, error) {
	builder := sqlbuild.NewSelectBuilder("t_blog_article as ba").
		Select("ba.*").
		LeftJoin("t_blog_comment as bc").On("bc.article_id").EqRaw("ba.article_id").And("bc.status").EqRaw("0").And("bc.delete_at").EqRaw("0").BuildAsSelect().
		Select("count(DISTINCT bc.comment_id) as comment_num").
		LeftJoin("t_blog_category as ct").On("ct.category_id").EqRaw("ANY(ba.category_ids)").And("ct.status").EqRaw("0").And("ct.delete_at").EqRaw("0").BuildAsSelect().
		Select("jsonb_agg(DISTINCT jsonb_build_object('categoryId', ct.category_id, 'categoryName', ct.category_name)) AS categories").
		LeftJoin("t_blog_tag as bt").On("bt.tag_id").EqRaw("ANY(ba.tag_ids)").And("bt.status").EqRaw("0").And("bt.delete_at").EqRaw("0").BuildAsSelect().
		Select("jsonb_agg(DISTINCT jsonb_build_object('tagId', bt.tag_id, 'tagName', bt.tag_name, 'color', bt.color)) AS tags").
		Where("ba.article_id").Eq(articleId).
		AndByCondition(checkStatus, "ba.status").EqRaw("0").
		And("ba.delete_at").EqRaw("0").BuildAsSelect().
		GroupBy("ba.article_id")
	rows, err := self.db.Query(context.Background(), builder.Sql(), articleId)
	if err == nil && rows.Next() {
		return pgx.RowToAddrOfStructByNameLax[usercase.ArticleVo](rows)
	}
	return nil, err
}

func (self *ArticleRepo) CountByTagId(tagId int) (int64, error) {
	builder := sqlbuild.NewSelectBuilder("t_blog_article").
		Select("count(*)").
		Where(strconv.Itoa(tagId)).EqRaw("ANY(tag_ids)").
		And("delete_at").EqRaw("0").BuildAsSelect()
	row := self.db.QueryRow(context.Background(), builder.Sql())
	var total int64
	err := row.Scan(&total)
	return total, err
}

func (self *ArticleRepo) CountByCategoryId(categoryId int) (int64, error) {
	builder := sqlbuild.NewSelectBuilder("t_blog_article").
		Select("count(*)").
		Where(strconv.Itoa(categoryId)).EqRaw("ANY(category_ids)").
		And("delete_at").EqRaw("0").BuildAsSelect()
	row := self.db.QueryRow(context.Background(), builder.Sql())
	var total int64
	err := row.Scan(&total)
	return total, err
}

func (self *ArticleRepo) CountByTitle(title string, articleId uint64) (uint8, error) {
	builder := sqlbuild.NewSelectBuilder("t_blog_article").
		Select("count(*)").
		Where("article_id").Ne(articleId).
		And("title").Eq(title).
		And("delete_at").EqRaw("0").BuildAsSelect()
	row := self.db.QueryRow(context.Background(), builder.Sql(), builder.Args()...)
	var total uint8
	err := row.Scan(&total)
	return total, err
}

func (self ArticleRepo) DeleteById(articleId uint64) error {
	builder := sqlbuild.NewUpdateBuilder("t_blog_article").
		Set("delete_at", time.Now().UnixMilli()).
		Where("article_id").Eq(articleId).BuildAsUpdate()
	result, err := self.db.Exec(context.Background(), builder.Sql(), builder.Args()...)
	if err == nil {
		slog.Info("删除博客文章完成", "row", result.RowsAffected(), "articleId", articleId)
	}
	return err
}
