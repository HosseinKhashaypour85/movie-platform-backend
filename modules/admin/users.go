package admin

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

	"movie-platform/internal/httpx"
	"movie-platform/internal/middleware"
	"movie-platform/modules/users"
)

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page < 1 { page = 1 }
	if limit < 1 || limit > 100 { limit = 20 }
	q := h.db.Model(&users.User{})
	if status := r.URL.Query().Get("status"); status == "active" { q = q.Where("is_active = ?", true) } else if status == "disabled" { q = q.Where("is_active = ?", false) }
	if search := strings.TrimSpace(r.URL.Query().Get("q")); search != "" {
		if len(search) > 120 { httpx.Error(w,400,"invalid_query","search query is too long"); return }
		like := "%" + search + "%"
		q = q.Where("name LIKE ? OR email LIKE ?",like,like)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil { httpx.Error(w,500,"internal_error","internal server error"); return }
	var rows []users.User
	if err := q.Select("id","name","email","role","is_active","created_at","updated_at").Order("created_at DESC").Limit(limit).Offset((page-1)*limit).Find(&rows).Error; err != nil {
		httpx.Error(w,500,"internal_error","internal server error"); return
	}
	httpx.JSON(w,200,map[string]any{"items":rows,"total":total,"page":page,"limit":limit})
}

func (h *Handler) SetUserStatus(w http.ResponseWriter, r *http.Request) {
	id, ok := adminID(r); if !ok { httpx.Error(w,400,"invalid_id","invalid user id"); return }
	var body struct { IsActive *bool `json:"is_active"` }
	if !decodeAdminJSON(w,r,&body) { return }
	if body.IsActive == nil { httpx.Error(w,400,"validation_error","is_active is required"); return }
	res := h.db.Model(&users.User{}).Where("id = ?",id).Update("is_active",*body.IsActive)
	if res.Error != nil { httpx.Error(w,500,"internal_error","internal server error"); return }
	if res.RowsAffected == 0 { httpx.Error(w,404,"not_found","user not found"); return }
	if !*body.IsActive { _ = h.db.Model(&SessionRecord{}).Where("user_id = ? AND revoked_at IS NULL",id).Update("revoked_at", gorm.Expr("UTC_TIMESTAMP(3)")).Error }
	httpx.JSON(w,200,map[string]bool{"is_active":*body.IsActive})
}

type SessionRecord struct { ID uint; UserID uint; RevokedAt *time.Time }
func (SessionRecord) TableName() string { return "sessions" }

func (h *Handler) SetUserRole(w http.ResponseWriter, r *http.Request) {
	if role, _ := middleware.RoleFromContext(r.Context()); role != "super_admin" { httpx.Error(w,403,"forbidden","super_admin role required"); return }
	id, ok := adminID(r); if !ok { httpx.Error(w,400,"invalid_id","invalid user id"); return }
	var body struct { Role string `json:"role"` }
	if !decodeAdminJSON(w,r,&body) { return }
	if body.Role != "user" && body.Role != "admin" && body.Role != "super_admin" { httpx.Error(w,400,"validation_error","role is invalid"); return }
	if res := h.db.Model(&users.User{}).Where("id = ?",id).Update("role",body.Role); res.Error != nil { httpx.Error(w,500,"internal_error","internal server error"); return } else if res.RowsAffected == 0 { httpx.Error(w,404,"not_found","user not found"); return }
	httpx.JSON(w,200,map[string]string{"role":body.Role})
}
