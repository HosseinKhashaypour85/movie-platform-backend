
package watch_history

import (
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

// CreateOrUpdate آخرین وضعیت تماشای کاربر را ذخیره می‌کند
func (r *Repository) CreateOrUpdate(history *WatchHistory) error {
	var existing WatchHistory

	err := r.DB.
		Where("user_id = ? AND movie_id = ?", history.UserID, history.MovieID).
		First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		return r.DB.Create(history).Error
	}

	if err != nil {
		return err
	}

	return r.DB.Model(&existing).Updates(map[string]interface{}{
		"progress_seconds": history.ProgressSeconds,
		"duration_seconds": history.DurationSeconds,
		"completed":        history.Completed,
		"watched_at":       history.WatchedAt,
	}).Error
}

// FindByUserID تاریخچه تماشای یک کاربر را دریافت می‌کند
func (r *Repository) FindByUserID(userID uint) ([]WatchHistory, error) {
	var histories []WatchHistory

	err := r.DB.
		Where("user_id = ?", userID).
		Order("watched_at DESC").
		Find(&histories).Error

	return histories, err
}

// FindByUserAndMovie یک فیلم مشخص از تاریخچه کاربر را پیدا می‌کند
func (r *Repository) FindByUserAndMovie(userID uint, movieID uint) (*WatchHistory, error) {
	var history WatchHistory

	err := r.DB.
		Where("user_id = ? AND movie_id = ?", userID, movieID).
		First(&history).Error

	if err != nil {
		return nil, err
	}

	return &history, nil
}

// Delete حذف یک رکورد از تاریخچه
func (r *Repository) Delete(userID uint, movieID uint) error {
	return r.DB.
		Where("user_id = ? AND movie_id = ?", userID, movieID).
		Delete(&WatchHistory{}).Error
}
