package usercase

import "time"

// Tag 博客标签
type Tag struct {
	TagId    uint   `json:"tagId,omitempty" db:"tag_id"`       // 标签ID
	TagName  string `json:"tagName,omitempty" db:"tag_name"`   // 标签名称
	CoverUrl string `json:"coverUrl,omitempty" db:"cover_url"` // 标签封面地址
	ViewNum  int64  `json:"viewNum" db:"view_num"`             // 标签的查看次数
	Color    string `json:"color,omitempty" db:"color"`        // 标签颜色
	BaseEntity
}

// TagForm 标签的新增、修改表单
type TagForm struct {
	TagId    uint   `json:"tagId,omitempty"`
	TagName  string `json:"tagName,omitempty" validate:"required,min=1,max=64"`
	CoverUrl string `json:"coverUrl,omitempty" validate:"required"`
	Color    string `json:"color,omitempty" validate:"required,len=6"`
	CommonField
}

// TagQueryForm 标签查询表单
type TagQueryForm struct {
	TagName         string     `json:"tagName,omitempty"`
	CreateTimeBegin *time.Time `json:"createTimeBegin,omitempty"`
	CreateTimeEnd   *time.Time `json:"createTimeEnd,omitempty"`
}

// ITagRepo 标签Repo层接口
type ITagRepo interface {
	// Save 保存标签
	Save(form *TagForm) error
	// Update 更新标签
	Update(form *TagForm) error
	// UpdateStatus 更新标签状态
	UpdateStatus(tagId int, status uint8) error
	// UpdateViewNum 更新标签的查看次数
	UpdateViewNum(tagId int, addNum int) error
	// SelectById 通过Id获取标签数据
	SelectById(id int) (*Tag, error)
	// ManageList 后台获取标签列表
	ManageList(form *TagQueryForm) ([]*Tag, error)
	// List 博客页码获取标签列表
	List() ([]*Tag, error)
	// ListByIds 通过标签Id列表获取标签列表
	ListByIds(ids []uint) ([]*Tag, error)
	// CountByTagName 通过名称查询标签数量
	CountByTagName(name string, tagId uint) (uint8, error)
	// DeleteById 通过Id删除标签
	DeleteById(id int) error
	// DeleteByIds 通过Id列表批量删除
	DeleteByIds(ids []int) (int64, error)
}

// ITagService 标签Service层接口
type ITagService interface {
	// CreateTag 新增标签
	CreateTag(form *TagForm) error
	// UpdateTag 更新标签
	UpdateTag(form *TagForm) error
	// QueryTagInfo 查询标签详情
	QueryTagInfo(id int) (*Tag, error)
	// ManageListTag 查询标签列表
	ManageListTag(form *TagQueryForm) ([]*Tag, error)
	// AllTag 博客获取所有标签
	AllTag() []*Tag
	// Delete 删除单个标签
	Delete(id int) error
	// BatchDelete 批量删除标签
	BatchDelete(ids string) error
}
