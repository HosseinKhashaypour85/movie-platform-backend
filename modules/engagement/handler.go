package engagement

import (
 "encoding/json"
 "net/http"
 "strconv"

 "movie-platform/internal/httpx"
 "movie-platform/internal/middleware"
 "github.com/go-chi/chi/v5"
)

type Handler struct { service *Service }
func NewHandler(s *Service)*Handler{return &Handler{service:s}}
func parseID(r *http.Request,key string)(uint,bool){v,e:=strconv.ParseUint(chi.URLParam(r,key),10,64);return uint(v),e==nil&&v>0}
func(h *Handler)Favorites(w http.ResponseWriter,r *http.Request){uid,_:=middleware.UserIDFromContext(r.Context());limit,_:=strconv.Atoi(r.URL.Query().Get("limit"));offset,_:=strconv.Atoi(r.URL.Query().Get("offset"));rows,e:=h.service.Favorites(uid,limit,offset);if e!=nil{httpx.Error(w,500,"internal_error","internal server error");return};httpx.JSON(w,200,map[string]any{"items":rows})}
func(h *Handler)AddFavorite(w http.ResponseWriter,r *http.Request){uid,_:=middleware.UserIDFromContext(r.Context());mid,ok:=parseID(r,"movieID");if !ok{httpx.Error(w,400,"invalid_movie_id","movie id is invalid");return};if e:=h.service.AddFavorite(uid,mid);e!=nil{httpx.Error(w,500,"internal_error","internal server error");return};httpx.JSON(w,201,map[string]string{"status":"added"})}
func(h *Handler)RemoveFavorite(w http.ResponseWriter,r *http.Request){uid,_:=middleware.UserIDFromContext(r.Context());mid,ok:=parseID(r,"movieID");if !ok{httpx.Error(w,400,"invalid_movie_id","movie id is invalid");return};if e:=h.service.RemoveFavorite(uid,mid);e!=nil{httpx.Error(w,500,"internal_error","internal server error");return};w.WriteHeader(http.StatusNoContent)}
func(h *Handler)Rate(w http.ResponseWriter,r *http.Request){uid,_:=middleware.UserIDFromContext(r.Context());mid,ok:=parseID(r,"movieID");if !ok{httpx.Error(w,400,"invalid_movie_id","movie id is invalid");return};var body struct{Score uint8 `json:"score"`};d:=json.NewDecoder(http.MaxBytesReader(w,r.Body,1024));d.DisallowUnknownFields();if d.Decode(&body)!=nil{httpx.Error(w,400,"invalid_json","invalid request body");return};if e:=h.service.Rate(uid,mid,body.Score);e!=nil{httpx.Error(w,400,"invalid_score","score must be between 1 and 10");return};httpx.JSON(w,200,map[string]string{"status":"rated"})}
func(h *Handler)RatingSummary(w http.ResponseWriter,r *http.Request){mid,ok:=parseID(r,"movieID");if !ok{httpx.Error(w,400,"invalid_movie_id","movie id is invalid");return};summary,e:=h.service.RatingSummary(mid);if e!=nil{httpx.Error(w,500,"internal_error","internal server error");return};httpx.JSON(w,200,summary)}
func RegisterRoutes(r chi.Router,h *Handler){r.Route("/favorites",func(r chi.Router){r.Get("/",h.Favorites);r.Put("/{movieID}",h.AddFavorite);r.Delete("/{movieID}",h.RemoveFavorite)});r.Route("/movies/{movieID}/ratings",func(r chi.Router){r.Post("/",h.Rate);r.Get("/summary",h.RatingSummary)})}
