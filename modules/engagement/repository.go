package engagement

import (
 "gorm.io/gorm"
 "gorm.io/gorm/clause"
)

type Repository struct { db *gorm.DB }
func NewRepository(db *gorm.DB) *Repository { return &Repository{db:db} }

func (r *Repository) AddFavorite(userID,movieID uint) error {
 return r.db.Clauses(clause.OnConflict{Columns:[]clause.Column{{Name:"user_id"},{Name:"movie_id"}},DoNothing:true}).Create(&Favorite{UserID:userID,MovieID:movieID}).Error
}
func (r *Repository) RemoveFavorite(userID,movieID uint) error {
 return r.db.Where("user_id = ? AND movie_id = ?",userID,movieID).Delete(&Favorite{}).Error
}
func (r *Repository) Favorites(userID uint,limit,offset int)([]Favorite,error) {
 var rows []Favorite
 err:=r.db.Where("user_id = ?",userID).Order("created_at DESC").Limit(limit).Offset(offset).Find(&rows).Error
 return rows,err
}
func (r *Repository) Rate(userID,movieID uint,score uint8) error {
 row:=Rating{UserID:userID,MovieID:movieID,Score:score}
 return r.db.Clauses(clause.OnConflict{Columns:[]clause.Column{{Name:"user_id"},{Name:"movie_id"}},DoUpdates:clause.AssignmentColumns([]string{"score","updated_at"})}).Create(&row).Error
}
func (r *Repository) RatingSummary(movieID uint)(map[string]any,error) {
 var result struct { Average float64; Count int64 }
 err:=r.db.Model(&Rating{}).Select("COALESCE(AVG(score),0) AS average, COUNT(*) AS count").Where("movie_id = ?",movieID).Scan(&result).Error
 return map[string]any{"average":result.Average,"count":result.Count},err
}
