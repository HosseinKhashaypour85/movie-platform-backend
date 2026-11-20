package media

import "time"

type TranscodeJob struct {
	ID uint `gorm:"primaryKey" json:"id"`
	SourceAssetID uint `gorm:"not null;uniqueIndex" json:"source_asset_id"`
	Status string `gorm:"size:20;not null;index" json:"status"`
	Attempts uint8 `gorm:"not null;default:0" json:"attempts"`
	ErrorMessage string `gorm:"size:1000" json:"error_message,omitempty"`
	LockedAt *time.Time `json:"locked_at,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
