package media

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"movie-platform/internal/httpx"
	"movie-platform/internal/middleware"
)

type Handler struct { db *gorm.DB; workerToken string }
func NewHandler(db *gorm.DB) *Handler { return &Handler{db:db,workerToken:os.Getenv("SUBTITLE_WORKER_TOKEN")} }

func (h *Handler) PublicSubtitles(w http.ResponseWriter,r *http.Request) {
	kind:=chi.URLParam(r,"kind"); id,e:=strconv.ParseUint(chi.URLParam(r,"id"),10,64)
	if e!=nil || id==0 || (kind!="movies" && kind!="episodes") { httpx.Error(w,400,"invalid_target","invalid media target"); return }
	q:=h.db.Model(&Subtitle{}).Where("is_published = ?",true)
	if kind=="movies" { q=q.Where("movie_id = ?",id) } else { q=q.Where("episode_id = ?",id) }
	var rows []Subtitle
	if err:=q.Select("id","movie_id","episode_id","language","label","format","download_url","is_generated","created_at").Order("language ASC").Find(&rows).Error; err!=nil { httpx.Error(w,500,"internal_error","internal server error"); return }
	httpx.JSON(w,200,map[string]any{"items":rows})
}

func (h *Handler) CreateAsset(w http.ResponseWriter,r *http.Request) {
	var in Asset
	if !decodeMedia(w,r,&in) { return }
	if (in.MovieID==nil)==(in.EpisodeID==nil) || (in.Kind!="video" && in.Kind!="trailer" && in.Kind!="poster") || strings.TrimSpace(in.StorageKey)=="" || len(in.StorageKey)>1024 || (in.PlaybackURL!="" && !httpsURL(in.PlaybackURL)) || len(in.MimeType)>120 {
		httpx.Error(w,400,"validation_error","asset must target exactly one movie or episode and include valid metadata"); return
	}
	in.ID=0
	if err:=h.db.Create(&in).Error; err!=nil { httpx.Error(w,500,"internal_error","asset could not be created"); return }
	httpx.JSON(w,201,in)
}
func (h *Handler) CreateSubtitleJob(w http.ResponseWriter,r *http.Request) {
	var in struct { SourceAssetID uint `json:"source_asset_id"`; SourceLanguage string `json:"source_language"`; TargetLanguage string `json:"target_language"` }
	if !decodeMedia(w,r,&in) { return }
	if in.SourceAssetID==0 || len(in.SourceLanguage)>12 || len(in.TargetLanguage)>12 || in.TargetLanguage!="fa" || in.SourceLanguage=="" { httpx.Error(w,400,"validation_error","valid source asset and target language fa are required"); return }
	var asset Asset
	if err:=h.db.First(&asset,in.SourceAssetID).Error; err!=nil || asset.Kind!="video" { httpx.Error(w,404,"source_not_found","source video asset not found"); return }
	job:=SubtitleJob{MovieID:asset.MovieID,EpisodeID:asset.EpisodeID,SourceAssetID:asset.ID,SourceLanguage:in.SourceLanguage,TargetLanguage:in.TargetLanguage,Status:"queued"}
	if err:=h.db.Create(&job).Error; err!=nil { httpx.Error(w,500,"internal_error","job could not be queued"); return }
	httpx.JSON(w,202,job)
}

func (h *Handler) ClaimJob(w http.ResponseWriter,r *http.Request) {
	if !h.authorizeWorker(w,r) { return }
	var job SubtitleJob
	now:=time.Now().UTC()
	err:=h.db.Transaction(func(tx *gorm.DB) error {
		if err:=tx.Clauses(clause.Locking{Strength:"UPDATE",Options:"SKIP LOCKED"}).Where("status = ? OR (status = ? AND locked_at < ?)", "queued","processing",now.Add(-15*time.Minute)).Order("id ASC").First(&job).Error; err!=nil { return err }
		return tx.Model(&job).Updates(map[string]any{"status":"processing","locked_at":now,"attempts":gorm.Expr("attempts + 1")}).Error
	})
	if errors.Is(err,gorm.ErrRecordNotFound) { httpx.JSON(w,200,map[string]any{"job":nil}); return }
	if err!=nil { httpx.Error(w,500,"internal_error","job claim failed"); return }
	httpx.JSON(w,200,map[string]any{"job":job})
}

