package media

import "time"

type Asset struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	MovieID         *uint     `gorm:"index" json:"movie_id,omitempty"`
	EpisodeID       *uint     `gorm:"index" json:"episode_id,omitempty"`
	Kind            string    `gorm:"size:20;not null;index" json:"kind"`
	StorageKey      string    `gorm:"size:1024;not null" json:"-"`
	PlaybackURL     string    `gorm:"size:2048" json:"playback_url,omitempty"`
	MimeType        string    `gorm:"size:120" json:"mime_type,omitempty"`
	SizeBytes       uint64    `json:"size_bytes"`
	DurationSeconds uint32    `json:"duration_seconds"`
	IsActive        bool      `gorm:"not null;default:false;index" json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Subtitle struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	MovieID     *uint     `gorm:"index" json:"movie_id,omitempty"`
	EpisodeID   *uint     `gorm:"index" json:"episode_id,omitempty"`
	Language    string    `gorm:"size:12;not null;index" json:"language"`
	Label       string    `gorm:"size:80;not null" json:"label"`
	Format      string    `gorm:"size:12;not null" json:"format"`
	StorageKey  string    `gorm:"size:1024;not null" json:"-"`
	DownloadURL string    `gorm:"size:2048" json:"download_url,omitempty"`
	IsGenerated bool      `gorm:"not null;default:false" json:"is_generated"`
	IsPublished bool      `gorm:"not null;default:false;index" json:"is_published"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SubtitleJob struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	MovieID          *uint      `gorm:"index" json:"movie_id,omitempty"`
	EpisodeID        *uint      `gorm:"index" json:"episode_id,omitempty"`
	SourceAssetID    uint       `gorm:"not null;index" json:"source_asset_id"`
	SourceLanguage   string     `gorm:"size:12;not null" json:"source_language"`
	TargetLanguage   string     `gorm:"size:12;not null" json:"target_language"`
	Status           string     `gorm:"size:20;not null;index" json:"status"`
	Attempts         uint8      `gorm:"not null;default:0" json:"attempts"`
	ErrorMessage     string     `gorm:"size:1000" json:"error_message,omitempty"`
	OutputSubtitleID *uint      `json:"output_subtitle_id,omitempty"`
	LockedAt         *time.Time `json:"locked_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}
