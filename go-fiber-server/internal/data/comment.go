package data

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-fiber-ent-web-layout/internal/usercase"
	sqlbuild "go-fiber-ent-web-layout/pkg/sql-build"
)

type CommentRepo struct {
	db *pgxpool.Pool
}

func NewCommentRepo(data *Data) usercase.ICommentRepo {
	return &CommentRepo{
		db: data.Db,
	}
}

func (self *CommentRepo) Save(comment *usercase.Comment) error {
	builder := sqlbuild.NewInsertBuilder("t_blog_comment").
		Fields("content", "user_id", "fid", "rid", "location", "comment_ip", "comment_ua", "comment_type").
		Values(comment.Content, comment.UserId, comment.Fid, comment.Rid, comment.Location, comment.CommentIp,
			comment.CommentUa, comment.CommentType).
		InsertByCondition(comment.TopicId > 0, "topic_id", comment.TopicId).
		InsertByCondition(comment.ArticleId > 0, "article_id", comment.ArticleId).
		Returning("comment_id")
	var commentId int64
	row := self.db.QueryRow(context.Background(), builder.Sql(), builder.Args()...)
	err := row.Scan(&commentId)
	if err == nil {
		comment.CommentId = commentId
	}
	return err
}

func (self *CommentRepo) Page(query *usercase.CommentQueryForm) (*usercase.PageData[usercase.CommentVo], error) {
	builder := sqlbuild.NewSelectBuilder("t_blog_comment as bc").
		Select("bc.comment_id", "bc.content", "bc.user_id", "bc.fid", "bc.rid", "bc.location", "bc.comment_ua", "bc.vote_up", "bc.comment_type", "bc.is_hot", "bc.is_top", "bc.is_coll", "bc.create_time").
		LeftJoin("t_blog_user as bu").On("bu.user_id").EqRaw("bc.user_id").BuildAsSelect().
		LeftJoin("t_blog_user_extend as bue").On("bue.user_id").EqRaw("bc.user_id").BuildAsSelect().
		Select("jsonb_build_object('nickname', bu.nick_name, 'avatar', bu.avatar, 'level', bue.level, 'labels', bu.labels, 'link', bu.link) as user").
		LeftJoin("t_blog_comment as rbc").On("rbc.comment_id").EqRaw("bc.rid").And("bc.rid").LeRaw("0").BuildAsSelect().
		LeftJoin("t_blog_user as rbu").On("rbu.user_id").EqRaw("rbc.user_id").BuildAsSelect().
		LeftJoin("t_blog_user_extend as rbue").On("rbue.user_id").EqRaw("rbc.user_id").BuildAsSelect().
		Select("CASE WHEN bc.rid > 0 THEN jsonb_build_object('nickname', rbu.nick_name, 'avatarUrl', rbu.avatar, 'level', rbue.level, 'labels', rbu.labels, 'link', rbu.link) END as parentUser").
		WhereByCondition(query.CommentType > 0, "bc.comment_type").Eq(query.CommentType).
		AndByCondition(query.ArticleId > 0, "bc.article_id").Eq(query.ArticleId).
		AndByCondition(query.TopicId > 0, "bc.topic_id").Eq(query.TopicId).
		And("bc.fid").Eq(query.Fid).
		And("bc.delete_at").EqRaw("0").
		And("bc.status").EqRaw("0").BuildAsSelect().
		OrderBy("bc.is_top desc", "bc.sort", "bc.create_time")
	return SelectPage[usercase.CommentVo](builder, query.Page, query.Size, true, self.db)
}

func (self *CommentRepo) TotalComment(query *usercase.CommentQueryForm) (uint64, error) {
	builder := sqlbuild.NewSelectBuilder("t_blog_comment").
		Select("count(*) as total").
		WhereByCondition(query.CommentType > 0, "comment_type").Eq(query.CommentType).
		AndByCondition(query.ArticleId > 0, "article_id").Eq(query.ArticleId).
		AndByCondition(query.TopicId > 0, "topic_id").Eq(query.TopicId).
		And("status").EqRaw("0").
		And("delete_at").EqRaw("0").BuildAsSelect()
	var total uint64
	row := self.db.QueryRow(context.Background(), builder.Sql(), builder.Args()...)
	err := row.Scan(&total)
	return total, err
}
