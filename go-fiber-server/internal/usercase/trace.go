package usercase

import (
	"mime/multipart"
	"time"
)

// UploadFile 文件上传记录
type UploadFile struct {
	ID         int64      `json:"id,omitempty" db:"id"`                  // 主键ID
	FileMd5    string     `json:"fileMd5,omitempty" db:"file_md5"`       // 文件MD5
	OriginName string     `json:"originName,omitempty" db:"origin_name"` // 文件原始名称
	FileName   string     `json:"fileName,omitempty" db:"file_name"`     // 格式化后的文件名称
	FilePath   string     `json:"filePath,omitempty" db:"file_path"`     // 文件的存储路径
	FileSize   int64      `json:"fileSize,omitempty" db:"file_size"`     // 文件大小
	FileType   string     `json:"fileType,omitempty" db:"file_type"`     // 文件类型
	UpdateTime *time.Time `json:"uploadTime,omitempty" db:"upload_time"` // 文件的保存时间
}

// LoginLog 登录日志
type LoginLog struct {
	ID         int64      `json:"id" db:"id"`                  // 主键ID
	UserId     uint64     `json:"userId" db:"user_id"`         // 用户Id
	UserType   int        `json:"userType" db:"user_type"`     // 用户类型
	Username   string     `json:"username" db:"username"`      // 用户名
	LoginIP    string     `json:"loginIp" db:"login_ip"`       // 登录IP
	Location   string     `json:"location" db:"location"`      // 登录位置
	LoginUa    string     `json:"loginUa" db:"login_ua"`       // 登录UA
	LoginType  uint8      `json:"loginType" db:"login_type"`   // 登录类型 1:博客登录 2:管理端登录
	CreateTime *time.Time `json:"createTime" db:"create_time"` // 登录时间
	Remark     string     `json:"remark" db:"remark"`          // 备注
	Result     int        `json:"result" db:"result"`          // 登录结果 0：成功 1：失败
}

// AccessLog 访问日志
type AccessLog struct {
	ID         int64      `json:"id" db:"id"`                  // 主键Id
	Location   string     `json:"location" db:"location"`      // 访问地址
	Referee    string     `json:"referee" db:"referee"`        // 访问来源
	AccessIp   string     `json:"accessIp" db:"access_ip"`     // 访问IP
	AccessUa   string     `json:"accessUa" db:"access_ua"`     // 访问UA
	CreateTime *time.Time `json:"createTime" db:"create_time"` // 访问时间
}

type IOtherRepo interface {
	SaveFileRecord(file *UploadFile)

	QueryFileByMd5(fileMd5 string) (*UploadFile, error)

	DeleteFileByName(filename string) error

	SaveLoginRecord(record *LoginLog)

	SaveAccessRecord(record *AccessLog)
}

type IOtherService interface {
	UploadImage(fileHeader *multipart.FileHeader) (string, error)

	UploadFile(fileHeader *multipart.FileHeader) (string, error)

	DeleteFile(filename string)

	TraceLogin(record *LoginLog)

	TraceAccess(referee, ip, ua string)
}
