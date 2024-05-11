package usercase

import "time"

// Tree 数据数据接口 用于工具函数统一处理树形数据
type Tree interface {
	// GetId 获取节点Id
	GetId() int64
	// GetParentId 获取上级节点Id
	GetParentId() int64
	// AppendChild 添加到子节点
	AppendChild(t Tree)
}

// UploadFile 文件上传记录
type UploadFile struct {
	ID         int64      `json:"id,omitempty" db:"id"`                  // 主键ID
	FileMd5    string     `json:"fileMd5,omitempty" db:"file_md5"`       // 文件MD5
	OriginName string     `json:"originName,omitempty" db:"origin_name"` // 文件原始名称
	FileName   string     `json:"fileName,omitempty" db:"file_name"`     // 格式化后的文件名称
	FilePath   string     `json:"filePath,omitempty" db:"file_path"`     // 文件的存储路径
	FileSize   int64      `json:"fileSize,omitempty" db:"file_size"`     // 文件大小
	FileType   string     `json:"fileType,omitempty" db:"file_type"`     // 文件类型
	CreateTime *time.Time `json:"createTime,omitempty" db:"create_time"` // 文件的保存时间
}

// BaseEntity 通用字段
type BaseEntity struct {
	DeleteAt   string     `json:"deleteAt,omitempty" db:"delete_at"`     // 是否删除
	CreateTime *time.Time `json:"createTime,omitempty" db:"create_time"` // 创建时间
	UpdateTime *time.Time `json:"updateTime,omitempty" db:"update_time"` // 更新时间
	Sort       uint       `json:"sort" db:"sort"`                        // 排序
	Status     uint8      `json:"status" db:"status"`                    // 状态
}

// BaseForm 基础表单 包含公共数据
type BaseForm struct {
	Sort   *uint  `json:"sort,omitempty" validate:"required,gte=0"`
	Status *uint8 `json:"status,omitempty" validate:"required,min=0,max=1"`
}

// PageQueryForm 分页查询表单
type PageQueryForm struct {
	Page int `json:"page,omitempty" query:"page" validate:"required,gte=1"`         // 页码
	Size int `json:"size,omitempty" query:"size" validate:"required,gte=1,lte=100"` // 每页记录数
}

// PageData 分页返回数据
type PageData[T any] struct {
	Current int   `json:"current"` // 当前页码
	Pages   int   `json:"pages"`   // 总页数
	Total   int64 `json:"total"`   // 总记录数
	Size    int   `json:"size"`    // 每页记录数
	Records []*T  `json:"records"` // 当前页数据
}
