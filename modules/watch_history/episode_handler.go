package watch_history

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"movie-platform/internal/httpx"
	"movie-platform/internal/middleware"
	"movie-platform/modules/catalog"
)

func (h *Handler) SaveEpisode(w http.ResponseWriter, r *http.Request) {
	uid, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.Error(w, 401, "unauthorized", "authentication required")
		return
	}
	id, e := strconv.ParseUint(chi.URLParam(r, "episodeID"), 10, 64)
	if e != nil || id == 0 {
		httpx.Error(w, 400, "invalid_episode_id", "episode id is invalid")
		return
	}
	var in struct {
		Progress int64 `json:"progress_seconds"`
		Duration int64 `json:"duration_seconds"`
	}
	if decodeBody(r, &in) != nil || in.Duration <= 0 || in.Progress < 0 || in.Progress > in.Duration {
		httpx.Error(w, 400, "invalid_progress", "progress or duration is invalid")
		return
	}
	var episode catalog.Episode
	if err := h.Service.Repo.DB.Joins("JOIN seasons ON seasons.id = episodes.season_id").Joins("JOIN series ON series.id = seasons.series_id").Where("episodes.id = ? AND episodes.is_published = ? AND series.is_published = ?", id, true, true).First(&episode).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			httpx.Error(w, 404, "not_found", "published episode not found")
		} else {
			httpx.Error(w, 500, "internal_error", "internal server error")
		}
		return
	}
	completed := (in.Duration > 30 && in.Progress >= in.Duration-30) || (in.Duration <= 30 && in.Progress >= int64(float64(in.Duration)*0.9))
	row := EpisodeWatchHistory{UserID: uid, EpisodeID: uint(id), ProgressSeconds: in.Progress, DurationSeconds: in.Duration, Completed: completed, WatchedAt: time.Now().UTC()}
	err := h.Service.Repo.DB.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "user_id"}, {Name: "episode_id"}}, DoUpdates: clause.AssignmentColumns([]string{"progress_seconds", "duration_seconds", "completed", "watched_at", "updated_at"})}).Create(&row).Error
	if err != nil {
		httpx.Error(w, 500, "internal_error", "internal server error")
		return
	}
	httpx.JSON(w, 200, map[string]string{"status": "saved"})
}
func (h *Handler) EpisodeHistory(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFromContext(r.Context())
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	var rows []EpisodeWatchHistory
	err := h.Service.Repo.DB.Where("user_id = ?", uid).Order("watched_at DESC").Limit(limit).Offset(offset).Find(&rows).Error
	if err != nil {
		httpx.Error(w, 500, "internal_error", "internal server error")
		return
	}
	httpx.JSON(w, 200, map[string]any{"items": rows, "limit": limit, "offset": offset})
}
func (h *Handler) DeleteEpisodeHistory(w http.ResponseWriter, r *http.Request) {
	uid, _ := middleware.UserIDFromContext(r.Context())
	id, e := strconv.ParseUint(chi.URLParam(r, "episodeID"), 10, 64)
	if e != nil || id == 0 {
		httpx.Error(w, 400, "invalid_episode_id", "episode id is invalid")
		return
	}
	res := h.Service.Repo.DB.Where("user_id = ? AND episode_id = ?", uid, id).Delete(&EpisodeWatchHistory{})
	if res.Error != nil {
		httpx.Error(w, 500, "internal_error", "internal server error")
		return
	}
	if res.RowsAffected == 0 {
		httpx.Error(w, 404, "not_found", "episode history not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
