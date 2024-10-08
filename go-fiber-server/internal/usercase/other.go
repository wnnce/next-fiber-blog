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

// LoginLogQueryForm 登录日志查询表单
type LoginLogQueryForm struct {
	Username        string `json:"username"`        // 用户名
	LoginType       *uint8 `json:"loginType"`       // 扽路过类型
	Result          *int   `json:"result"`          // 登录结果
	CreateTimeBegin string `json:"createTimeBegin"` // 开始时间
	CreateTimeEnd   string `json:"createTimeEnd"`   // 结束时间
	PageQueryForm
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

// AccessLogQueryForm 访问日志查询表单
type AccessLogQueryForm struct {
	Ip              string `json:"ip"`
	CreateTimeBegin string `json:"createTimeBegin"`
	CreateTimeEnd   string `json:"createTimeEnd"`
	PageQueryForm
}

type SiteConfigurationItem struct {
	Name  string `json:"name"`  // 配置项名称
	Type  string `json:"type"`  // 配置项类型 image:图片 text:普通字符串 number:数字 markdown:markdown格式文本 html:html格式文本 color:颜色
	Value any    `json:"value"` // 配置项内容
}

// DefaultSiteConfiguration 默认站点配置
var DefaultSiteConfiguration = map[string]SiteConfigurationItem{
	"tabTitle":     {Name: "浏览器标题", Value: "demo", Type: "text"},
	"logo":         {Name: "Logo", Type: "image", Value: "demo"},
	"avatar":       {Name: "博客头像", Type: "image", Value: "demo"},
	"title":        {Name: "博客标题", Type: "text", Value: "demo"},
	"summary":      {Name: "博客简介", Type: "text", Value: "demo"},
	"about":        {Name: "关于", Type: "markdown", Value: "demo"},
	"powered":      {Name: "底部Powered", Type: "html", Value: "demo"},
	"icp":          {Name: "ICP备案", Type: "html", Value: ""},
	"articleSize":  {Name: "文章数量", Type: "number", Value: 5},
	"topicSize":    {Name: "动态数量", Type: "number", Value: 10},
	"commentSize":  {Name: "评论数量", Type: "number", Value: 5},
	"primaryColor": {Name: "主要颜色", Type: "color", Value: "#ffffff"},
}

// SiteStats 站点信息统计数据
type SiteStats struct {
	ArticleCount  uint   `json:"articleCount" db:"article_count"`
	CategoryCount uint   `json:"categoryCount" db:"category_count"`
	TagCount      uint   `json:"tagCount" db:"tag_count"`
	CommentCount  uint   `json:"commentCount" db:"comment_count"`
	VisitorCount  uint   `json:"visitorCount" db:"visitor_count"`
	AccessCount   uint64 `json:"accessCount" db:"access_count"`
	WordTotal     uint64 `json:"wordTotal" db:"word_total"`
}

type IOtherRepo interface {
	// SaveFileRecord 保存文件上传记录
	SaveFileRecord(file *UploadFile)

	// QueryFileByMd5 通过文件md5查询文件
	QueryFileByMd5(fileMd5 string) (*UploadFile, error)

	// DeleteFileByName 通过文件名称删除文件
	DeleteFileByName(filename string) error

	// SaveLoginRecord 保存登录记录
	SaveLoginRecord(record *LoginLog)

	// SaveAccessRecord 保存访问记录
	SaveAccessRecord(record *AccessLog)

	// PageLoginRecord 分页查询登录记录
	PageLoginRecord(query *LoginLogQueryForm) ([]*LoginLog, int64, error)

	// PageAccessRecord 分页查询访问记录
	PageAccessRecord(query *AccessLogQueryForm) ([]*AccessLog, int64, error)

	// SiteStats 获取站点的统计信息
	SiteStats() (SiteStats, error)
}

type IOtherService interface {

	// SiteConfiguration 获取站点配置
	SiteConfiguration() map[string]SiteConfigurationItem

	// UpdateSiteConfiguration 更新站点配置
	UpdateSiteConfiguration(config map[string]SiteConfigurationItem) error

	// UploadImage 上传图片
	UploadImage(fileHeader *multipart.FileHeader) (string, error)

	// UploadFile 上传文件
	UploadFile(fileHeader *multipart.FileHeader) (string, error)

	// DeleteFile 删除上传文件
	DeleteFile(filename string)

	// TraceLogin 记录登录日志
	TraceLogin(record *LoginLog)

	// TraceAccess 记录访问日志
	TraceAccess(referee, ip, ua string)

	// PageLogin 分页查询登录日志
	PageLogin(query *LoginLogQueryForm) (*PageData[LoginLog], error)

	// PageAccess 分页查询访问日志
	PageAccess(query *AccessLogQueryForm) (*PageData[AccessLog], error)

	// SiteStats 获取站点的统计信息
	SiteStats() (SiteStats, error)
}