func (h *Handler) CompleteJob(w http.ResponseWriter,r *http.Request) {
	if !h.authorizeWorker(w,r) { return }
	id,e:=strconv.ParseUint(chi.URLParam(r,"id"),10,64); if e!=nil || id==0 { httpx.Error(w,400,"invalid_id","invalid job id"); return }
	var in struct { StorageKey string `json:"storage_key"`; DownloadURL string `json:"download_url"`; Label string `json:"label"`; Format string `json:"format"`; Error string `json:"error"` }
	if !decodeMedia(w,r,&in) { return }
	var job SubtitleJob
	if err:=h.db.First(&job,id).Error; err!=nil { httpx.Error(w,404,"not_found","job not found"); return }
	if job.Status!="processing" { httpx.Error(w,409,"invalid_state","job is not processing"); return }
	if in.Error!="" {
		if len(in.Error)>1000 { in.Error=in.Error[:1000] }
		h.db.Model(&job).Updates(map[string]any{"status":"failed","error_message":in.Error,"locked_at":nil})
		httpx.JSON(w,200,map[string]string{"status":"failed"}); return
	}
	if strings.TrimSpace(in.StorageKey)=="" || len(in.StorageKey)>1024 || (in.DownloadURL!="" && !httpsURL(in.DownloadURL)) || (in.Format!="vtt" && in.Format!="srt") || strings.TrimSpace(in.Label)=="" || len(in.Label)>80 {
		httpx.Error(w,400,"validation_error","subtitle output metadata is invalid"); return
	}
	sub:=Subtitle{MovieID:job.MovieID,EpisodeID:job.EpisodeID,Language:job.TargetLanguage,Label:in.Label,Format:in.Format,StorageKey:in.StorageKey,DownloadURL:in.DownloadURL,IsGenerated:true,IsPublished:false}
	if err:=h.db.Transaction(func(tx *gorm.DB) error {
		if err:=tx.Create(&sub).Error; err!=nil { return err }
		return tx.Model(&job).Updates(map[string]any{"status":"completed","output_subtitle_id":sub.ID,"locked_at":nil,"error_message":""}).Error
	}); err!=nil { httpx.Error(w,500,"internal_error","subtitle output could not be saved"); return }
	httpx.JSON(w,201,map[string]any{"subtitle":sub,"status":"completed"})
}

func (h *Handler) PublishSubtitle(w http.ResponseWriter,r *http.Request) {
	id,e:=strconv.ParseUint(chi.URLParam(r,"id"),10,64); if e!=nil || id==0 { httpx.Error(w,400,"invalid_id","invalid subtitle id"); return }
	var body struct { IsPublished *bool `json:"is_published"` }
	if !decodeMedia(w,r,&body) || body.IsPublished==nil { httpx.Error(w,400,"validation_error","is_published is required"); return }
	res:=h.db.Model(&Subtitle{}).Where("id = ?",id).Update("is_published",*body.IsPublished)
	if res.Error!=nil { httpx.Error(w,500,"internal_error","internal server error"); return }
	if res.RowsAffected==0 { httpx.Error(w,404,"not_found","subtitle not found"); return }
	httpx.JSON(w,200,map[string]bool{"is_published":*body.IsPublished})
}

func (h *Handler) authorizeWorker(w http.ResponseWriter,r *http.Request) bool {
	if h.workerToken=="" { httpx.Error(w,503,"worker_disabled","subtitle worker is not configured"); return false }
	if r.Header.Get("Authorization")!="Bearer "+h.workerToken { httpx.Error(w,401,"unauthorized","worker token required"); return false }
	return true
}
func decodeMedia(w http.ResponseWriter,r *http.Request,dst any) bool {
	r.Body=http.MaxBytesReader(w,r.Body,1<<20); d:=json.NewDecoder(r.Body); d.DisallowUnknownFields()
	if err:=d.Decode(dst); err!=nil { httpx.Error(w,400,"invalid_body","invalid request body"); return false }
	var extra any; if err:=d.Decode(&extra); err!=io.EOF { httpx.Error(w,400,"invalid_body","request must contain one JSON object"); return false }
	return true
}
func httpsURL(v string) bool { return strings.HasPrefix(v,"https://") && len(v)<=2048 }

func RegisterPublicRoutes(r chi.Router,h *Handler) {
	r.Get("/media/{kind}/{id}/subtitles",h.PublicSubtitles)
}
func RegisterAdminRoutes(r chi.Router,h *Handler) {
	r.Route("/media",func(a chi.Router) {
		a.Use(middleware.RequireRole("admin","super_admin"))
		a.Post("/assets",h.CreateAsset)
		a.Post("/subtitle-jobs",h.CreateSubtitleJob)
		a.Patch("/subtitles/{id}/publish",h.PublishSubtitle)
	})
}
func RegisterWorkerRoutes(r chi.Router,h *Handler) {
	r.Post("/worker/subtitle-jobs/claim",h.ClaimJob)
	r.Post("/worker/subtitle-jobs/{id}/complete",h.CompleteJob)
}
