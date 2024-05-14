package usercase

// Notice 通知
type Notice struct {
	NoticeId   uint64 `json:"noticeId" db:"notice_id"`                                               // 通知ID
	Title      string `json:"title,omitempty" db:"title" validate:"required"`                        // 通知标题
	Message    string `json:"message,omitempty" db:"message" validate:"required"`                    // 通知内容
	Level      int    `json:"level,omitempty" db:"level" validate:"required,gte=1,lte=4"`            // 通知级别 1:普通 2:警告 3:严重
	NoticeType int    `json:"noticeType,omitempty" db:"notice_type" validate:"requires,gte=1,lte=3"` // 通知类型 1：首页弹窗通知 2：公告板通知 3：后台通知
	BaseEntity
}

// NotifyQueryForm 通知管理后台查询参数
type NotifyQueryForm struct {
	Title      string `json:"title"`
	Level      int    `json:"level"`
	NoticeType int    `json:"noticeType"`
	PageQueryForm
}

type NoticeRepo interface {
	Save(notice *Notice) error

	Update(notice *Notice) error

	ListByType(t int) ([]*Notice, error)

	ManageList(query *NotifyQueryForm) ([]*Notice, int64, error)
}
