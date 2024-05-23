package data

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-fiber-ent-web-layout/internal/usercase"
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

func (ot *OtherRepo) SaveFileRecord(file *usercase.UploadFile) {
	row := ot.db.QueryRow(context.Background(), "insert into t_upload_file (file_md5, origin_name, file_name, file_path, file_size, file_type) VALUES ($1, $2, $3, $4, $5, $6) returning id",
		file.FileMd5, file.OriginName, file.FileName, file.FilePath, file.FileSize, file.FileType)
	var fileId int64
	if err := row.Scan(&fileId); err != nil {
		slog.Error("保存文件上传记录信息失败", "err", err.Error())
	} else {
		slog.Info("保存文件上传记录信息成功", "fileId", fileId)
	}
}

func (ot *OtherRepo) QueryFileByMd5(fileMd5 string) (*usercase.UploadFile, error) {
	rows, err := ot.db.Query(context.Background(), "select * from t_upload_file where file_md5 = $1", fileMd5)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		return pgx.RowToAddrOfStructByName[usercase.UploadFile](rows)
	}
	return nil, nil
}

func (ot *OtherRepo) DeleteFileByName(filename string) error {
	result, err := ot.db.Exec(context.Background(), "delete from t_upload_file where file_path = $1", filename)
	if err == nil {
		slog.Info("文件上传记录删除成功", "row", result.RowsAffected())
	}
	return err
}

func (ot *OtherRepo) SaveLoginRecord(record *usercase.LoginLog) {
	_, err := ot.db.Exec(context.Background(), "insert into t_login_log (user_id, user_type, username, login_ip, location, login_ua, remark, result, login_type) values ($1, $2, $3, $4, $5, $6, $7, $8, $9 ) returning id",
		record.UserId, record.UserType, record.Username, record.LoginIP, record.Location, record.LoginUa, record.Remark, record.Result, record.LoginType)
	if err != nil {
		slog.Error("保存登录日志失败", "err", err)
	}
}

func (ot *OtherRepo) SaveAccessRecord(record *usercase.AccessLog) {
	_, err := ot.db.Exec(context.Background(), "insert into t_blog_access_log (location, referee, access_ip, access_ua) values ($1, $2, $3, $4)",
		record.Location, record.Referee, record.AccessIp, record.AccessUa)
	if err != nil {
		slog.Error("保存访问日志失败", "err", err)
	}
}
