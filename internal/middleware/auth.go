package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"movie-platform/internal/httpx"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string
const UserIDKey contextKey="user_id"
type Claims struct{UserID uint `json:"user_id"`;Role string `json:"role"`;jwt.RegisteredClaims}
func Auth(secret []byte,issuer string)func(http.Handler)http.Handler{return func(next http.Handler)http.Handler{return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){p:=strings.Fields(r.Header.Get("Authorization"));if len(p)!=2||!strings.EqualFold(p[0],"Bearer"){httpx.Error(w,401,"unauthorized","valid bearer token required");return};c:=&Claims{};t,e:=jwt.ParseWithClaims(p[1],c,func(t *jwt.Token)(any,error){if t.Method!=jwt.SigningMethodHS256{return nil,jwt.ErrSignatureInvalid};return secret,nil},jwt.WithValidMethods([]string{"HS256"}),jwt.WithIssuer(issuer),jwt.WithExpirationRequired(),jwt.WithIssuedAt(),jwt.WithLeeway(30*time.Second));if e!=nil||t==nil||!t.Valid||c.UserID==0{httpx.Error(w,401,"unauthorized","invalid or expired token");return};ctx:=context.WithValue(r.Context(),UserIDKey,c.UserID);next.ServeHTTP(w,r.WithContext(ctx))})}}
func UserIDFromContext(ctx context.Context)(uint,bool){id,ok:=ctx.Value(UserIDKey).(uint);return id,ok&&id>0}
