package data

import (
	"context"
	"fmt"
	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"go-fiber-ent-web-layout/internal/conf"
)

var InjectSet = wire.NewSet(NewData, NewTagRepo, NewCategoryRepo, NewConcatRepo, NewLinkRepo)

type Data struct {
	Db *sqlx.DB      // gorm连接
	Rc *redis.Client // 封装的redis操作
}

func NewData(conf *conf.Data) (*Data, func(), error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", conf.Database.Host, conf.Database.Port, conf.Database.Username, conf.Database.Password, conf.Database.DbName)
	db, err := sqlx.Connect(conf.Database.Driver, dsn)
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
	}
	return &Data{
		Db: db,
		Rc: rdb,
	}, cleanup, nil
}
