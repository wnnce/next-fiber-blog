package usercase

import "time"

// Notice 通知
type Notice struct {
	NoticeId   int64      `json:"noticeId" db:"notice_id"`               // 通知ID
	Title      string     `json:"title,omitempty" db:"title"`            // 通知标题
	Message    string     `json:"message,omitempty" db:"message"`        // 通知内容
	Level      int        `json:"level,omitempty" db:"level"`            // 通知级别
	NoticeType int        `json:"noticeType,omitempty" db:"notice_type"` // 通知类型 1：首页弹窗通知 2：公告板通知 3：后台通知
	CreateTime *time.Time `json:"createTime,omitempty" db:"create_time"` // 创建时间
	UpdateTime *time.Time `json:"updateTime,omitempty" db:"update_time"` // 更新时间
	Sort       int        `json:"sort" db:"sort"`                        // 排序
	Status     int        `json:"status" db:"status"`                    // 状态
}
