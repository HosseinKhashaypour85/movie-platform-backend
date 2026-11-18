package watch_history

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"io"
	"movie-platform/internal/httpx"
	"movie-platform/internal/middleware"
	"net/http"
	"strconv"
)

type Handler struct{ Service *Service }

func NewHandler(s *Service) *Handler   { return &Handler{Service: s} }
func uid(r *http.Request) (uint, bool) { return middleware.UserIDFromContext(r.Context()) }
func movieID(r *http.Request) (uint, bool) {
	v, e := strconv.ParseUint(chi.URLParam(r, "movieID"), 10, 64)
	return uint(v), e == nil && v > 0
}
func decodeBody(r *http.Request, v any) error {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if e := d.Decode(v); e != nil {
		return e
	}
	var x any
	if e := d.Decode(&x); e != io.EOF {
		return errors.New("multiple JSON values")
	}
	return nil
}
func (h *Handler) Save(w http.ResponseWriter, r *http.Request) {
	u, ok := uid(r)
	if !ok {
		httpx.Error(w, 401, "unauthorized", "authentication required")
		return
	}
	m, ok := movieID(r)
	if !ok {
		httpx.Error(w, 400, "invalid_movie_id", "movie id is invalid")
		return
	}
	var q struct {
		Progress int64 `json:"progress_seconds"`
		Duration int64 `json:"duration_seconds"`
	}
	if decodeBody(r, &q) != nil {
		httpx.Error(w, 400, "invalid_json", "invalid request body")
		return
	}
	if h.Service.Save(u, m, q.Progress, q.Duration) != nil {
		httpx.Error(w, 400, "invalid_progress", "progress or duration is invalid")
		return
	}
	httpx.JSON(w, 200, map[string]string{"status": "saved"})
}
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	u, _ := uid(r)
	l, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	o, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if l < 1 || l > 100 {
		l = 20
	}
	if o < 0 {
		o = 0
	}
	items, e := h.Service.List(u, l, o)
	if e != nil {
		httpx.Error(w, 500, "internal_error", "internal server error")
		return
	}
	httpx.JSON(w, 200, map[string]any{"items": items, "limit": l, "offset": o})
}
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	u, _ := uid(r)
	m, ok := movieID(r)
	if !ok {
		httpx.Error(w, 400, "invalid_movie_id", "movie id is invalid")
		return
	}
	item, e := h.Service.Get(u, m)
	if errors.Is(e, gorm.ErrRecordNotFound) {
		httpx.Error(w, 404, "not_found", "watch history not found")
		return
	}
	if e != nil {
		httpx.Error(w, 500, "internal_error", "internal server error")
		return
	}
	httpx.JSON(w, 200, item)
}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	u, _ := uid(r)
	m, ok := movieID(r)
	if !ok {
		httpx.Error(w, 400, "invalid_movie_id", "movie id is invalid")
		return
	}
	if e := h.Service.Delete(u, m); e != nil {
		httpx.Error(w, 500, "internal_error", "internal server error")
		return
	}
	w.WriteHeader(204)
}
