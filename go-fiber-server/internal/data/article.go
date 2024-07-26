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

type ArticleRepo struct {
	db *pgxpool.Pool
}

func NewArticleRepo(data *Data) usercase.IArticleRepo {
	return &ArticleRepo{
		db: data.Db,
	}
}

func (self *ArticleRepo) Save(article *usercase.Article) error {
	sql := `insert into t_blog_article 
    		(title, summary, cover_url, category_ids, tag_ids, content, protocol, tips, password, is_hot, is_top, is_comment, is_private, sort, status) 
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15 ) returning article_id`
	row := self.db.QueryRow(context.Background(), sql, article.Title, article.Summary, article.CoverUrl, article.CategoryIds,
		article.TagIds, article.Content, article.Protocol, article.Tips, article.Password, article.IsHot, article.IsTop,
		article.IsComment, article.IsPrivate, *article.Sort, *article.Status)
	var articleId uint64
	err := row.Scan(&articleId)
	if err == nil {
		slog.Info("保存博客文章完成", "articleId", articleId)
		article.ArticleId = articleId
	}
	return err
}

func (self *ArticleRepo) Update(article *usercase.Article) error {
	var builder strings.Builder
	builder.WriteString(`update t_blog_article set update_time = now(), title = $1, summary = $2, cover_url = $3, category_ids = $4,
                          tag_ids = $5, protocol = $6, tips = $7, password = $8, is_hot = $9, is_top = $10, is_comment = $11, 
                          is_private = $12, sort = $13, status = $14`)
	args := []any{article.Title, article.Summary, article.CoverUrl, article.CategoryIds, article.TagIds, article.Protocol,
		article.Tips, article.Password, article.IsHot, article.IsTop, article.IsComment, article.IsPrivate, *article.Sort,
		*article.Status}
	if strings.TrimSpace(article.Content) != "" {
		args = append(args, article.Content)
		builder.WriteString(fmt.Sprintf(", content = $%d", len(args)))
	}
	builder.WriteString(fmt.Sprintf(" where article_id = $%d", len(args)+1))
	args = append(args, article.ArticleId)
	result, err := self.db.Exec(context.Background(), builder.String(), args...)
	if err == nil {
		slog.Info("更新博客文章完成", "row", result.RowsAffected(), "articleId", article.ArticleId)
	}
	return err
}

func (self *ArticleRepo) UpdateSelective(form *usercase.ArticleUpdateForm) error {
	var builder strings.Builder
	builder.WriteString("update t_blog_article set update_time = now() ")
	args := make([]any, 0)
	if form.IsHot != nil {
		args = append(args, *form.IsHot)
		builder.WriteString(fmt.Sprintf(", is_hot = $%d", len(args)))
	}
	if form.IsTop != nil {
		args = append(args, *form.IsTop)
		builder.WriteString(fmt.Sprintf(", is_top = $%d", len(args)))
	}
	if form.IsComment != nil {
		args = append(args, *form.IsComment)
		builder.WriteString(fmt.Sprintf(", is_comment = $%d", len(args)))
	}
	if form.Status != nil {
		args = append(args, *form.Status)
		builder.WriteString(fmt.Sprintf(", status = $%d", len(args)))
	}
	builder.WriteString(fmt.Sprintf(" where article_id = $%d", len(args)+1))
	args = append(args, form.ArticleId)
	result, err := self.db.Exec(context.Background(), builder.String(), args...)
	if err == nil {
		slog.Info("快捷更新博客文章完成", "row", result.RowsAffected(), "articleId", form.ArticleId)
	}
	return err
}

func (self *ArticleRepo) Page(query *usercase.ArticleQueryForm) ([]*usercase.Article, int64, error) {
	var condition strings.Builder
	condition.WriteString("where delete_at = 0")
	args := make([]any, 0)
	if query.Title != "" {
		args = append(args, "%"+query.Title+"%")
		condition.WriteString(fmt.Sprintf(" and title like $%d", len(args)))
	}
	if query.TagId != nil {
		args = append(args, *query.TagId)
		condition.WriteString(fmt.Sprintf(" and $%d = ANY(tag_ids)", len(args)))
	}
	if query.CategoryId != nil {
		args = append(args, *query.CategoryId)
		condition.WriteString(fmt.Sprintf(" and $%d = ANY(category_ids)", len(args)))
	}
	timeQueryConditionBuilder(query.CreateTimeBegin, query.CreateTimeEnd, &condition, &args)
	row := self.db.QueryRow(context.Background(), "select count(*) from t_blog_article "+condition.String(), args...)
	var total int64
	if err := row.Scan(&total); err != nil {
		return nil, 0, err
	}
	articles := make([]*usercase.Article, 0)
	if total == 0 {
		return articles, 0, nil
	}
	offset := tools.ComputeOffset(total, query.Page, query.Size, false)
	condition.WriteString(fmt.Sprintf("order by is_top desc, sort, create_time desc limit $%d offset $%d", len(args)+1, len(args)+2))
	args = append(args, query.Size, offset)
	rows, err := self.db.Query(context.Background(), `select article_id, title, summary, cover_url, category_ids, 
       tag_ids, view_num, share_num, protocol, tips, password, is_hot, is_top, is_comment, is_private, create_time, sort, 
       status from t_blog_article `+condition.String(), args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	articles, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (*usercase.Article, error) {
		return pgx.RowToAddrOfStructByNameLax[usercase.Article](row)
	})
	return articles, total, err
}

func (self *ArticleRepo) SelectById(articleId int64, checkStatus bool) (*usercase.Article, error) {
	sql := "select * from t_blog_article where article_id = $1 and delete_at = 0"
	if checkStatus {
		sql += " and status = 0"
	}
	rows, err := self.db.Query(context.Background(), sql, articleId)
	if err == nil && rows.Next() {
		return pgx.RowToAddrOfStructByName[usercase.Article](rows)
	}
	return nil, err
}

func (self *ArticleRepo) CountByTagId(tagId int) (int64, error) {
	sql := "select count(*) from t_blog_article where $1 = ANY(tag_ids) and status = 0 and delete_at = 0"
	row := self.db.QueryRow(context.Background(), sql, tagId)
	var total int64
	err := row.Scan(&total)
	return total, err
}

func (self *ArticleRepo) CountByCategoryId(categoryId int) (int64, error) {
	sql := "select count(*) from t_blog_article where $1 = ANY(category_ids) and status = 0 and delete_at = 0"
	row := self.db.QueryRow(context.Background(), sql, categoryId)
	var total int64
	err := row.Scan(&total)
	return total, err
}

func (self ArticleRepo) DeleteById(articleId int64) error {
	sql := "update t_blog_article set delete_at = $1 where article_id = $2"
	result, err := self.db.Exec(context.Background(), sql, time.Now().UnixMilli(), articleId)
	if err == nil {
		slog.Info("删除博客文章完成", "row", result.RowsAffected(), "articleId", articleId)
	}
	return err
}
