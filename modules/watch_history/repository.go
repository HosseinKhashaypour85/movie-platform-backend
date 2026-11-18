package watch_history

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type Repository struct{ DB *gorm.DB }

func NewRepository(db *gorm.DB) *Repository { return &Repository{DB: db} }
func (r *Repository) CreateOrUpdate(h *WatchHistory) error {
	if h == nil {
		return errors.New("history is required")
	}
	if h.WatchedAt.IsZero() {
		h.WatchedAt = time.Now().UTC()
	}
	return r.DB.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "user_id"}, {Name: "movie_id"}}, DoUpdates: clause.AssignmentColumns([]string{"progress_seconds", "duration_seconds", "completed", "watched_at", "updated_at"})}).Create(h).Error
}
func (r *Repository) FindByUserID(uid uint, limit, offset int) ([]WatchHistory, error) {
	if limit < 1 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	var rows []WatchHistory
	e := r.DB.Where("user_id = ?", uid).Order("watched_at DESC").Limit(limit).Offset(offset).Find(&rows).Error
	return rows, e
}
func (r *Repository) FindByUserAndMovie(uid, mid uint) (*WatchHistory, error) {
	var h WatchHistory
	if e := r.DB.Where("user_id = ? AND movie_id = ?", uid, mid).First(&h).Error; e != nil {
		return nil, e
	}
	return &h, nil
}
func (r *Repository) Delete(uid, mid uint) error {
	return r.DB.Where("user_id = ? AND movie_id = ?", uid, mid).Delete(&WatchHistory{}).Error
}
