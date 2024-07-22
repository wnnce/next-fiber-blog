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
	"strings"
)

var InjectSet = wire.NewSet(NewData, NewRedisTemplate, NewTagRepo, NewCategoryRepo, NewConcatRepo, NewLinkRepo, NewSysMenuRepo, NewOtherRepo,
	NewSysConfigRepo, NewSysRoleRepo, NewSysUserRepo, NewSysDictRepo, NewNoticeRepo)

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

// commonFieldUpdateBuilder 动态构建通过字段的更新sql
func commonFieldUpdateBuilder(sort *uint, status *uint8, builder *strings.Builder, args *[]any) {
	if sort != nil {
		*args = append(*args, *sort)
		builder.WriteString(fmt.Sprintf(", sort = $%d", len(*args)))
	}
	if status != nil {
		*args = append(*args, *status)
		builder.WriteString(fmt.Sprintf(", status = $%d", len(*args)))
	}
}

// timeQueryConditionBuilder 动态构建时间查询SQL
func timeQueryConditionBuilder(begin, end string, builder *strings.Builder, args *[]any) {
	if begin != "" {
		*args = append(*args, begin)
		builder.WriteString(fmt.Sprintf(" and create_time >= $%d", len(*args)))
	}
	if end != "" {
		*args = append(*args, end)
		builder.WriteString(fmt.Sprintf(" and create_time <= $%d", len(*args)))
	}
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
