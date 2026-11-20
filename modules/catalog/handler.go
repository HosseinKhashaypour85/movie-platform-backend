package catalog

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"movie-platform/internal/httpx"
)

type Handler struct{ Service *Service }

func NewHandler(s *Service) *Handler { return &Handler{Service: s} }
func intParam(r *http.Request, key string) int {
	v, _ := strconv.Atoi(r.URL.Query().Get(key))
	return v
}
func uintParam(r *http.Request, key string) uint {
	v, _ := strconv.ParseUint(r.URL.Query().Get(key), 10, 64)
	return uint(v)
}
func pathID(r *http.Request) (uint, bool) {
	v, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	return uint(v), err == nil && v > 0
}

func (h *Handler) Movies(w http.ResponseWriter, r *http.Request) {
	rows, total, err := h.Service.Movies(r.URL.Query().Get("q"), uintParam(r, "genre_id"), uint16(uintParam(r, "year")), intParam(r, "limit"), intParam(r, "offset"))
	if err != nil {
		httpx.Error(w, 400, "invalid_query", "invalid catalog query")
		return
	}
	httpx.JSON(w, 200, map[string]any{"items": rows, "total": total})
}

func (h *Handler) Movie(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		httpx.Error(w, 400, "invalid_id", "invalid movie id")
		return
	}
	row, err := h.Service.Movie(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		httpx.Error(w, 404, "not_found", "movie not found")
		return
	}
	if err != nil {
		httpx.Error(w, 500, "internal_error", "internal server error")
		return
	}
	httpx.JSON(w, 200, row)
}

func (h *Handler) Series(w http.ResponseWriter, r *http.Request) {
	rows, total, err := h.Service.SeriesList(r.URL.Query().Get("q"), uintParam(r, "genre_id"), uint16(uintParam(r, "year")), intParam(r, "limit"), intParam(r, "offset"))
	if err != nil {
		httpx.Error(w, 400, "invalid_query", "invalid catalog query")
		return
	}
	httpx.JSON(w, 200, map[string]any{"items": rows, "total": total})
}

func (h *Handler) SeriesItem(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		httpx.Error(w, 400, "invalid_id", "invalid series id")
		return
	}
	row, err := h.Service.SeriesItem(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		httpx.Error(w, 404, "not_found", "series not found")
		return
	}
	if err != nil {
		httpx.Error(w, 500, "internal_error", "internal server error")
		return
	}
	httpx.JSON(w, 200, row)
}

func (h *Handler) Seasons(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		httpx.Error(w, 400, "invalid_id", "invalid series id")
		return
	}
	rows, err := h.Service.Seasons(id)
	if err != nil {
		httpx.Error(w, 500, "internal_error", "internal server error")
		return
	}
	httpx.JSON(w, 200, map[string]any{"items": rows})
}

func (h *Handler) Episodes(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		httpx.Error(w, 400, "invalid_id", "invalid season id")
		return
	}
	rows, err := h.Service.Episodes(id)
	if err != nil {
		httpx.Error(w, 500, "internal_error", "internal server error")
		return
	}
	httpx.JSON(w, 200, map[string]any{"items": rows})
}

func (h *Handler) Genres(w http.ResponseWriter, r *http.Request) {
	rows, err := h.Service.Genres()
	if err != nil { httpx.Error(w,500,"internal_error","internal server error"); return }
	httpx.JSON(w,200,map[string]any{"items":rows})
}
