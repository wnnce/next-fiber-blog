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
	row := self.db.QueryRow(context.Background(), "insert into t_upload_file (file_md5, origin_name, file_name, file_path, file_size, file_type) VALUES ($1, $2, $3, $4, $5, $6) returning id",
		file.FileMd5, file.OriginName, file.FileName, file.FilePath, file.FileSize, file.FileType)
	var fileId int64
	if err := row.Scan(&fileId); err != nil {
		slog.Error("保存文件上传记录信息失败", "err", err.Error())
	} else {
		slog.Info("保存文件上传记录信息成功", "fileId", fileId)
	}
}

func (self *OtherRepo) QueryFileByMd5(fileMd5 string) (*usercase.UploadFile, error) {
	rows, err := self.db.Query(context.Background(), "select * from t_upload_file where file_md5 = $1", fileMd5)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		return pgx.RowToAddrOfStructByName[usercase.UploadFile](rows)
	}
	return nil, nil
}

func (self *OtherRepo) DeleteFileByName(filename string) error {
	result, err := self.db.Exec(context.Background(), "delete from t_upload_file where file_path = $1", filename)
	if err == nil {
		slog.Info("文件上传记录删除成功", "row", result.RowsAffected())
	}
	return err
}

func (self *OtherRepo) SaveLoginRecord(record *usercase.LoginLog) {
	_, err := self.db.Exec(context.Background(), "insert into t_login_log (user_id, user_type, username, login_ip, location, login_ua, remark, result, login_type) values ($1, $2, $3, $4, $5, $6, $7, $8, $9 ) returning id",
		record.UserId, record.UserType, record.Username, record.LoginIP, record.Location, record.LoginUa, record.Remark, record.Result, record.LoginType)
	if err != nil {
		slog.Error("保存登录日志失败", "err", err)
	}
}

func (self *OtherRepo) SaveAccessRecord(record *usercase.AccessLog) {
	_, err := self.db.Exec(context.Background(), "insert into t_blog_access_log (location, referee, access_ip, access_ua) values ($1, $2, $3, $4)",
		record.Location, record.Referee, record.AccessIp, record.AccessUa)
	if err != nil {
		slog.Error("保存访问日志失败", "err", err)
	}
}

func (self *OtherRepo) PageLoginRecord(query *usercase.LoginLogQueryForm) ([]*usercase.LoginLog, int64, error) {
	var condition strings.Builder
	condition.WriteString(" where id > 0")
	if query.Username != "" {
		condition.WriteString(fmt.Sprintf(" and username = '%s'", query.Username))
	}
	if query.LoginType != nil {
		condition.WriteString(fmt.Sprintf(" and login_type = %d", *query.LoginType))
	}
	if query.Result != nil {
		condition.WriteString(fmt.Sprintf(" and result = %d", *query.Result))
	}
	if query.CreateTimeBegin != "" {
		condition.WriteString(fmt.Sprintf(" and create_time >= '%s'", query.CreateTimeBegin))
	}
	if query.CreateTimeEnd != "" {
		condition.WriteString(fmt.Sprintf(" and create_time <= '%s'", query.CreateTimeEnd))
	}
	row := self.db.QueryRow(context.Background(), "select count(*) from t_login_log "+condition.String())
	var total int64
	if err := row.Scan(&total); err != nil {
		return nil, 0, err
	}
	records := make([]*usercase.LoginLog, 0)
	if total == 0 {
		return records, total, nil
	}
	offset := tools.ComputeOffset(total, query.Page, query.Size, false)
	condition.WriteString(" order by create_time desc limit $1 offset $2")
	rows, err := self.db.Query(context.Background(), "select * from t_login_log "+condition.String(), query.Size, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	records, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (*usercase.LoginLog, error) {
		return pgx.RowToAddrOfStructByName[usercase.LoginLog](row)
	})
	return records, total, err
}
