package usercase

// Article 博客文章
type Article struct {
	ArticleId   uint64 `json:"articleId,omitempty" db:"article_id"`                         // 博客文章Id
	Title       string `json:"title,omitempty" db:"title" validate:"required"`              // 博客文章标题
	Summary     string `json:"summary,omitempty" db:"summary" validate:"required"`          // 博客文章简介
	CoverUrl    string `json:"coverUrl,omitempty" db:"cover_url" validate:"required"`       // 博客文章封面地址
	CategoryIds []uint `json:"categoryIds,omitempty" db:"category_ids" validate:"required"` // 博客文章关联的分类Id
	TagIds      []uint `json:"tagIds,omitempty" db:"tag_ids" validate:"required"`           // 博客文章关联的标签id
	ViewNum     int64  `json:"viewNum" db:"view_num"`                                       // 文章的查看次数
	ShareNum    int64  `json:"shareNum" db:"share_num"`                                     // 文章的分享次数
	VoteUp      int64  `json:"voteUp" db:"vote_up"`                                         // 文章的点赞次数
	Content     string `json:"content,omitempty" db:"content"`                              // 文章正文 数据库会经过gzip压缩
	Protocol    string `json:"protocol,omitempty" db:"protocol"`                            // 文章的许可协议
	Tips        string `json:"tips,omitempty" db:"tips"`                                    // 文章底部自定义提示
	Password    string `json:"password,omitempty" db:"password"`                            // 文章密码 私密文章需要
	IsHot       bool   `json:"isHot" db:"is_hot"`                                           // 是否热门文章
	IsTop       bool   `json:"isTop" db:"is_top"`                                           // 文章是否置顶
	IsComment   bool   `json:"isComment" db:"is_comment"`                                   // 文章是否开启评论
	IsPrivate   bool   `json:"isPrivate" db:"is_private"`                                   // 是否私密文章
	*BaseEntity
}

// ArticleUpdateForm 文章更新表单
type ArticleUpdateForm struct {
	ArticleId uint64 `json:"articleId" validate:"required"`
	IsHot     *bool  `json:"isHot"`
	IsTop     *bool  `json:"isTop"`
	IsComment *bool  `json:"isComment"`
	Status    *uint8 `json:"status"`
}

// ArticleQueryForm 文章查询表单
type ArticleQueryForm struct {
	Title           string `json:"title"`
	TagId           *uint  `json:"tagId"`
	CategoryId      *uint  `json:"categoryId"`
	Status          *uint8 `json:"status"`
	CreateTimeBegin string `json:"createTimeBegin"`
	CreateTimeEnd   string `json:"createTimeEnd"`
	PageQueryForm
}

type IArticleRepo interface {
	Save(article *Article) error

	Update(article *Article) error

	UpdateSelective(form *ArticleUpdateForm) error

	Page(query *ArticleQueryForm) ([]*Article, int64, error)

	SelectById(articleId int64, checkStatus bool) (*Article, error)

	CountByTagId(tagId int) (int64, error)

	CountByCategoryId(categoryId int) (int64, error)

	DeleteById(articleId int64) error
}
