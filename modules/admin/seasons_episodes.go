package admin

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"movie-platform/internal/httpx"
	"movie-platform/modules/catalog"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type seasonInput struct {
	SeriesID uint `json:"series_id"`
	Number uint16 `json:"number"`
	Title string `json:"title"`
}
type episodeInput struct {
	SeasonID uint `json:"season_id"`
	Number uint16 `json:"number"`
	Title string `json:"title"`
	Synopsis string `json:"synopsis"`
	DurationMinutes uint16 `json:"duration_minutes"`
	ThumbnailURL string `json:"thumbnail_url"`
	IsPublished bool `json:"is_published"`
}
func routeUint(r *http.Request,key string)(uint,bool){
	v,err:=strconv.ParseUint(chi.URLParam(r,key),10,64)
	return uint(v),err==nil && v>0
}
func (h *Handler) CreateSeason(w http.ResponseWriter,r *http.Request){
	var in seasonInput
	if !decodeAdminJSON(w,r,&in){return}
	if in.SeriesID==0 || in.Number==0 || len(in.Title)>240 {httpx.Error(w,400,"validation_error","series_id and positive season number are required");return}
	var parent catalog.Series
	if err:=h.db.First(&parent,in.SeriesID).Error;err!=nil {if errors.Is(err,gorm.ErrRecordNotFound){httpx.Error(w,404,"not_found","series not found")}else{httpx.Error(w,500,"internal_error","internal server error")};return}
	row:=catalog.Season{SeriesID:in.SeriesID,Number:in.Number,Title:strings.TrimSpace(in.Title)}
	if err:=h.db.Create(&row).Error;err!=nil {httpx.Error(w,409,"conflict","season could not be created; check season number uniqueness");return}
	httpx.JSON(w,201,row)
}
func (h *Handler) UpdateSeason(w http.ResponseWriter,r *http.Request){
	id,ok:=routeUint(r,"id");if !ok {httpx.Error(w,400,"invalid_id","invalid season id");return}
	var in seasonInput;if !decodeAdminJSON(w,r,&in){return}
	if in.SeriesID==0 || in.Number==0 || len(in.Title)>240 {httpx.Error(w,400,"validation_error","series_id and positive season number are required");return}
	var parent catalog.Series
	if err:=h.db.First(&parent,in.SeriesID).Error;err!=nil {if errors.Is(err,gorm.ErrRecordNotFound){httpx.Error(w,404,"not_found","series not found")}else{httpx.Error(w,500,"internal_error","internal server error")};return}
	var row catalog.Season
	if err:=h.db.First(&row,id).Error;err!=nil {if errors.Is(err,gorm.ErrRecordNotFound){httpx.Error(w,404,"not_found","season not found")}else{httpx.Error(w,500,"internal_error","internal server error")};return}
	row.SeriesID=in.SeriesID;row.Number=in.Number;row.Title=strings.TrimSpace(in.Title)
	if err:=h.db.Save(&row).Error;err!=nil {httpx.Error(w,409,"conflict","season could not be updated; check season number uniqueness");return}
	httpx.JSON(w,200,row)
}
func (h *Handler) DeleteSeason(w http.ResponseWriter,r *http.Request){
	id,ok:=routeUint(r,"id");if !ok {httpx.Error(w,400,"invalid_id","invalid season id");return}
	var count int64
	if err:=h.db.Model(&catalog.Episode{}).Where("season_id = ?",id).Count(&count).Error;err!=nil {httpx.Error(w,500,"internal_error","internal server error");return}
	if count>0 {httpx.Error(w,409,"has_episodes","delete or move episodes before deleting this season");return}
	res:=h.db.Delete(&catalog.Season{},id);if res.Error!=nil {httpx.Error(w,500,"internal_error","internal server error");return};if res.RowsAffected==0 {httpx.Error(w,404,"not_found","season not found");return}
	w.WriteHeader(http.StatusNoContent)
}
func (h *Handler) CreateEpisode(w http.ResponseWriter,r *http.Request){
	var in episodeInput;if !decodeAdminJSON(w,r,&in){return}
	if in.SeasonID==0 || in.Number==0 || strings.TrimSpace(in.Title)=="" || len(in.Title)>240 || len(in.ThumbnailURL)>2048 {httpx.Error(w,400,"validation_error","season_id, positive episode number and title are required");return}
	var parent catalog.Season
	if err:=h.db.First(&parent,in.SeasonID).Error;err!=nil {if errors.Is(err,gorm.ErrRecordNotFound){httpx.Error(w,404,"not_found","season not found")}else{httpx.Error(w,500,"internal_error","internal server error")};return}
	row:=catalog.Episode{SeasonID:in.SeasonID,Number:in.Number,Title:strings.TrimSpace(in.Title),Synopsis:in.Synopsis,DurationMinutes:in.DurationMinutes,ThumbnailURL:in.ThumbnailURL,IsPublished:in.IsPublished}
	if err:=h.db.Create(&row).Error;err!=nil {httpx.Error(w,409,"conflict","episode could not be created");return}
	httpx.JSON(w,201,row)
}
func (h *Handler) UpdateEpisode(w http.ResponseWriter,r *http.Request){
	id,ok:=routeUint(r,"id");if !ok {httpx.Error(w,400,"invalid_id","invalid episode id");return}
	var in episodeInput;if !decodeAdminJSON(w,r,&in){return}
	if in.SeasonID==0 || in.Number==0 || strings.TrimSpace(in.Title)=="" || len(in.Title)>240 || len(in.ThumbnailURL)>2048 {httpx.Error(w,400,"validation_error","season_id, positive episode number and title are required");return}
	var parent catalog.Season
	if err:=h.db.First(&parent,in.SeasonID).Error;err!=nil {if errors.Is(err,gorm.ErrRecordNotFound){httpx.Error(w,404,"not_found","season not found")}else{httpx.Error(w,500,"internal_error","internal server error")};return}
	var row catalog.Episode
	if err:=h.db.First(&row,id).Error;err!=nil {if errors.Is(err,gorm.ErrRecordNotFound){httpx.Error(w,404,"not_found","episode not found")}else{httpx.Error(w,500,"internal_error","internal server error")};return}
	row.SeasonID=in.SeasonID;row.Number=in.Number;row.Title=strings.TrimSpace(in.Title);row.Synopsis=in.Synopsis;row.DurationMinutes=in.DurationMinutes;row.ThumbnailURL=in.ThumbnailURL;row.IsPublished=in.IsPublished
	if err:=h.db.Save(&row).Error;err!=nil {httpx.Error(w,409,"conflict","episode could not be updated");return}
	httpx.JSON(w,200,row)
}
func (h *Handler) DeleteEpisode(w http.ResponseWriter,r *http.Request){
	id,ok:=routeUint(r,"id");if !ok {httpx.Error(w,400,"invalid_id","invalid episode id");return}
	res:=h.db.Delete(&catalog.Episode{},id);if res.Error!=nil {httpx.Error(w,500,"internal_error","internal server error");return};if res.RowsAffected==0 {httpx.Error(w,404,"not_found","episode not found");return}
	w.WriteHeader(http.StatusNoContent)
}
