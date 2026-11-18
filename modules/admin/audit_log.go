package admin

import (
	"context"
	"log/slog"
	"net"
	"net/http"

	"movie-platform/internal/httpx"
	"strconv"
	"strings"
	"time"

	"movie-platform/internal/middleware"
	"gorm.io/gorm"
)

type AuditLog struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	ActorUserID  uint      `gorm:"not null;index" json:"actor_user_id"`
	Action       string    `gorm:"size:10;not null;index" json:"action"`
	ResourceType string    `gorm:"size:40;not null;index" json:"resource_type"`
	ResourceID   string    `gorm:"size:64" json:"resource_id,omitempty"`
	Path         string    `gorm:"size:255;not null" json:"path"`
	StatusCode   int       `gorm:"not null;index" json:"status_code"`
	IPAddress    string    `gorm:"size:45" json:"ip_address,omitempty"`
	UserAgent    string    `gorm:"size:255" json:"user_agent,omitempty"`
	RequestID    string    `gorm:"size:64;index" json:"request_id,omitempty"`
	CreatedAt    time.Time `gorm:"index" json:"created_at"`
}

type auditResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *auditResponseWriter) WriteHeader(status int) {
	if w.status != 0 {
		return
	}
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *auditResponseWriter) Write(p []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(p)
}

func AuditMiddleware(db *gorm.DB, logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
			default:
				next.ServeHTTP(w, r)
				return
			}

			recorder := &auditResponseWriter{ResponseWriter: w}
			next.ServeHTTP(recorder, r)
			status := recorder.status
			if status == 0 {
				status = http.StatusOK
			}

			actorID, ok := middleware.UserIDFromContext(r.Context())
			if !ok {
				logger.Error("admin audit skipped: actor missing", "path", r.URL.Path)
				return
			}

			parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
			resourceType, resourceID := "unknown", ""
			if len(parts) >= 4 {
				resourceType = parts[3]
			}
			if len(parts) >= 5 {
				resourceID = parts[4]
			}
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr
			}
			userAgent := r.UserAgent()
			if len(userAgent) > 255 {
				userAgent = userAgent[:255]
			}
			requestID := r.Header.Get("X-Request-ID")
			if len(requestID) > 64 {
				requestID = requestID[:64]
			}

			entry := AuditLog{
				ActorUserID: actorID,
				Action: r.Method,
				ResourceType: resourceType,
				ResourceID: resourceID,
				Path: r.URL.Path,
				StatusCode: status,
				IPAddress: ip,
				UserAgent: userAgent,
				RequestID: requestID,
			}
			ctx, cancel := context.WithTimeout(context.WithoutCancel(r.Context()), 2*time.Second)
			defer cancel()
			if err := db.WithContext(ctx).Create(&entry).Error; err != nil {
				logger.Error("admin audit write failed", "error", err, "actor_user_id", actorID, "request_id", requestID)
			}
		})
	}
}

func (h *Handler) ListAuditLogs(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 25
	}
	if limit > 100 {
		limit = 100
	}

	query := h.db.Model(&AuditLog{})
	if actor := r.URL.Query().Get("actor_user_id"); actor != "" {
		id, err := strconv.ParseUint(actor, 10, 64)
		if err != nil || id == 0 {
			httpx.Error(w, http.StatusBadRequest, "validation_error", "invalid actor_user_id")
			return
		}
		query = query.Where("actor_user_id = ?", id)
	}
	if action := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("action"))); action != "" {
		switch action {
		case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
			query = query.Where("action = ?", action)
		default:
			httpx.Error(w, http.StatusBadRequest, "validation_error", "invalid action")
			return
		}
	}
	if resource := strings.TrimSpace(r.URL.Query().Get("resource_type")); resource != "" {
		if len(resource) > 40 {
			httpx.Error(w, http.StatusBadRequest, "validation_error", "invalid resource_type")
			return
		}
		query = query.Where("resource_type = ?", resource)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal_error", "internal server error")
		return
	}
	var rows []AuditLog
	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").Order("id DESC").Limit(limit).Offset(offset).Find(&rows).Error; err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	httpx.JSON(w, http.StatusOK, map[string]any{
		"items": rows,
		"pagination": map[string]any{"page": page, "limit": limit, "total": total},
	})
}
