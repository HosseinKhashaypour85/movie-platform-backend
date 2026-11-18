package catalog

import (
	"gorm.io/gorm"
	"time"
)

type Genre struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:80;not null;uniqueIndex" json:"name"`
	Slug      string    `gorm:"size:100;not null;uniqueIndex" json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Movie struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	Title           string         `gorm:"size:240;not null;index" json:"title"`
	OriginalTitle   string         `gorm:"size:240" json:"original_title,omitempty"`
	Slug            string         `gorm:"size:260;not null;uniqueIndex" json:"slug"`
	Synopsis        string         `gorm:"type:text" json:"synopsis,omitempty"`
	ReleaseYear     uint16         `gorm:"index" json:"release_year,omitempty"`
	DurationMinutes uint16         `json:"duration_minutes,omitempty"`
	AgeRating       string         `gorm:"size:20" json:"age_rating,omitempty"`
	PosterURL       string         `gorm:"size:2048" json:"poster_url,omitempty"`
	BackdropURL     string         `gorm:"size:2048" json:"backdrop_url,omitempty"`
	Language        string         `gorm:"size:80" json:"language,omitempty"`
	Country         string         `gorm:"size:100" json:"country,omitempty"`
	IMDbRating      float32        `json:"imdb_rating"`
	IsPublished     bool           `gorm:"not null;default:false;index" json:"is_published"`
	PublishedAt     *time.Time     `json:"published_at,omitempty"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	Genres          []Genre        `gorm:"many2many:movie_genres;" json:"genres,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

type Series struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Title         string         `gorm:"size:240;not null;index" json:"title"`
	OriginalTitle string         `gorm:"size:240" json:"original_title,omitempty"`
	Slug          string         `gorm:"size:260;not null;uniqueIndex" json:"slug"`
	Synopsis      string         `gorm:"type:text" json:"synopsis,omitempty"`
	ReleaseYear   uint16         `gorm:"index" json:"release_year,omitempty"`
	EndYear       uint16         `json:"end_year,omitempty"`
	AgeRating     string         `gorm:"size:20" json:"age_rating,omitempty"`
	PosterURL     string         `gorm:"size:2048" json:"poster_url,omitempty"`
	BackdropURL   string         `gorm:"size:2048" json:"backdrop_url,omitempty"`
	Language      string         `gorm:"size:80" json:"language,omitempty"`
	Country       string         `gorm:"size:100" json:"country,omitempty"`
	IMDbRating    float32        `json:"imdb_rating"`
	IsPublished   bool           `gorm:"not null;default:false;index" json:"is_published"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	Genres        []Genre        `gorm:"many2many:series_genres;" json:"genres,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

type Season struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	SeriesID  uint      `gorm:"not null;uniqueIndex:idx_series_season" json:"series_id"`
	Number    uint16    `gorm:"not null;uniqueIndex:idx_series_season" json:"number"`
	Title     string    `gorm:"size:240" json:"title,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Episode struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	SeasonID        uint       `gorm:"not null;index" json:"season_id"`
	Number          uint16     `gorm:"not null" json:"number"`
	Title           string     `gorm:"size:240;not null" json:"title"`
	Synopsis        string     `gorm:"type:text" json:"synopsis,omitempty"`
	DurationMinutes uint16     `json:"duration_minutes,omitempty"`
	ThumbnailURL    string     `gorm:"size:2048" json:"thumbnail_url,omitempty"`
	ReleaseAt       *time.Time `json:"release_at,omitempty"`
	IsPublished     bool       `gorm:"not null;default:false;index" json:"is_published"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
