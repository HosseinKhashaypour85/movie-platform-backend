package media

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"movie-platform/internal/httpx"
	"movie-platform/internal/middleware"
	"movie-platform/modules/catalog"
)

type UploadHandler struct {
	db *gorm.DB
	s3 *S3Presigner
}

func NewUploadHandler(db *gorm.DB, s3 *S3Presigner) *UploadHandler {
	return &UploadHandler{db:db,s3:s3}
}

type uploadRequest struct {
	TargetType string `json:"target_type"`
	TargetID uint `json:"target_id"`
	Kind string `json:"kind"`
	ContentType string `json:"content_type"`
}

func (h *UploadHandler) CreateUpload(w http.ResponseWriter,r *http.Request) {
	w.Header().Set("Cache-Control","no-store")
	if h.s3==nil {
		httpx.Error(w,http.StatusServiceUnavailable,"storage_unavailable","object storage is not configured")
		return
	}
	var in uploadRequest
	if !decodeMedia(w,r,&in) { return }
	if in.TargetID==0 || (in.TargetType!="movie" && in.TargetType!="episode") {
		httpx.Error(w,http.StatusBadRequest,"validation_error","target_type and target_id are invalid")
		return
	}
	if in.Kind!="video" && in.Kind!="trailer" && in.Kind!="poster" {
		httpx.Error(w,http.StatusBadRequest,"validation_error","kind must be video, trailer or poster")
		return
	}
	allowed:=map[string]string{
		"video/mp4":".mp4",
		"image/jpeg":".jpg",
		"image/png":".png",
		"image/webp":".webp",
	}
	ext,ok:=allowed[strings.ToLower(strings.TrimSpace(in.ContentType))]
	if !ok || (in.Kind=="video" && ext!=".mp4") || (in.Kind=="poster" && ext==".mp4") {
		httpx.Error(w,http.StatusBadRequest,"validation_error","unsupported content type for asset kind")
		return
	}
	var prefix string
	if in.TargetType=="movie" {
		var row catalog.Movie
		if err:=h.db.WithContext(r.Context()).First(&row,in.TargetID).Error; err!=nil {
			if err==gorm.ErrRecordNotFound { httpx.Error(w,404,"target_not_found","movie not found") } else { httpx.Error(w,500,"internal_error","internal server error") }
			return
		}
		prefix="movies"
	} else {
		var row catalog.Episode
		if err:=h.db.WithContext(r.Context()).First(&row,in.TargetID).Error; err!=nil {
			if err==gorm.ErrRecordNotFound { httpx.Error(w,404,"target_not_found","episode not found") } else { httpx.Error(w,500,"internal_error","internal server error") }
			return
		}
		prefix="episodes"
	}
	var random [16]byte
	if _,err:=rand.Read(random[:]); err!=nil {
		httpx.Error(w,500,"internal_error","could not create upload key")
		return
	}
	objectKey:=path.Join(prefix,strconv.FormatUint(uint64(in.TargetID),10),in.Kind,time.Now().UTC().Format("2006/01"),hex.EncodeToString(random[:])+ext)
	uploadURL,err:=h.s3.PresignPut(objectKey,in.ContentType,time.Now(),10*time.Minute)
	if err!=nil {
		httpx.Error(w,500,"storage_error","could not create upload URL")
		return
	}
	httpx.JSON(w,http.StatusCreated,map[string]any{
		"method":"PUT",
		"upload_url":uploadURL,
		"headers":map[string]string{"Content-Type":in.ContentType},
		"storage_key":objectKey,
		"expires_in_seconds":600,
		"asset_kind":in.Kind,
		"target_type":in.TargetType,
		"target_id":in.TargetID,
	})
}

func RegisterUploadRoutes(r chi.Router,h *UploadHandler) {
	r.Route("/media",func(a chi.Router) {
		a.Use(middleware.RequireRole("admin","super_admin"))
		a.Post("/uploads",h.CreateUpload)
	})
}
