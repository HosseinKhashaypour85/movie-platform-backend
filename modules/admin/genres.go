package admin

import (
	"net/http"
	"strings"

	"movie-platform/internal/httpx"
	"movie-platform/modules/catalog"
)

func (h *Handler) CreateGenre(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	}
	if !decodeAdminJSON(w, r, &in) {
		return
	}
	in.Name = strings.TrimSpace(in.Name)
	in.Slug = strings.TrimSpace(in.Slug)
	if in.Name == "" || len(in.Name) > 80 || in.Slug == "" || len(in.Slug) > 100 {
		httpx.Error(w, 400, "validation_error", "name and slug are required")
		return
	}
	row := catalog.Genre{Name: in.Name, Slug: in.Slug}
	if err := h.db.Create(&row).Error; err != nil {
		httpx.Error(w, 409, "conflict", "genre could not be created; check unique name and slug")
		return
	}
	httpx.JSON(w, 201, row)
}

func (h *Handler) UpdateGenre(w http.ResponseWriter, r *http.Request) {
	id, ok := adminID(r)
	if !ok {
		httpx.Error(w, 400, "invalid_id", "invalid genre id")
		return
	}
	var in struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	}
	if !decodeAdminJSON(w, r, &in) {
		return
	}
	in.Name = strings.TrimSpace(in.Name)
	in.Slug = strings.TrimSpace(in.Slug)
	if in.Name == "" || len(in.Name) > 80 || in.Slug == "" || len(in.Slug) > 100 {
		httpx.Error(w, 400, "validation_error", "name and slug are required")
		return
	}
	var row catalog.Genre
	if err := h.db.First(&row, id).Error; err != nil {
		httpx.Error(w, 404, "not_found", "genre not found")
		return
	}
	row.Name = in.Name
	row.Slug = in.Slug
	if err := h.db.Save(&row).Error; err != nil {
		httpx.Error(w, 409, "conflict", "genre could not be updated")
		return
	}
	httpx.JSON(w, 200, row)
}

func (h *Handler) DeleteGenre(w http.ResponseWriter, r *http.Request) {
	id, ok := adminID(r)
	if !ok {
		httpx.Error(w, 400, "invalid_id", "invalid genre id")
		return
	}
	var movieCount, seriesCount int64
	if err := h.db.Table("movie_genres").Where("genre_id = ?", id).Count(&movieCount).Error; err != nil {
		httpx.Error(w, 500, "internal_error", "internal server error")
		return
	}
	if err := h.db.Table("series_genres").Where("genre_id = ?", id).Count(&seriesCount).Error; err != nil {
		httpx.Error(w, 500, "internal_error", "internal server error")
		return
	}
	if movieCount+seriesCount > 0 {
		httpx.Error(w, 409, "genre_in_use", "remove genre assignments before deleting")
		return
	}
	res := h.db.Delete(&catalog.Genre{}, id)
	if res.Error != nil {
		httpx.Error(w, 500, "internal_error", "internal server error")
		return
	}
	if res.RowsAffected == 0 {
		httpx.Error(w, 404, "not_found", "genre not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
