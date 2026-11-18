package middleware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"movie-platform/internal/middleware"

	"github.com/golang-jwt/jwt/v5"
)

func TestAuthMiddleware(t *testing.T) {
	secret := []byte("12345678901234567890123456789012")
	issuer := "movie-platform"

	now := time.Now()
	claims := jwt.MapClaims{
		"user_id": uint(42),
		"role":    "admin",
		"iss":     issuer,
		"sub":     fmt.Sprint(42),
		"iat":     now.Unix(),
		"exp":     now.Add(15 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	handler := middleware.Auth(secret, issuer)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uid, ok := middleware.UserIDFromContext(r.Context())
		if !ok || uid != 42 {
			t.Errorf("expected uid 42, got %v (ok=%v)", uid, ok)
		}
		role, ok := middleware.RoleFromContext(r.Context())
		if !ok || role != "admin" {
			t.Errorf("expected role 'admin', got '%v' (ok=%v)", role, ok)
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d, body: %s", rec.Code, rec.Body.String())
	}
}
