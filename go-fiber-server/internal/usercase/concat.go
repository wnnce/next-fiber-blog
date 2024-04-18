package usercase

// Concat 联系方式
type Concat struct {
	ConcatId  int64  `json:"concatId,omitempty" db:"concat_id"`   // 联系方式Id
	Name      string `json:"name,omitempty" db:"name"`            // 联系方式名称
	LogoUrl   string `json:"logoUrl,omitempty" db:"logo_url"`     // 联系方式Logo地址
	TargetUrl string `json:"targetUrl,omitempty" db:"target_url"` // 联系方式源链接
	IsMain    string `json:"isMain" db:"is_main"`                 // 是否为主要联系方式
	*BaseEntity
}
