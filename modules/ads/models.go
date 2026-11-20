package ads

import "time"

type Campaign struct {
	ID uint `gorm:"primaryKey" json:"id"`
	Name string `gorm:"size:160;not null" json:"name"`
	Placement string `gorm:"size:40;not null;index" json:"placement"`
	Format string `gorm:"size:20;not null" json:"format"`
	MediaURL string `gorm:"size:2048;not null" json:"media_url"`
	ClickURL string `gorm:"size:2048" json:"click_url,omitempty"`
	Advertiser string `gorm:"size:160" json:"advertiser,omitempty"`
	IsActive bool `gorm:"not null;default:false;index" json:"is_active"`
	StartsAt time.Time `gorm:"not null;index" json:"starts_at"`
	EndsAt time.Time `gorm:"not null;index" json:"ends_at"`
	ImpressionCount uint64 `gorm:"not null;default:0" json:"impression_count"`
	ClickCount uint64 `gorm:"not null;default:0" json:"click_count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Event struct {
	ID uint64 `gorm:"primaryKey" json:"id"`
	CampaignID uint `gorm:"not null;index" json:"campaign_id"`
	EventType string `gorm:"size:16;not null;index" json:"event_type"`
	Placement string `gorm:"size:40;not null" json:"placement"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`
}
