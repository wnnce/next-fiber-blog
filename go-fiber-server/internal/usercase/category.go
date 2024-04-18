package usercase

type Category struct {
	CategoryId   int    `json:"categoryId,omitempty" db:"category_id"`     // 分类ID
	CategoryName string `json:"categoryName,omitempty" db:"category_name"` // 分类名称
	Description  string `json:"description,omitempty" db:"description"`    // 分类描述
	CoverUrl     string `json:"coverUrl,omitempty" db:"cover_url"`         // 分类封面路径
	ViewNum      int64  `json:"viewNum" db:"view_num"`                     // 分类查看次数
	PatentId     int    `json:"patentId" db:"parent_id"`                   // 分类上级Id
	IsHot        bool   `json:"isHot" db:"is_hot"`                         // 是否热门
	IsTop        bool   `json:"isTop" db:"is_top"`                         // 是否置顶
	*BaseEntity
}
