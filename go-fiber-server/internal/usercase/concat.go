package usercase

// Concat 联系方式
type Concat struct {
	ConcatId  uint   `json:"concatId,omitempty" db:"concat_id"`                               // 联系方式Id
	Name      string `json:"name,omitempty" db:"name" validate:"required,max=64"`             // 联系方式名称
	LogoUrl   string `json:"logoUrl,omitempty" db:"logo_url" validate:"required,max=255"`     // 联系方式Logo地址
	TargetUrl string `json:"targetUrl,omitempty" db:"target_url" validate:"required,max=255"` // 联系方式源链接
	IsMain    bool   `json:"isMain" db:"is_main"`                                             // 是否为主要联系方式
	BaseEntity
}

// ConcatQueryForm 联系方式查询表单
type ConcatQueryForm struct {
	Name            string `json:"name"`
	CreateTimeBegin string `json:"createTimeBegin"`
	CreateTimeEnd   string `json:"createTimeEnd"`
}

type IConcatRepo interface {
	Save(concat *Concat) error

	Update(concat *Concat) error

	UpdateStatus(cid int, status uint) error

	List() ([]*Concat, error)

	ManageList(query *ConcatQueryForm) ([]*Concat, error)

	CountByName(name string, cid uint) (uint8, error)

	DeleteById(cid int) error
}

type IConcatService interface {
	CreateConcat(concat *Concat) error

	UpdateConcat(concat *Concat) error

	ListConcat() ([]*Concat, error)

	ManageListConcat(query *ConcatQueryForm) ([]*Concat, error)

	Delete(cid int) error
}
