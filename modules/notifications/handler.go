package notifications

import (
	"net/http"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"movie-platform/internal/httpx"
	"movie-platform/internal/middleware"
)

type Handler struct { db *gorm.DB }
func NewHandler(db *gorm.DB) *Handler { return &Handler{db:db} }

func (h *Handler) List(w http.ResponseWriter,r *http.Request) {
	uid,_:=middleware.UserIDFromContext(r.Context())
	page,_:=strconv.Atoi(r.URL.Query().Get("page")); limit,_:=strconv.Atoi(r.URL.Query().Get("limit"))
	if page<1 { page=1 }; if limit<1 || limit>100 { limit=20 }
	q:=h.db.Model(&Notification{}).Where("user_id = ?",uid)
	var total int64; if err:=q.Count(&total).Error; err!=nil { httpx.Error(w,500,"internal_error","internal server error"); return }
	var rows []Notification
	if err:=q.Order("created_at DESC").Limit(limit).Offset((page-1)*limit).Find(&rows).Error; err!=nil { httpx.Error(w,500,"internal_error","internal server error"); return }
	httpx.JSON(w,200,map[string]any{"items":rows,"total":total,"page":page,"limit":limit})
}
func (h *Handler) MarkRead(w http.ResponseWriter,r *http.Request) {
	uid,_:=middleware.UserIDFromContext(r.Context()); id,e:=strconv.ParseUint(chi.URLParam(r,"id"),10,64)
	if e!=nil || id==0 { httpx.Error(w,400,"invalid_id","invalid notification id"); return }
	now:=time.Now().UTC()
	res:=h.db.Model(&Notification{}).Where("id = ? AND user_id = ?",id,uid).Update("read_at",now)
	if res.Error!=nil { httpx.Error(w,500,"internal_error","internal server error"); return }
	if res.RowsAffected==0 { httpx.Error(w,404,"not_found","notification not found"); return }
	w.WriteHeader(http.StatusNoContent)
}
func (h *Handler) MarkAllRead(w http.ResponseWriter,r *http.Request) {
	uid,_:=middleware.UserIDFromContext(r.Context())
	if err:=h.db.Model(&Notification{}).Where("user_id = ? AND read_at IS NULL",uid).Update("read_at",time.Now().UTC()).Error; err!=nil { httpx.Error(w,500,"internal_error","internal server error"); return }
	w.WriteHeader(http.StatusNoContent)
}
func (h *Handler) Broadcast(w http.ResponseWriter,r *http.Request) {
	var in struct { Title string `json:"title"`; Body string `json:"body"`; Type string `json:"type"`; ActionURL string `json:"action_url"` }
	r.Body=http.MaxBytesReader(w,r.Body,4096)
	d:=json.NewDecoder(r.Body); d.DisallowUnknownFields()
	if d.Decode(&in)!=nil { httpx.Error(w,400,"invalid_body","invalid request body"); return }
	in.Title=strings.TrimSpace(in.Title); in.Body=strings.TrimSpace(in.Body)
	if in.Title=="" || len(in.Title)>160 || in.Body=="" || len(in.Body)>4000 || (in.ActionURL!="" && (!strings.HasPrefix(in.ActionURL,"https://") || len(in.ActionURL)>2048)) { httpx.Error(w,400,"validation_error","notification fields are invalid"); return }
	if in.Type=="" { in.Type="general" }; if len(in.Type)>40 { httpx.Error(w,400,"validation_error","type is too long"); return }
	// Broadcast in bounded batches to avoid a single unbounded insert.
	var lastID uint
	batch:=500
	for {
		var ids []uint
		if err:=h.db.Model(&struct{ID uint; IsActive bool}{ }).Where("id > ? AND is_active = ?",lastID,true).Order("id ASC").Limit(batch).Pluck("id",&ids).Error; err!=nil { httpx.Error(w,500,"internal_error","internal server error"); return }
		if len(ids)==0 { break }
		rows:=make([]Notification,0,len(ids))
		for _,id:=range ids { rows=append(rows,Notification{UserID:id,Title:in.Title,Body:in.Body,Type:in.Type,ActionURL:in.ActionURL}) }
		if err:=h.db.CreateInBatches(&rows,batch).Error; err!=nil { httpx.Error(w,500,"internal_error","broadcast failed"); return }
		lastID=ids[len(ids)-1]
	}
	httpx.JSON(w,202,map[string]string{"status":"queued"})
}

func RegisterUserRoutes(r chi.Router,h *Handler) {
	r.Route("/notifications",func(n chi.Router) { n.Get("/",h.List); n.Patch("/{id}/read",h.MarkRead); n.Post("/read-all",h.MarkAllRead) })
}
func RegisterAdminRoutes(r chi.Router,h *Handler) {
	r.Route("/notifications",func(n chi.Router) { n.Use(middleware.RequireRole("admin","super_admin")); n.Post("/broadcast",h.Broadcast) })
}
