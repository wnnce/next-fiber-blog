package usercase

// Comment 评论
type Comment struct {
	CommentId   int64  `json:"commentId,omitempty" db:"comment_id"` // 评论ID
	Content     string `json:"content,omitempty" db:"content"`      // 评论内容 markdown格式
	UserId      uint64 `json:"userId,omitempty" db:"user_id"`       // 评论的用户Id
	ArticleId   *int64 `json:"articleId,omitempty" db:"article_id"` // 评论对应的文章Id
	TopicId     *int64 `json:"topicId,omitempty" db:"topic_id"`     // 评论对应的动态Id
	Fid         int64  `json:"fid" db:"fid"`                        // 评论的顶级Id
	Rid         int64  `json:"rid" db:"rid"`                        // 评论的上级Id
	Location    string `json:"location,omitempty" db:"location"`    // 发表评论的位置
	CommentIp   string `json:"commentIp,omitempty" db:"comment_ip"` // 发表评论的IP
	CommentUa   string `json:"commentUa,omitempty" db:"comment_ua"` // 发表评论的UA
	VoteUp      int64  `json:"voteUp" db:"vote_up"`                 // 评论的点赞数
	VoteDown    int64  `json:"voteDown" db:"vote_down"`             // 评论的点踩数
	CommentType int    `json:"commentType" db:"comment_type"`       // 评论类型 1：文章评论 2：动态评论 3：关于页面评论 4：友情链接评论
	IsHot       bool   `json:"isHot" db:"is_hot"`                   // 是否热门
	IsTop       bool   `json:"isTop" db:"is_top"`                   // 是否置顶
	IsColl      bool   `json:"isColl" db:"is_coll"`                 // 是否折叠
	BaseEntity
}

type CommentUser struct {
	Nickname string   `json:"nickname" db:"nick_name"`
	Avatar   string   `json:"avatar" db:"avatar"`
	Level    uint8    `json:"level" db:"level"`
	Labels   []string `json:"labels,omitempty" db:"labels"`
	Link     string   `json:"link" db:"link"`
}

// CommentVo 评论博客前端Vo类
type CommentVo struct {
	Comment
	User       *CommentUser         `json:"user"`
	ParentUser *CommentUser         `json:"parentUser"`
	Children   *PageData[CommentVo] `json:"children"`
}

// CommentManageVo 评论管理端Vo类
type CommentManageVo struct {
	Comment
	Username     string  `json:"username" db:"username"`
	ArticleTitle *string `json:"articleTitle" db:"article_title"`
}

type CommentQueryForm struct {
	PageQueryForm
	CommentType     int    `json:"commentType"`     // 评论类型
	ArticleId       int64  `json:"articleId"`       // 博客文章Id
	TopicId         int64  `json:"topicId"`         // 博客动态Id
	Fid             int64  `json:"fid"`             // 顶层评论Id
	CreateTimeBegin string `json:"createTimeBegin"` // 开始时间
	CreateTimeEnd   string `json:"createTimeEnd"`   // 结束时间
}

// CommentUpdateForm 评论快捷更新表单
type CommentUpdateForm struct {
	CommentId int64  `json:"commentId"` // 评论Id
	IsHot     *bool  `json:"isHot"`     // 是否热门
	IsTop     *bool  `json:"isTop"`     // 是否置顶
	IsColl    *bool  `json:"isColl"`    // 是否折叠
	Status    *uint8 `json:"status"`    // 评论状态
}

type ICommentRepo interface {
	// Save 保存评论
	Save(comment *Comment) error
	// Page 分页查询评论
	Page(query *CommentQueryForm) (*PageData[CommentVo], error)
	// TotalComment 查询评论总数
	TotalComment(query *CommentQueryForm) (uint64, error)

	// ManagePage 管理端分页查询评论
	ManagePage(query *CommentQueryForm) (*PageData[CommentManageVo], error)

	UpdateSelective(form *CommentUpdateForm) error

	DeleteById(commentId int64) error

	// CommentVoteUp 评论点赞
	CommentVoteUp(commentId int64, num int) error
}

type ICommentService interface {
	// SaveComment 保存评论
	SaveComment(comment *Comment) error
	// Page 前端页面分页查询评论
	Page(query *CommentQueryForm) (*PageData[CommentVo], error)
	// TotalComment 获取某个文章下的评论总数
	TotalComment(query *CommentQueryForm) (uint64, error)

	UpdateSelectiveComment(form *CommentUpdateForm) error

	Delete(commentId int64) error

	ManagePage(query *CommentQueryForm) (*PageData[CommentManageVo], error)

	CommentVoteUp(commentId int64, userId uint64) error
}
