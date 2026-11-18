package watch_history

import "time"

type WatchHistory struct {
	ID uint `gorm:"primaryKey" json:"id"`

	UserID uint `gorm:"uniqueIndex:idx_user_movie" json:"user_id"`

	MovieID uint `gorm:"uniqueIndex:idx_user_movie" json:"movie_id"`

	ProgressSeconds int64 `gorm:"default:0" json:"progress_seconds"`

	DurationSeconds int64 `gorm:"default:0" json:"duration_seconds"`

	Completed bool `gorm:"default:false" json:"completed"`

	WatchedAt time.Time `json:"watched_at"`

	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`
}
