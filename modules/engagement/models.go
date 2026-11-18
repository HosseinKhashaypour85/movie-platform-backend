package engagement

import "time"

type Favorite struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;uniqueIndex:idx_user_movie_favorite" json:"user_id"`
	MovieID   uint      `gorm:"not null;uniqueIndex:idx_user_movie_favorite;index" json:"movie_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Rating struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;uniqueIndex:idx_user_movie_rating" json:"user_id"`
	MovieID   uint      `gorm:"not null;uniqueIndex:idx_user_movie_rating;index" json:"movie_id"`
	Score     uint8     `gorm:"not null" json:"score"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
