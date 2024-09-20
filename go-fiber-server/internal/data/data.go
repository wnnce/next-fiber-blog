package data

import (
	"context"
	"fmt"
	"github.com/google/wire"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go-fiber-ent-web-layout/internal/conf"
	"go-fiber-ent-web-layout/internal/tools"
	"go-fiber-ent-web-layout/internal/usercase"
	sqlbuild "go-fiber-ent-web-layout/pkg/sql-build"
	"math"
)

var InjectSet = wire.NewSet(NewData, NewRedisTemplate, NewTagRepo, NewCategoryRepo, NewConcatRepo, NewLinkRepo,
	NewSysMenuRepo, NewOtherRepo, NewSysConfigRepo, NewSysRoleRepo, NewSysUserRepo, NewSysDictRepo, NewNoticeRepo,
	NewArticleRepo, NewTopicRepo, NewUserRepo, NewCommentRepo)

type Data struct {
	Db *pgxpool.Pool // pgx连接
	Rc *redis.Client // 封装的redis操作
}

func NewData(conf *conf.Data) (*Data, func(), error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", conf.Database.Host, conf.Database.Port, conf.Database.Username, conf.Database.Password, conf.Database.DbName)
	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, nil, err
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", conf.Redis.Host, conf.Redis.Port),
		DB:           conf.Redis.Index,
		Username:     conf.Redis.Username,
		Password:     conf.Redis.Password,
		ReadTimeout:  conf.Redis.ReadTimeout,
		WriteTimeout: conf.Redis.WireTimeout,
	})
	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		_ = rdb.Close()
		db.Close()
	}
	return &Data{
		Db: db,
		Rc: rdb,
	}, cleanup, nil
}

// SelectPage 通用泛型分页查询方法
func SelectPage[T any](builder sqlbuild.SelectBuilder, page, size int, safe bool, db *pgxpool.Pool) (*usercase.PageData[T], error) {
	var total int64
	row := db.QueryRow(context.Background(), builder.CountSql(), builder.Args()...)
	if err := row.Scan(&total); err != nil {
		return nil, err
	}
	if total == 0 {
		return &usercase.PageData[T]{
			Current: page,
			Size:    size,
			Total:   total,
			Pages:   0,
			Records: make([]*T, 0),
		}, nil
	}
	offset := tools.ComputeOffset(total, page, size, safe)
	builder.Limit(int64(size)).Offset(offset)
	rows, err := db.Query(context.Background(), builder.Sql(), builder.Args()...)
	if err != nil {
		return nil, err
	}
	records, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*T, error) {
		return pgx.RowToAddrOfStructByNameLax[T](row)
	})
	if err != nil {
		return nil, err
	}
	pages := int(math.Ceil(float64(total) / float64(size)))
	if page > pages && safe {
		page = pages
	}
	return &usercase.PageData[T]{
		Current: page,
		Size:    size,
		Total:   total,
		Pages:   pages,
		Records: records,
	}, nil
}

// smartExec 通用执行方法 用于区分事务执行和普通执行
func smartExec(pool *pgxpool.Pool, tx pgx.Tx, ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	if tx == nil {
		return pool.Exec(ctx, sql, args...)
	} else {
		return tx.Exec(ctx, sql, args...)
	}
}

// smartQuery 通用查询方法 区分事务执行和普通执行
func smartQuery(pool *pgxpool.Pool, tx pgx.Tx, ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	if tx == nil {
		return pool.Query(ctx, sql, args...)
	} else {
		return tx.Query(ctx, sql, args...)
	}
}

// smartQueryRow 通用查询方法
func smartQueryRow(pool *pgxpool.Pool, tx pgx.Tx, ctx context.Context, sql string, args ...any) pgx.Row {
	if tx == nil {
		return pool.QueryRow(ctx, sql, args...)
	} else {
		return tx.QueryRow(ctx, sql, args...)
	}
}
