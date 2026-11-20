package watch_history

import "time"

type EpisodeWatchHistory struct {
	ID uint `gorm:"primaryKey" json:"id"`
	UserID uint `gorm:"not null;uniqueIndex:idx_user_episode" json:"user_id"`
	EpisodeID uint `gorm:"not null;uniqueIndex:idx_user_episode" json:"episode_id"`
	ProgressSeconds int64 `gorm:"not null;default:0" json:"progress_seconds"`
	DurationSeconds int64 `gorm:"not null;default:0" json:"duration_seconds"`
	Completed bool `gorm:"not null;default:false" json:"completed"`
	WatchedAt time.Time `json:"watched_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
