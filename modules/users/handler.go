package users

import (
	"encoding/json"
		"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"movie-platform/internal/httpx"
	"movie-platform/internal/middleware"
)

type Handler struct { db *gorm.DB }
func NewHandler(db *gorm.DB) *Handler { return &Handler{db:db} }

func (h *Handler) Profile(w http.ResponseWriter,r *http.Request) {
	id,ok:=middleware.UserIDFromContext(r.Context()); if !ok { httpx.Error(w,401,"unauthorized","authentication required"); return }
	var user User
	if err:=h.db.Select("id","name","email","role","is_active","created_at","updated_at").First(&user,id).Error; err!=nil || !user.IsActive { httpx.Error(w,401,"unauthorized","authentication required"); return }
	httpx.JSON(w,200,user)
}
func (h *Handler) UpdateProfile(w http.ResponseWriter,r *http.Request) {
	id,ok:=middleware.UserIDFromContext(r.Context()); if !ok { httpx.Error(w,401,"unauthorized","authentication required"); return }
	var in struct { Name string `json:"name"` }
	if !decodeUser(w,r,&in) { return }; in.Name=strings.TrimSpace(in.Name)
	if len(in.Name)<2 || len(in.Name)>100 { httpx.Error(w,400,"validation_error","name must be between 2 and 100 characters"); return }
	res:=h.db.Model(&User{}).Where("id = ? AND is_active = ?",id,true).Update("name",in.Name)
	if res.Error!=nil { httpx.Error(w,500,"internal_error","internal server error"); return }
	if res.RowsAffected==0 { httpx.Error(w,404,"not_found","user not found"); return }
	httpx.JSON(w,200,map[string]string{"name":in.Name})
}
func (h *Handler) ChangePassword(w http.ResponseWriter,r *http.Request) {
	id,ok:=middleware.UserIDFromContext(r.Context()); if !ok { httpx.Error(w,401,"unauthorized","authentication required"); return }
	var in struct { CurrentPassword string `json:"current_password"`; NewPassword string `json:"new_password"` }
	if !decodeUser(w,r,&in) { return }
	if len(in.NewPassword)<12 || len([]byte(in.NewPassword))>72 { httpx.Error(w,400,"validation_error","new password must be 12-72 bytes"); return }
	var user User
	if err:=h.db.First(&user,id).Error; err!=nil || !user.IsActive { httpx.Error(w,401,"unauthorized","authentication required"); return }
	if bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(in.CurrentPassword))!=nil { httpx.Error(w,400,"invalid_password","current password is incorrect"); return }
	hash,err:=bcrypt.GenerateFromPassword([]byte(in.NewPassword),bcrypt.DefaultCost); if err!=nil { httpx.Error(w,500,"internal_error","internal server error"); return }
	now:=time.Now().UTC()
	err=h.db.Transaction(func(tx *gorm.DB) error {
		if e:=tx.Model(&User{}).Where("id = ?",id).Update("password",string(hash)).Error; e!=nil { return e }
		return tx.Table("sessions").Where("user_id = ? AND revoked_at IS NULL",id).Update("revoked_at",now).Error
	})
	if err!=nil { httpx.Error(w,500,"internal_error","internal server error"); return }
	w.WriteHeader(http.StatusNoContent)
}
func decodeUser(w http.ResponseWriter,r *http.Request,dst any) bool {
	r.Body=http.MaxBytesReader(w,r.Body,4096); d:=json.NewDecoder(r.Body); d.DisallowUnknownFields()
	if err:=d.Decode(dst); err!=nil { httpx.Error(w,400,"invalid_body","invalid request body"); return false }
	var extra any; if err:=d.Decode(&extra); err!=io.EOF { httpx.Error(w,400,"invalid_body","request must contain one JSON object"); return false }
	return true
}
func RegisterRoutes(r interface{ Get(string,http.HandlerFunc); Patch(string,http.HandlerFunc); Post(string,http.HandlerFunc) },h *Handler) {
	r.Get("/profile",h.Profile)
	r.Patch("/profile",h.UpdateProfile)
	r.Post("/profile/change-password",h.ChangePassword)
}

