package media

import (
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"movie-platform/internal/httpx"
	"movie-platform/modules/catalog"
)

const (
	maxVideoUploadBytes = uint64(2 * 1024 * 1024 * 1024)
	maxImageUploadBytes = uint64(20 * 1024 * 1024)
)

type completeUploadRequest struct {
	TargetType string `json:"target_type"`
	TargetID uint `json:"target_id"`
	Kind string `json:"kind"`
	StorageKey string `json:"storage_key"`
}

func (h *UploadHandler) CompleteUpload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	if h.s3 == nil {
		httpx.Error(w, http.StatusServiceUnavailable, "storage_unavailable", "object storage is not configured")
		return
	}
	var in completeUploadRequest
	if !decodeMedia(w, r, &in) { return }
	if in.TargetID == 0 || (in.TargetType != "movie" && in.TargetType != "episode") ||
		(in.Kind != "video" && in.Kind != "trailer" && in.Kind != "poster") ||
		len(in.StorageKey) == 0 || len(in.StorageKey) > 1024 || strings.Contains(in.StorageKey, "\\") ||
		path.Clean(in.StorageKey) != in.StorageKey || strings.HasPrefix(in.StorageKey, "/") {
		httpx.Error(w, http.StatusBadRequest, "validation_error", "invalid upload target, kind or storage key")
		return
	}
	prefix := "movies/"
	if in.TargetType == "episode" { prefix = "episodes/" }
	expectedPrefix := prefix + strconv.FormatUint(uint64(in.TargetID), 10) + "/" + in.Kind + "/"
	if !strings.HasPrefix(in.StorageKey, expectedPrefix) {
		httpx.Error(w, http.StatusBadRequest, "invalid_storage_key", "storage key does not belong to the requested target")
		return
	}
	if in.TargetType == "movie" {
		var target catalog.Movie
		if err := h.db.WithContext(r.Context()).First(&target, in.TargetID).Error; err != nil {
			if err == gorm.ErrRecordNotFound { httpx.Error(w, 404, "target_not_found", "movie not found") } else { httpx.Error(w, 500, "internal_error", "internal server error") }
			return
		}
	} else {
		var target catalog.Episode
		if err := h.db.WithContext(r.Context()).First(&target, in.TargetID).Error; err != nil {
			if err == gorm.ErrRecordNotFound { httpx.Error(w, 404, "target_not_found", "episode not found") } else { httpx.Error(w, 500, "internal_error", "internal server error") }
			return
		}
	}
	var existing Asset
	if err := h.db.WithContext(r.Context()).Where("storage_key = ?", in.StorageKey).First(&existing).Error; err == nil {
		httpx.Error(w, http.StatusConflict, "upload_already_finalized", "this upload has already been finalized")
		return
	} else if err != gorm.ErrRecordNotFound {
		httpx.Error(w, 500, "internal_error", "internal server error")
		return
	}
	meta, err := h.s3.HeadObject(r.Context(), in.StorageKey)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") { httpx.Error(w, 422, "upload_not_found", "uploaded object was not found in storage") } else { httpx.Error(w, 502, "storage_verification_failed", "could not verify uploaded object") }
		return
	}
	allowed := map[string]bool{"video/mp4":true, "image/jpeg":true, "image/png":true, "image/webp":true}
	if !allowed[meta.ContentType] || (in.Kind != "poster" && meta.ContentType != "video/mp4") ||
		(in.Kind == "poster" && meta.ContentType == "video/mp4") {
		httpx.Error(w, 422, "invalid_uploaded_content", "uploaded content type does not match asset kind")
		return
	}
	limit := maxVideoUploadBytes
	if meta.ContentType != "video/mp4" { limit = maxImageUploadBytes }
	if meta.SizeBytes == 0 || meta.SizeBytes > limit {
		httpx.Error(w, 413, "invalid_upload_size", "uploaded object is empty or exceeds the permitted size")
		return
	}
	asset := Asset{Kind:in.Kind, StorageKey:in.StorageKey, MimeType:meta.ContentType, SizeBytes:meta.SizeBytes, IsActive:false}
	if in.TargetType == "movie" { asset.MovieID = &in.TargetID } else { asset.EpisodeID = &in.TargetID }
	if err := h.db.WithContext(r.Context()).Create(&asset).Error; err != nil {
		httpx.Error(w, 500, "internal_error", "asset could not be finalized")
		return
	}
	httpx.JSON(w, http.StatusCreated, map[string]any{
		"asset": asset,
		"status": "uploaded_pending_processing",
		"message": "source upload verified; asset remains inactive until HLS processing completes",
	})
}
