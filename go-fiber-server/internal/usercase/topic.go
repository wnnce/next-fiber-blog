package usercase

// Topic 博客动态
type Topic struct {
	TopicId   uint64   `json:"topicId,omitempty" db:"topic_id"`                    // 动态ID
	Content   string   `json:"content,omitempty" db:"content" validate:"required"` // 动态内容 markdown格式
	ImageUrls []string `json:"imageUrls,omitempty" db:"image_urls"`                // 动态图片数组
	Location  string   `json:"location,omitempty" db:"location"`                   // 发布动态的位置
	IsHot     bool     `json:"isHot" db:"is_hot"`                                  // 是否热门动态
	IsTop     bool     `json:"isTop" db:"is_top"`                                  // 动态是否置顶
	VoteUp    int      `json:"voteUp" db:"vote_up"`                                // 动态点赞数
	Mode      int      `json:"mode" db:"mode" validate:"required,min=1"`           // 动态模式 1:图文动态 2:照片墙
	BaseEntity
}

// TopicUpdateForm 动态快捷更新表单
type TopicUpdateForm struct {
	TopicId uint64 `json:"topicId" validate:"required"`
	IsHot   *bool  `json:"isHot"`
	IsTop   *bool  `json:"isTop"`
	Status  *uint8 `json:"status"`
}

// TopicQueryForm 动态查询表单
type TopicQueryForm struct {
	Location        string `json:"location"`
	Status          *uint8 `json:"status"`
	CreateTimeBegin string `json:"createTimeBegin"`
	CreateTimeEnd   string `json:"createTimeEnd"`
	PageQueryForm
}

// ITopicRepo 动态持久层接口
type ITopicRepo interface {
	Save(topic *Topic) error

	Update(topic *Topic) error

	UpdateSelective(form *TopicUpdateForm) error

	Page(query *TopicQueryForm) ([]*Topic, int64, error)

	DeleteById(topicId int64) error
}

// ITopicService 动态Service层接口
type ITopicService interface {
	SaveTopic(topic *Topic) error

	UpdateTopic(topic *Topic) error

	UpdateSelectiveTopic(form *TopicUpdateForm) error

	PageTopic(query *TopicQueryForm) (*PageData[Topic], error)

	DeleteTopicById(topicId int64) error
}
