package admin

import (
	"net/http"

	"movie-platform/internal/httpx"
	"movie-platform/modules/catalog"
	"movie-platform/modules/users"

	"gorm.io/gorm"
)

type Handler struct{ db *gorm.DB }

func NewHandler(db *gorm.DB) *Handler { return &Handler{db: db} }

func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	var userCount, movieCount, seriesCount, episodeCount int64
	if err := h.db.Model(&users.User{}).Count(&userCount).Error; err != nil {
		httpx.Error(w, 500, "internal_error", "internal server error")
		return
	}
	if err := h.db.Model(&catalog.Movie{}).Count(&movieCount).Error; err != nil {
		httpx.Error(w, 500, "internal_error", "internal server error")
		return
	}
	if err := h.db.Model(&catalog.Series{}).Count(&seriesCount).Error; err != nil {
		httpx.Error(w, 500, "internal_error", "internal server error")
		return
	}
	if err := h.db.Model(&catalog.Episode{}).Count(&episodeCount).Error; err != nil {
		httpx.Error(w, 500, "internal_error", "internal server error")
		return
	}
	httpx.JSON(w, 200, map[string]any{
		"users":    userCount,
		"movies":   movieCount,
		"series":   seriesCount,
		"episodes": episodeCount,
	})
}
