package usercase

// Tag 博客标签
type Tag struct {
	TagId    int    `json:"tagId,omitempty" db:"tag_id"`       // 标签ID
	TagName  string `json:"tagName,omitempty" db:"tag_name"`   // 标签名称
	CoverUrl string `json:"coverUrl,omitempty" db:"cover_url"` // 标签封面地址
	ViewNum  int64  `json:"viewNum" db:"view_num"`             // 标签的查看次数
	Color    string `json:"color,omitempty" db:"color"`        // 标签颜色
	*BaseEntity
}
