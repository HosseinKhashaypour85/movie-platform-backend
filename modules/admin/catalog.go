package admin

import (
	"encoding/json"\n\t"strconv"
	"errors"
	"io"
	"net/http"
	"strings"

	"movie-platform/internal/httpx"
	"movie-platform/modules/catalog"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type movieInput struct {
	Title string `json:"title"`
	OriginalTitle string `json:"original_title"`
	Slug string `json:"slug"`
	Synopsis string `json:"synopsis"`
	ReleaseYear uint16 `json:"release_year"`
	DurationMinutes uint16 `json:"duration_minutes"`
	AgeRating string `json:"age_rating"`
	PosterURL string `json:"poster_url"`
	BackdropURL string `json:"backdrop_url"`
	Language string `json:"language"`
	Country string `json:"country"`
	IMDbRating float32 `json:"imdb_rating"`
	IsPublished bool `json:"is_published"`
}
type seriesInput struct {
	Title string `json:"title"`
	OriginalTitle string `json:"original_title"`
	Slug string `json:"slug"`
	Synopsis string `json:"synopsis"`
	ReleaseYear uint16 `json:"release_year"`
	EndYear uint16 `json:"end_year"`
	AgeRating string `json:"age_rating"`
	PosterURL string `json:"poster_url"`
	BackdropURL string `json:"backdrop_url"`
	Language string `json:"language"`
	Country string `json:"country"`
	IMDbRating float32 `json:"imdb_rating"`
	IsPublished bool `json:"is_published"`
}

func decodeAdminJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil { httpx.Error(w,400,"invalid_body","invalid request body"); return false }
	var extra any
	if err := dec.Decode(&extra); !errors.Is(err,io.EOF) { httpx.Error(w,400,"invalid_body","request must contain one JSON object"); return false }
	return true
}
func validCatalogFields(title, slug string, rating float32) bool {
	return strings.TrimSpace(title)!="" && len(title)<=240 && strings.TrimSpace(slug)!="" && len(slug)<=260 && rating>=0 && rating<=10
}
func adminID(r *http.Request) (uint,bool) {
	var id uint64
	parsed,err:=strconv.ParseUint(chi.URLParam(r,"id"),10,64); if err!=nil || parsed==0 { return 0,false }; id=parsed
	return uint(id),true
}

