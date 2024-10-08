package usercase

// Category 分类
type Category struct {
	CategoryId   uint   `json:"categoryId,omitempty" db:"category_id"`                                // 分类ID
	CategoryName string `json:"categoryName,omitempty" db:"category_name" validate:"required,max=64"` // 分类名称
	Description  string `json:"description,omitempty" db:"description" validate:"max=255,omitempty"`  // 分类描述
	CoverUrl     string `json:"coverUrl,omitempty" db:"cover_url" validate:"required"`                // 分类封面路径
	ViewNum      uint64 `json:"viewNum" db:"view_num"`                                                // 分类查看次数
	ParentId     uint   `json:"parentId" db:"parent_id"`                                              // 分类上级Id
	IsHot        bool   `json:"isHot" db:"is_hot"`                                                    // 是否热门
	IsTop        bool   `json:"isTop" db:"is_top"`                                                    // 是否置顶
	BaseEntity
}

type CategoryVo struct {
	Category
	ArticleNum int           `json:"articleNum"`         // 分类关联的文章数量
	Children   []*CategoryVo `json:"children,omitempty"` // 子分类
}

// ArticleCategoryVo 博客文章分类Vo类 用于查询文章时一起返回
type ArticleCategoryVo struct {
	CategoryId   uint   `json:"categoryId,omitempty" db:"category_id"`
	CategoryName string `json:"categoryName,omitempty" db:"category_name"`
}

// CategoryUpdateForm 分类快捷更新表单
type CategoryUpdateForm struct {
	CategoryId uint  `json:"categoryId" validate:"required"`
	IsHot      *bool `json:"isHot"`
	IsTop      *bool `json:"isTop"`
	Status     *int  `json:"status"`
}

func (c *CategoryVo) GetId() uint {
	return c.CategoryId
}

func (c *CategoryVo) GetParentId() uint {
	return c.ParentId
}

func (c *CategoryVo) AppendChild(t Tree[uint]) {
	if cat, ok := t.(*CategoryVo); ok {
		if c.Children == nil {
			c.Children = make([]*CategoryVo, 0)
		}
		c.Children = append(c.Children, cat)
	}
}

type ICategoryRepo interface {
	// Save 保存分类
	Save(cat *Category) error
	// Update 更新分类
	Update(cat *Category) error
	// UpdateSelective 更新分类状态
	UpdateSelective(from *CategoryUpdateForm) error
	// UpdateViewNum 更新分类的查看数量
	UpdateViewNum(catId uint, addNum int) error
	// SelectById 通过Id查询分类
	SelectById(catId int) (*Category, error)
	// List 查询分类列表
	List() ([]*CategoryVo, error)
	// ManageList 管理端查询分类
	ManageList() ([]*CategoryVo, error)
	// ListByIds 通过分类Id列表查询分类列表
	ListByIds([]uint) ([]Category, error)

	CountByName(name string, catId uint) (uint8, error)

	DeleteById(catId int) error

	BatchDelete(ids []int) (int64, error)
}

type ICategoryService interface {
	CreateCategory(cat *Category) error

	UpdateCategory(cat *Category) error

	UpdateSelectiveCategory(form *CategoryUpdateForm) error

	TreeCategory() ([]*CategoryVo, error)

	ManageTreeCategory() ([]*CategoryVo, error)

	QueryCategoryInfo(catId int) (*Category, error)

	Delete(catId int) error
}
