package ads

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"movie-platform/internal/httpx"
)

type Handler struct{ db *gorm.DB }

func NewHandler(db *gorm.DB) *Handler { return &Handler{db: db} }

func (h *Handler) Serve(w http.ResponseWriter, r *http.Request) {
	placement := strings.TrimSpace(r.URL.Query().Get("placement"))
	if placement == "" || len(placement) > 40 {
		httpx.Error(w, 400, "invalid_placement", "placement is required")
		return
	}
	now := time.Now().UTC()
	var rows []Campaign
	err := h.db.Where("is_active = ? AND placement = ? AND starts_at <= ? AND ends_at > ?", true, placement, now, now).
		Select("id", "name", "placement", "format", "media_url", "click_url", "advertiser").Order("id ASC").Limit(5).Find(&rows).Error
	if err != nil {
		httpx.Error(w, 500, "internal_error", "internal server error")
		return
	}
	httpx.JSON(w, 200, map[string]any{"items": rows})
}

func (h *Handler) Track(w http.ResponseWriter, r *http.Request) {
	id64, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id64 == 0 {
		httpx.Error(w, 400, "invalid_id", "invalid campaign id")
		return
	}
	var body struct {
		Type      string `json:"type"`
		Placement string `json:"placement"`
	}
	if !decodeEvent(w, r, &body) {
		return
	}
	if body.Type != "impression" && body.Type != "click" {
		httpx.Error(w, 400, "invalid_event", "type must be impression or click")
		return
	}
	if len(body.Placement) > 40 {
		httpx.Error(w, 400, "invalid_event", "placement is invalid")
		return
	}
	now := time.Now().UTC()
	err = h.db.Transaction(func(tx *gorm.DB) error {
		var c Campaign
		if e := tx.Where("id = ? AND is_active = ? AND starts_at <= ? AND ends_at > ?", id64, true, now, now).First(&c).Error; e != nil {
			return e
		}
		if e := tx.Create(&Event{CampaignID: c.ID, EventType: body.Type, Placement: body.Placement}).Error; e != nil {
			return e
		}
		column := "impression_count"
		if body.Type == "click" {
			column = "click_count"
		}
		return tx.Model(&Campaign{}).Where("id = ?", c.ID).UpdateColumn(column, gorm.Expr(column+" + 1")).Error
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		httpx.Error(w, 404, "not_found", "campaign not found")
		return
	}
	if err != nil {
		httpx.Error(w, 500, "internal_error", "internal server error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func decodeEvent(w http.ResponseWriter, r *http.Request, dst any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, 2048)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		httpx.Error(w, 400, "invalid_body", "invalid request body")
		return false
	}
	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		httpx.Error(w, 400, "invalid_body", "request must contain one JSON object")
		return false
	}
	return true
}

func (h *Handler) CreateCampaign(w http.ResponseWriter, r *http.Request) {
	var in Campaign
	if !decodeEvent(w, r, &in) {
		return
	}
	in.Name = strings.TrimSpace(in.Name)
	in.Placement = strings.TrimSpace(in.Placement)
	if in.Name == "" || len(in.Name) > 160 || in.Placement == "" || len(in.Placement) > 40 || (in.Format != "image" && in.Format != "video" && in.Format != "banner") || !validURL(in.MediaURL) || (in.ClickURL != "" && !validURL(in.ClickURL)) || !in.EndsAt.After(in.StartsAt) {
		httpx.Error(w, 400, "validation_error", "campaign fields, URLs or date range are invalid")
		return
	}
	in.ID = 0
	in.ImpressionCount = 0
	in.ClickCount = 0
	if err := h.db.Create(&in).Error; err != nil {
		httpx.Error(w, 500, "internal_error", "internal server error")
		return
	}
	httpx.JSON(w, 201, in)
}
func (h *Handler) ListCampaigns(w http.ResponseWriter, r *http.Request) {
	var rows []Campaign
	if err := h.db.Order("created_at DESC").Limit(200).Find(&rows).Error; err != nil {
		httpx.Error(w, 500, "internal_error", "internal server error")
		return
	}
	httpx.JSON(w, 200, map[string]any{"items": rows})
}
func (h *Handler) UpdateCampaign(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id == 0 {
		httpx.Error(w, 400, "invalid_id", "invalid campaign id")
		return
	}
	var in Campaign
	if !decodeEvent(w, r, &in) {
		return
	}
	in.Name = strings.TrimSpace(in.Name)
	in.Placement = strings.TrimSpace(in.Placement)
	if in.Name == "" || len(in.Name) > 160 || in.Placement == "" || len(in.Placement) > 40 || (in.Format != "image" && in.Format != "video" && in.Format != "banner") || !validURL(in.MediaURL) || (in.ClickURL != "" && !validURL(in.ClickURL)) || !in.EndsAt.After(in.StartsAt) {
		httpx.Error(w, 400, "validation_error", "campaign fields, URLs or date range are invalid")
		return
	}
	var old Campaign
	if err := h.db.First(&old, id).Error; err != nil {
		httpx.Error(w, 404, "not_found", "campaign not found")
		return
	}
	old.Name = in.Name
	old.Placement = in.Placement
	old.Format = in.Format
	old.MediaURL = in.MediaURL
	old.ClickURL = in.ClickURL
	old.Advertiser = in.Advertiser
	old.IsActive = in.IsActive
	old.StartsAt = in.StartsAt
	old.EndsAt = in.EndsAt
	if err := h.db.Save(&old).Error; err != nil {
		httpx.Error(w, 500, "internal_error", "internal server error")
		return
	}
	httpx.JSON(w, 200, old)
}
func validURL(v string) bool { return strings.HasPrefix(v, "https://") && len(v) <= 2048 }
