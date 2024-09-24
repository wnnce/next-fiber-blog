package data

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-fiber-ent-web-layout/internal/usercase"
	sqlbuild "go-fiber-ent-web-layout/pkg/sql-build"
	"log/slog"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(data *Data) usercase.IUserRepo {
	return &UserRepo{
		db: data.Db,
	}
}

func (self *UserRepo) Transaction(ctx context.Context, fn func(tx pgx.Tx) error) error {
	return pgx.BeginFunc(ctx, self.db, func(tx pgx.Tx) error {
		return fn(tx)
	})
}

func (self *UserRepo) QueryUserByUsername(username string) (*usercase.User, error) {
	builder := sqlbuild.NewSelectBuilder("t_blog_user").
		Select("user_id", "nick_name", "summary", "avatar", "email", "link", "username", "labels", "user_type", "create_time", "status").
		Where("username").Eq(username).BuildAsSelect()
	rows, err := self.db.Query(context.Background(), builder.Sql(), builder.Args()...)
	if err == nil && rows.Next() {
		defer rows.Close()
		return pgx.RowToAddrOfStructByNameLax[usercase.User](rows)
	}
	return nil, err
}

func (self *UserRepo) Save(user *usercase.UserVo) error {
	builder := sqlbuild.NewInsertBuilder("t_blog_user").
		Fields("nick_name", "summary", "avatar", "email", "link", "username", "user_type").
		Values(user.Nickname, user.Summary, user.Avatar, user.Email, user.Link, user.Username, 0).
		Returning("user_id")
	return pgx.BeginFunc(context.Background(), self.db, func(tx pgx.Tx) error {
		row := self.db.QueryRow(context.Background(), builder.Sql(), builder.Args()...)
		var userId uint64
		err := row.Scan(&userId)
		if err != nil {
			return err
		}
		user.User.UserId = userId
		extendBuilder := sqlbuild.NewInsertBuilder("t_blog_user_extend").
			Fields("user_id", "level", "register_ip", "register_location").
			Values(userId, user.Level, user.RegisterIp, user.RegisterLocation)
		_, err = tx.Exec(context.Background(), extendBuilder.Sql(), extendBuilder.Args()...)
		return err
	})

}

func (self *UserRepo) QueryUserExtendById(userId uint64) (*usercase.UserExtend, error) {
	builder := sqlbuild.NewSelectBuilder("t_blog_user_extend").
		Select("level", "expertise", "register_ip", "register_location").
		Where("user_id").Eq(userId).BuildAsSelect()
	rows, err := self.db.Query(context.Background(), builder.Sql(), builder.Args()...)
	if err == nil && rows.Next() {
		return pgx.RowToAddrOfStructByNameLax[usercase.UserExtend](rows)
	}
	return nil, err
}

func (self *UserRepo) SaveExpertiseDetail(detail *usercase.ExpertiseDetail, tx pgx.Tx) error {
	builder := sqlbuild.NewInsertBuilder("t_expertise_detail").
		Fields("user_id", "detail", "detail_type", "source").
		Values(detail.UserId, detail.Detail, detail.DetailType, detail.Source).
		InsertByCondition("" != detail.Remark, "remark", detail.Remark).
		Returning("id")
	var detailId uint64
	row := smartQueryRow(self.db, tx, context.Background(), builder.Sql(), builder.Args()...)
	err := row.Scan(&detailId)
	if err == nil {
		detail.ID = detailId
	}
	return err
}

func (self *UserRepo) UpdateUserExpertise(count int64, userId uint64, tx pgx.Tx) (uint64, uint8, error) {
	builder := sqlbuild.NewUpdateBuilder("t_blog_user_extend").
		SetRaw("expertise", fmt.Sprintf("expertise + %d", count)).
		Where("user_id").Eq(userId).BuildAsUpdate().
		Returning("expertise", "level")
	row := smartQueryRow(self.db, tx, context.Background(), builder.Sql(), builder.Args()...)
	var expertise uint64
	var level uint8
	err := row.Scan(&expertise, &level)
	return expertise, level, err
}

func (self *UserRepo) UpdateUserLevel(level uint8, userId uint64, tx pgx.Tx) error {
	builder := sqlbuild.NewUpdateBuilder("t_user_extend").
		Set("level", level).
		Where("user_id").Eq(userId).BuildAsUpdate()
	result, err := smartExec(self.db, tx, context.Background(), builder.Sql(), builder.Args()...)
	if err == nil {
		slog.Info("更新用户等级完成", "row", result.RowsAffected(), "userId", userId, "level", level)
	}
	return err
}
