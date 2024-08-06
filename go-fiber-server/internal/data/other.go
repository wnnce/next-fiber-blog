package data

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-fiber-ent-web-layout/internal/tools"
	"go-fiber-ent-web-layout/internal/usercase"
	sqlbuild "go-fiber-ent-web-layout/pkg/sql-build"
	"log/slog"
)

type OtherRepo struct {
	db *pgxpool.Pool
}

func NewOtherRepo(data *Data) usercase.IOtherRepo {
	return &OtherRepo{
		db: data.Db,
	}
}

func (self *OtherRepo) SaveFileRecord(file *usercase.UploadFile) {
	builder := sqlbuild.NewInsertBuilder("t_upload_file").
		Fields("file_md5", "origin_name", "file_name", "file_path", "file_size", "file_type").
		Values(file.FileMd5, file.OriginName, file.FileName, file.FilePath, file.FileSize, file.FileType).
		Returning("id")
	row := self.db.QueryRow(context.Background(), builder.Sql(), builder.Args())
	var fileId int64
	if err := row.Scan(&fileId); err != nil {
		slog.Error("保存文件上传记录信息失败", "err", err.Error())
	} else {
		file.ID = fileId
		slog.Info("保存文件上传记录信息成功", "fileId", fileId)
	}
}

func (self *OtherRepo) QueryFileByMd5(fileMd5 string) (*usercase.UploadFile, error) {
	builder := sqlbuild.NewSelectBuilder("t_upload_file").
		Where("file_md5").Eq(fileMd5).BuildAsSelect()
	rows, err := self.db.Query(context.Background(), builder.Sql(), fileMd5)
	if err == nil && rows.Next() {
		defer rows.Close()
		return pgx.RowToAddrOfStructByName[usercase.UploadFile](rows)
	}
	return nil, err
}

func (self *OtherRepo) DeleteFileByName(filename string) error {
	builder := sqlbuild.NewDeleteBuilder("t_upload_file").
		Where("file_path").Eq(filename).BuildSaDelete()
	result, err := self.db.Exec(context.Background(), builder.Sql(), filename)
	if err == nil {
		slog.Info("文件上传记录删除成功", "row", result.RowsAffected())
	}
	return err
}

func (self *OtherRepo) SaveLoginRecord(record *usercase.LoginLog) {
	builder := sqlbuild.NewInsertBuilder("t_login_log").
		Fields("user_id", "user_type", "username", "login_ip", "location", "login_ua", "remark", "result", "login_type").
		Values(record.UserId, record.UserType, record.Username, record.LoginIP, record.Location, record.LoginUa, record.Remark, record.Result, record.LoginType)
	_, err := self.db.Exec(context.Background(), builder.Sql(), builder.Args()...)
	if err != nil {
		slog.Error("保存登录日志失败", "err", err)
	}
}

func (self *OtherRepo) SaveAccessRecord(record *usercase.AccessLog) {
	builder := sqlbuild.NewInsertBuilder("t_blog_access_log").
		Fields("location", "referee", "access_ip", "access_ua").
		Values(record.Location, record.Referee, record.AccessIp, record.AccessUa)
	_, err := self.db.Exec(context.Background(), builder.Sql(), builder.Args()...)
	if err != nil {
		slog.Error("保存访问日志失败", "err", err)
	}
}

func (self *OtherRepo) PageLoginRecord(query *usercase.LoginLogQueryForm) ([]*usercase.LoginLog, int64, error) {
	builder := sqlbuild.NewSelectBuilder("t_login_log").
		WhereByCondition(query.Username != "", "username").Eq(query.Username).
		AndByCondition(query.LoginType != nil, "login_type").Eq(query.LoginType).
		AndByCondition(query.Result != nil, "result").Eq(query.Result).
		AndByCondition(query.CreateTimeBegin != "", "create_time").Ge(query.CreateTimeBegin).
		AndByCondition(query.CreateTimeEnd != "", "create_time").Le(query.CreateTimeEnd).BuildAsSelect().
		OrderBy("create_time desc")
	row := self.db.QueryRow(context.Background(), builder.CountSql(), builder.Args()...)
	var total int64
	if err := row.Scan(&total); err != nil {
		return nil, 0, err
	}
	records := make([]*usercase.LoginLog, 0)
	if total == 0 {
		return records, total, nil
	}
	offset := tools.ComputeOffset(total, query.Page, query.Size, false)
	builder.Limit(int64(query.Size)).Offset(offset)
	rows, err := self.db.Query(context.Background(), builder.Sql(), builder.Args()...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	records, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (*usercase.LoginLog, error) {
		return pgx.RowToAddrOfStructByName[usercase.LoginLog](row)
	})
	return records, total, err
}

func (self *OtherRepo) PageAccessRecord(query *usercase.AccessLogQueryForm) ([]*usercase.AccessLog, int64, error) {
	builder := sqlbuild.NewSelectBuilder("t_blog_access_log").
		WhereByCondition(query.Ip != "", "access_ip").Eq(query.Ip).
		AndByCondition(query.CreateTimeBegin != "", "create_time").Ge(query.CreateTimeBegin).
		AndByCondition(query.CreateTimeEnd != "", "create_time").Le(query.CreateTimeEnd).BuildAsSelect().
		OrderBy("create_time desc")
	row := self.db.QueryRow(context.Background(), builder.CountSql(), builder.Args()...)
	var total int64
	if err := row.Scan(&total); err != nil {
		return nil, 0, err
	}
	records := make([]*usercase.AccessLog, 0)
	if total == 0 {
		return records, total, nil
	}
	offset := tools.ComputeOffset(total, query.Page, query.Size, false)
	builder.Limit(int64(query.Size)).Offset(offset)
	rows, err := self.db.Query(context.Background(), builder.Sql(), builder.Args()...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	records, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (*usercase.AccessLog, error) {
		return pgx.RowToAddrOfStructByName[usercase.AccessLog](row)
	})
	return records, total, err
}
