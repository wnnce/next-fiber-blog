package usercase

// Comment 评论
type Comment struct {
	CommentId   int64  `json:"commentId,omitempty" db:"comment_id"` // 评论ID
	Content     string `json:"content,omitempty" db:"content"`      // 评论内容 markdown格式
	UserId      int64  `json:"userId,omitempty" db:"user_id"`       // 评论的用户Id
	ArticleId   int64  `json:"articleId,omitempty" db:"article_id"` // 评论对应的文章Id
	TopicId     int64  `json:"topicId,omitempty" db:"topic_id"`     // 评论对应的动态Id
	Fid         int64  `json:"fid,omitempty" db:"fid"`              // 评论的顶级Id
	Rid         int64  `json:"rid,omitempty" db:"rid"`              // 评论的上级Id
	Location    string `json:"location,omitempty" db:"location"`    // 发表评论的位置
	CommentIp   string `json:"commentIp,omitempty" db:"comment_ip"` // 发表评论的IP
	CommentUa   string `json:"commentUa,omitempty" db:"comment_ua"` // 发表评论的UA
	VoteUp      int64  `json:"voteUp" db:"vote_up"`                 // 评论的点赞数
	VoteDown    int64  `json:"voteDown" db:"vote_down"`             // 评论的点踩数
	CommentType int    `json:"commentType" db:"comment_type"`       // 评论类型 1：文章评论 2：动态评论
	IsHot       bool   `json:"isHot" db:"is_hot"`                   // 是否热门
	IsTop       bool   `json:"isTop" db:"is_top"`                   // 是否置顶
	IsColl      bool   `json:"isColl" db:"is_coll"`                 // 是否折叠
	*BaseEntity
}