func (h *Handler) CreateMovie(w http.ResponseWriter,r *http.Request) {
	var in movieInput
	if !decodeAdminJSON(w,r,&in) { return }
	if !validCatalogFields(in.Title,in.Slug,in.IMDbRating) { httpx.Error(w,400,"validation_error","title, slug and IMDb rating are invalid"); return }
	m:=catalog.Movie{Title:strings.TrimSpace(in.Title),OriginalTitle:in.OriginalTitle,Slug:strings.TrimSpace(in.Slug),Synopsis:in.Synopsis,ReleaseYear:in.ReleaseYear,DurationMinutes:in.DurationMinutes,AgeRating:in.AgeRating,PosterURL:in.PosterURL,BackdropURL:in.BackdropURL,Language:in.Language,Country:in.Country,IMDbRating:in.IMDbRating,IsPublished:in.IsPublished}
	if err:=h.db.Create(&m).Error;err!=nil { httpx.Error(w,409,"conflict","movie could not be created; check slug uniqueness");return }
	httpx.JSON(w,201,m)
}
func (h *Handler) UpdateMovie(w http.ResponseWriter,r *http.Request) {
	id,ok:=adminID(r);if !ok {httpx.Error(w,400,"invalid_id","invalid movie id");return}
	var in movieInput;if !decodeAdminJSON(w,r,&in){return}
	if !validCatalogFields(in.Title,in.Slug,in.IMDbRating){httpx.Error(w,400,"validation_error","title, slug and IMDb rating are invalid");return}
	var m catalog.Movie
	if err:=h.db.First(&m,id).Error;err!=nil {if errors.Is(err,gorm.ErrRecordNotFound){httpx.Error(w,404,"not_found","movie not found")}else{httpx.Error(w,500,"internal_error","internal server error")};return}
	m.Title=strings.TrimSpace(in.Title);m.OriginalTitle=in.OriginalTitle;m.Slug=strings.TrimSpace(in.Slug);m.Synopsis=in.Synopsis;m.ReleaseYear=in.ReleaseYear;m.DurationMinutes=in.DurationMinutes;m.AgeRating=in.AgeRating;m.PosterURL=in.PosterURL;m.BackdropURL=in.BackdropURL;m.Language=in.Language;m.Country=in.Country;m.IMDbRating=in.IMDbRating;m.IsPublished=in.IsPublished
	if err:=h.db.Save(&m).Error;err!=nil {httpx.Error(w,409,"conflict","movie could not be updated; check slug uniqueness");return}
	httpx.JSON(w,200,m)
}
func (h *Handler) DeleteMovie(w http.ResponseWriter,r *http.Request) {
	id,ok:=adminID(r);if !ok {httpx.Error(w,400,"invalid_id","invalid movie id");return}
	res:=h.db.Delete(&catalog.Movie{},id);if res.Error!=nil {httpx.Error(w,500,"internal_error","internal server error");return};if res.RowsAffected==0 {httpx.Error(w,404,"not_found","movie not found");return}
	w.WriteHeader(http.StatusNoContent)
}
func (h *Handler) CreateSeries(w http.ResponseWriter,r *http.Request) {
	var in seriesInput;if !decodeAdminJSON(w,r,&in){return}
	if !validCatalogFields(in.Title,in.Slug,in.IMDbRating){httpx.Error(w,400,"validation_error","title, slug and IMDb rating are invalid");return}
	m:=catalog.Series{Title:strings.TrimSpace(in.Title),OriginalTitle:in.OriginalTitle,Slug:strings.TrimSpace(in.Slug),Synopsis:in.Synopsis,ReleaseYear:in.ReleaseYear,EndYear:in.EndYear,AgeRating:in.AgeRating,PosterURL:in.PosterURL,BackdropURL:in.BackdropURL,Language:in.Language,Country:in.Country,IMDbRating:in.IMDbRating,IsPublished:in.IsPublished}
	if err:=h.db.Create(&m).Error;err!=nil {httpx.Error(w,409,"conflict","series could not be created; check slug uniqueness");return}
	httpx.JSON(w,201,m)
}
func (h *Handler) UpdateSeries(w http.ResponseWriter,r *http.Request) {
	id,ok:=adminID(r);if !ok {httpx.Error(w,400,"invalid_id","invalid series id");return}
	var in seriesInput;if !decodeAdminJSON(w,r,&in){return}
	if !validCatalogFields(in.Title,in.Slug,in.IMDbRating){httpx.Error(w,400,"validation_error","title, slug and IMDb rating are invalid");return}
	var m catalog.Series
	if err:=h.db.First(&m,id).Error;err!=nil {if errors.Is(err,gorm.ErrRecordNotFound){httpx.Error(w,404,"not_found","series not found")}else{httpx.Error(w,500,"internal_error","internal server error")};return}
	m.Title=strings.TrimSpace(in.Title);m.OriginalTitle=in.OriginalTitle;m.Slug=strings.TrimSpace(in.Slug);m.Synopsis=in.Synopsis;m.ReleaseYear=in.ReleaseYear;m.EndYear=in.EndYear;m.AgeRating=in.AgeRating;m.PosterURL=in.PosterURL;m.BackdropURL=in.BackdropURL;m.Language=in.Language;m.Country=in.Country;m.IMDbRating=in.IMDbRating;m.IsPublished=in.IsPublished
	if err:=h.db.Save(&m).Error;err!=nil {httpx.Error(w,409,"conflict","series could not be updated; check slug uniqueness");return}
	httpx.JSON(w,200,m)
}
func (h *Handler) DeleteSeries(w http.ResponseWriter,r *http.Request) {
	id,ok:=adminID(r);if !ok {httpx.Error(w,400,"invalid_id","invalid series id");return}
	res:=h.db.Delete(&catalog.Series{},id);if res.Error!=nil {httpx.Error(w,500,"internal_error","internal server error");return};if res.RowsAffected==0 {httpx.Error(w,404,"not_found","series not found");return}
	w.WriteHeader(http.StatusNoContent)
}
