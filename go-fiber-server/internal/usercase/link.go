package usercase

import "time"

// Link 友情链接
type Link struct {
	LinkId    uint64 `json:"linkId" db:"link_id"`                                             // 友情链接Id
	Name      string `json:"name,omitempty" db:"name" validate:"required,max=64"`             // 友情链接名称
	Summary   string `json:"summary,omitempty" db:"summary" validate:"max=255,omitempty"`     // 友情链接简介
	CoverUrl  string `json:"coverUrl,omitempty" db:"cover_url" validate:"required,max=255"`   // 友情链接封面地址
	TargetUrl string `json:"targetUrl,omitempty" db:"target_url" validate:"required,max=255"` // 友情链接源地址
	ClickNum  uint64 `json:"clickNum" db:"click_num"`                                         // 友情链接的点击次数
	BaseEntity
}

// LinkQueryForm 友情链接后端查询参数
type LinkQueryForm struct {
	Name            string     `json:"name,omitempty"`
	CreateTimeBegin *time.Time `json:"createTimeBegin,omitempty"`
	CreateTimeEnd   *time.Time `json:"createTimeEnd,omitempty"`
	PageQueryForm
}

type ILinkRepo interface {
	// Save 保存
	Save(link *Link) error
	// Update 更新
	Update(link *Link) error
	// UpdateStatus 更新状态
	UpdateStatus(linkId int64, status uint8) error
	// Page 获取分页
	Page(query *PageQueryForm) ([]*Link, int64, error)
	// ManagePage 管理端获取分页
	ManagePage(query *LinkQueryForm) ([]*Link, int64, error)
	// DeleteById 通过Id删除
	DeleteById(linkId int64) error
	// BatchDelete 批量删除
	BatchDelete(linkIds []int64) (int64, error)
}

type ILinkService interface {
	CreateLike(link *Link) error

	UpdateLike(like *Link) error

	PageLike(query *PageQueryForm) (*PageData[Link], error)

	ManagePageLike(query *LinkQueryForm) (*PageData[Link], error)

	Delete(linkId int64) error
}
