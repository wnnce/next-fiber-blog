package usercase

import "time"

// Notice 通知
type Notice struct {
	NoticeId   int64      `json:"noticeId" db:"notice_id"`                                               // 通知ID
	Title      string     `json:"title,omitempty" db:"title" validate:"required"`                        // 通知标题
	Message    string     `json:"message,omitempty" db:"message" validate:"required"`                    // 通知内容
	Level      int        `json:"level,omitempty" db:"level" validate:"required,gte=1,lte=4"`            // 通知级别
	NoticeType int        `json:"noticeType,omitempty" db:"notice_type" validate:"requires,gte=1,lte=3"` // 通知类型 1：首页弹窗通知 2：公告板通知 3：后台通知
	CreateTime *time.Time `json:"createTime,omitempty" db:"create_time"`                                 // 创建时间
	UpdateTime *time.Time `json:"updateTime,omitempty" db:"update_time"`                                 // 更新时间
	*CommonField
}
