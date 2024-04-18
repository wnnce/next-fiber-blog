package usercase

import "github.com/lib/pq"

// Topic 博客动态
type Topic struct {
	TopicId   int64          `json:"topicId,omitempty" db:"topic_id"`     // 博客动态ID
	Title     string         `json:"title,omitempty" db:"title"`          // 博客动态标题
	Content   string         `json:"content,omitempty" db:"content"`      // 博客动态内容 markdown格式
	CoverUrls pq.StringArray `json:"coverUrls,omitempty" db:"cover_urls"` // 动态封面图片数组
	Location  string         `json:"location,omitempty" db:"location"`    // 发布动态的位置
	IsHot     bool           `json:"isHot" db:"is_hot"`                   // 是否热门动态
	IsTop     bool           `json:"isTop" db:"is_top"`                   // 动态是否置顶
	*BaseEntity
}
