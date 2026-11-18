package catalog

import (
	"strings"
	"gorm.io/gorm"
)

type Repository struct { DB *gorm.DB }
func NewRepository(db *gorm.DB) *Repository { return &Repository{DB: db} }

func (r *Repository) ListMovies(query string, genreID uint, year uint16, limit, offset int) ([]Movie, int64, error) {
	if limit < 1 || limit > 100 { limit = 20 }; if offset < 0 { offset = 0 }
	tx := r.DB.Model(&Movie{}).Where("is_published = ?", true)
	if query = strings.TrimSpace(query); query != "" { like := "%" + query + "%"; tx = tx.Where("title LIKE ? OR original_title LIKE ?", like, like) }
	if year > 0 { tx = tx.Where("release_year = ?", year) }
	if genreID > 0 { tx = tx.Joins("JOIN movie_genres ON movie_genres.movie_id = movies.id").Where("movie_genres.genre_id = ?", genreID) }
	var total int64; if err := tx.Distinct("movies.id").Count(&total).Error; err != nil { return nil, 0, err }
	var rows []Movie; err := tx.Preload("Genres").Distinct("movies.*").Order("created_at DESC").Limit(limit).Offset(offset).Find(&rows).Error
	return rows, total, err
}
func (r *Repository) MovieByID(id uint) (*Movie, error) { var row Movie; if err := r.DB.Preload("Genres").Where("is_published = ?", true).First(&row, id).Error; err != nil { return nil, err }; return &row, nil }
func (r *Repository) ListSeries(query string, limit, offset int) ([]Series, int64, error) {
	if limit < 1 || limit > 100 { limit = 20 }; if offset < 0 { offset = 0 }
	tx := r.DB.Model(&Series{}).Where("is_published = ?", true); if query = strings.TrimSpace(query); query != "" { like := "%" + query + "%"; tx = tx.Where("title LIKE ? OR original_title LIKE ?", like, like) }
	var total int64; if err := tx.Count(&total).Error; err != nil { return nil, 0, err }
	var rows []Series; err := tx.Preload("Genres").Order("created_at DESC").Limit(limit).Offset(offset).Find(&rows).Error
	return rows, total, err
}
func (r *Repository) SeriesByID(id uint) (*Series, error) { var row Series; if err := r.DB.Preload("Genres").Where("is_published = ?", true).First(&row, id).Error; err != nil { return nil, err }; return &row, nil }
func (r *Repository) Seasons(seriesID uint) ([]Season, error) { var rows []Season; err := r.DB.Where("series_id = ?", seriesID).Order("number ASC").Find(&rows).Error; return rows, err }
func (r *Repository) Episodes(seasonID uint) ([]Episode, error) { var rows []Episode; err := r.DB.Where("season_id = ? AND is_published = ?", seasonID, true).Order("number ASC").Find(&rows).Error; return rows, err }
