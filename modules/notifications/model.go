package notifications

import "time"

type Notification struct {
	ID uint `gorm:"primaryKey" json:"id"`
	UserID uint `gorm:"not null;index:idx_notifications_user_read" json:"user_id"`
	Title string `gorm:"size:160;not null" json:"title"`
	Body string `gorm:"type:text;not null" json:"body"`
	Type string `gorm:"size:40;not null;default:general" json:"type"`
	ActionURL string `gorm:"size:2048" json:"action_url,omitempty"`
	ReadAt *time.Time `gorm:"index:idx_notifications_user_read" json:"read_at,omitempty"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`
}
