package admin

import (
 "net/http"
 "movie-platform/internal/httpx"
 "gorm.io/gorm"
)

type Handler struct{db *gorm.DB}
func NewHandler(db *gorm.DB)*Handler{return &Handler{db:db}}
func(h *Handler)Dashboard(w http.ResponseWriter,r *http.Request){
 var users,movies,series,episodes int64
 if e:=h.db.Table("users").Count(&users).Error;e!=nil{httpx.Error(w,500,"internal_error","internal server error");return}
 if e:=h.db.Table("movies").Count(&movies).Error;e!=nil{httpx.Error(w,500,"internal_error","internal server error");return}
 if e:=h.db.Table("series").Count(&series).Error;e!=nil{httpx.Error(w,500,"internal_error","internal server error");return}
 if e:=h.db.Table("episodes").Count(&episodes).Error;e!=nil{httpx.Error(w,500,"internal_error","internal server error");return}
 httpx.JSON(w,200,map[string]any{"users":users,"movies":movies,"series":series,"episodes":episodes})
}
