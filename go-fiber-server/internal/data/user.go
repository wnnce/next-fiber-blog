package data

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-fiber-ent-web-layout/internal/usercase"
	sqlbuild "go-fiber-ent-web-layout/pkg/sql-build"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(data *Data) usercase.IUserRepo {
	return &UserRepo{
		db: data.Db,
	}
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
		Values(user.NickName, user.Summary, user.Avatar, user.Email, user.Link, user.Username, 0).
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
