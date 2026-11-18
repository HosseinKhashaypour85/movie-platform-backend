package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"

	"movie-platform/internal/httpx"
)

type requestIDKey struct{}
func RequestID(next http.Handler)http.Handler{return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){b:=make([]byte,16);if _,err:=rand.Read(b);err!=nil{httpx.Error(w,500,"internal_error","internal server error");return};id:=hex.EncodeToString(b);w.Header().Set("X-Request-ID",id);next.ServeHTTP(w,r.WithContext(context.WithValue(r.Context(),requestIDKey{},id)))})}
func Recover(l *slog.Logger)func(http.Handler)http.Handler{return func(next http.Handler)http.Handler{return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){defer func(){if recover()!=nil{l.Error("panic recovered","request_id",r.Header.Get("X-Request-ID"));httpx.Error(w,500,"internal_error","internal server error")}}();next.ServeHTTP(w,r)})}}
func SecurityHeaders(next http.Handler)http.Handler{return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){w.Header().Set("X-Content-Type-Options","nosniff");w.Header().Set("X-Frame-Options","DENY");w.Header().Set("Referrer-Policy","no-referrer");w.Header().Set("Permissions-Policy","camera=(), microphone=(), geolocation=()");w.Header().Set("Cache-Control","no-store");w.Header().Set("Content-Security-Policy","default-src 'none'; frame-ancestors 'none'");next.ServeHTTP(w,r)})}
func BodyLimit(n int64)func(http.Handler)http.Handler{return func(next http.Handler)http.Handler{return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){r.Body=http.MaxBytesReader(w,r.Body,n);next.ServeHTTP(w,r)})}}
func Timeout(d time.Duration)func(http.Handler)http.Handler{return func(next http.Handler)http.Handler{return http.TimeoutHandler(next,d,`{"error":{"code":"timeout","message":"request timed out"}}`)}}
func CORS(origins []string)func(http.Handler)http.Handler{allowed:=map[string]bool{};for _,o:=range origins{allowed[o]=true};return func(next http.Handler)http.Handler{return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){o:=r.Header.Get("Origin");if o!=""{if !allowed[o]{httpx.Error(w,403,"origin_denied","origin is not allowed");return};w.Header().Set("Access-Control-Allow-Origin",o);w.Header().Set("Vary","Origin");w.Header().Set("Access-Control-Allow-Headers","Authorization, Content-Type, X-Request-ID");w.Header().Set("Access-Control-Allow-Methods","GET, POST, PUT, DELETE, OPTIONS");w.Header().Set("Access-Control-Max-Age","300")};if r.Method==http.MethodOptions{w.WriteHeader(204);return};next.ServeHTTP(w,r)})}}
type rateBucket struct{tokens float64;last time.Time}
type rateLimiter struct{mu sync.Mutex;buckets map[string]*rateBucket;rate,burst float64}
func RateLimit(count int,window time.Duration)func(http.Handler)http.Handler{l:=&rateLimiter{buckets:map[string]*rateBucket{},rate:float64(count)/window.Seconds(),burst:float64(count)};return func(next http.Handler)http.Handler{return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){ip,_,e:=net.SplitHostPort(r.RemoteAddr);if e!=nil{ip=r.RemoteAddr};now:=time.Now();l.mu.Lock();b:=l.buckets[ip];if b==nil{b=&rateBucket{tokens:l.burst,last:now};l.buckets[ip]=b};b.tokens+=now.Sub(b.last).Seconds()*l.rate;if b.tokens>l.burst{b.tokens=l.burst};b.last=now;ok:=b.tokens>=1;if ok{b.tokens--};l.mu.Unlock();if !ok{w.Header().Set("Retry-After","60");httpx.Error(w,429,"rate_limited","too many requests");return};next.ServeHTTP(w,r)})}}
