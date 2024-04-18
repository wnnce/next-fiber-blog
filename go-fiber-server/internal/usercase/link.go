package usercase

// Link 友情链接
type Link struct {
	LinkId    int64  `json:"linkId" db:"link_id"`                 // 友情链接Id
	Name      string `json:"name,omitempty" db:"name"`            // 友情链接名称
	Summary   string `json:"summary,omitempty" db:"summary"`      // 友情链接简介
	CoverUrl  string `json:"coverUrl,omitempty" db:"cover_url"`   // 友情链接封面地址
	TargetUrl string `json:"targetUrl,omitempty" db:"target_url"` // 友情链接源地址
	ClickNum  int64  `json:"clickNum" db:"click_num"`             // 友情链接的点击次数
	*BaseEntity
}
