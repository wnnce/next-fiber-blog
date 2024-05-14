package usercase

// Article 博客文章
type Article struct {
	ArticleId   int64  `json:"articleId,omitempty" db:"article_id"`     // 博客文章Id
	Title       string `json:"title,omitempty" db:"title"`              // 博客文章标题
	Summary     string `json:"summary,omitempty" db:"summary"`          // 博客文章简介
	CoverUrl    string `json:"coverUrl,omitempty" db:"cover_url"`       // 博客文章封面地址
	CategoryIds []uint `json:"categoryIds,omitempty" db:"category_ids"` // 博客文章关联的分类Id
	TagIds      []int  `json:"tagIds,omitempty" db:"tag_ids"`           // 博客文章关联的标签id
	ViewNum     int64  `json:"viewNum" db:"view_num"`                   // 文章的查看次数
	ShareNum    int64  `json:"shareNum" db:"share_num"`                 // 文章的分享次数
	VoteUp      int64  `json:"voteUp" db:"vote_up"`                     // 文章的点赞次数
	ContentUrl  string `json:"contentUrl,omitempty" db:"content_url"`   // 文章正文的地址
	Protocol    string `json:"protocol,omitempty" db:"protocol"`        // 文章的许可协议
	Tips        string `json:"tips,omitempty" db:"tips"`                // 文章底部自定义提示
	Password    string `json:"password,omitempty" db:"password"`        // 文章密码 私密文章需要
	IsHot       bool   `json:"isHot" db:"is_hot"`                       // 是否热门文章
	IsTop       bool   `json:"isTop" db:"is_top"`                       // 文章是否置顶
	IsComment   bool   `json:"isComment" db:"is_comment"`               // 文章是否开启评论
	IsPrivate   bool   `json:"isPrivate" db:"is_private"`               // 是否私密文章
	*BaseEntity
}
