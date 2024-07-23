package usercase

// Notice 通知
type Notice struct {
	NoticeId   uint64 `json:"noticeId" db:"notice_id"`                                               // 通知ID
	Title      string `json:"title,omitempty" db:"title" validate:"required"`                        // 通知标题
	Message    string `json:"message,omitempty" db:"message" validate:"required"`                    // 通知内容
	Level      int    `json:"level,omitempty" db:"level" validate:"required,gte=1,lte=4"`            // 通知级别 1:普通 2:警告 3:严重
	NoticeType int    `json:"noticeType,omitempty" db:"notice_type" validate:"required,gte=1,lte=3"` // 通知类型 1：首页弹窗通知 2：公告板通知 3：后台通知
	BaseEntity
}

// NoticeQueryForm 通知管理后台查询参数
type NoticeQueryForm struct {
	Title      string `json:"title"`
	Level      *int   `json:"level"`
	NoticeType *int   `json:"noticeType"`
	PageQueryForm
}

type INoticeRepo interface {
	Save(notice *Notice) error

	Update(notice *Notice) error

	ListByType(t int) ([]Notice, error)

	ManagePage(query *NoticeQueryForm) ([]*Notice, int64, error)

	QueryNoticeTypeById(noticeId int64) int

	DeleteById(id int64) error
}

type INoticeService interface {
	SaveNotice(notice *Notice) error

	UpdateNotice(notice *Notice) error

	Page(query *NoticeQueryForm) (*PageData[Notice], error)

	ListNoticeByType(noticeType int) ([]Notice, error)

	Delete(noticeId int64) error
}
