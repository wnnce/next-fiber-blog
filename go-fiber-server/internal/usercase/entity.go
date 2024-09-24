package usercase

import (
	"math"
	"time"
)

// Tree 数据数据接口 用于工具函数统一处理树形数据
type Tree[T any] interface {
	// GetId 获取节点Id
	GetId() T
	// GetParentId 获取上级节点Id
	GetParentId() T
	// AppendChild 添加到子节点
	AppendChild(t Tree[T])
}

// BaseEntity 通用字段
type BaseEntity struct {
	DeleteAt   int64      `json:"deleteAt,omitempty" db:"delete_at"`     // 是否删除
	CreateTime *time.Time `json:"createTime,omitempty" db:"create_time"` // 创建时间
	UpdateTime *time.Time `json:"updateTime,omitempty" db:"update_time"` // 更新时间
	CommonField
}

// CommonField 通用字段 包含排序和状态
type CommonField struct {
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

// NewPageData 创建分页数据返回对象
func NewPageData[T any](records []*T, total int64, current, size int) *PageData[T] {
	var pages int
	if total > 0 {
		pages = int(math.Ceil(float64(total) / float64(size)))
	}
	// 如果当前页大于总页数 但是 记录仍然有值的话
	// 说明数据库查询时启用了offset安全检查 无论页码超过多大 始终会返回最后一页的数据
	if len(records) > 0 && current > pages {
		current = pages
	}
	return &PageData[T]{
		Current: current,
		Pages:   pages,
		Total:   total,
		Size:    size,
		Records: records,
	}
}
