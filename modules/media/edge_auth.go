package media

import (
	"encoding/base64"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
)

// EdgeAuthConfig configures the provider-neutral HLS authorization middleware.
// Mount it on the media origin or a self-hosted reverse proxy in front of
// private object storage. CDN vendors can implement the same contract at edge.
type EdgeAuthConfig struct {
	SigningKey []byte
	CookieName string
}

// HLSCookieMiddleware exchanges a signed master-manifest URL for a short-lived,
// path-scoped HttpOnly cookie. Subsequent playlists, segments and keys must pass
// through this middleware (or an equivalent CDN edge implementation).
func HLSCookieMiddleware(cfg EdgeAuthConfig) func(http.Handler) http.Handler {
	cookieName := cfg.CookieName
	if cookieName == "" {
		cookieName = "hls_playback"
	}
	key := append([]byte(nil), cfg.SigningKey...)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestPath, ok := canonicalHLSRequestPath(r)
			if !ok || len(key) < 32 {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			if expiresRaw := r.URL.Query().Get("expires"); expiresRaw != "" {
				expires, err := strconv.ParseInt(expiresRaw, 10, 64)
				scope := r.URL.Query().Get("scope")
				sig := r.URL.Query().Get("signature")
				if err != nil || !VerifyPlaybackScopeSignature(key, scope, requestPath, expires, sig, time.Now()) {
					http.Error(w, "forbidden", http.StatusForbidden)
					return
				}
				remaining := int(time.Until(time.Unix(expires, 0)).Seconds())
				if remaining < 1 {
					http.Error(w, "forbidden", http.StatusForbidden)
					return
				}
				value := encodeHLSCookie(scope, expires, sig)
				http.SetCookie(w, &http.Cookie{
					Name: cookieName, Value: value, Path: scope,
					MaxAge: remaining, Expires: time.Unix(expires, 0).UTC(),
					Secure: true, HttpOnly: true, SameSite: http.SameSiteLaxMode,
				})
				next.ServeHTTP(w, r)
				return
			}

			cookie, err := r.Cookie(cookieName)
			if err != nil {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			scope, expires, sig, ok := decodeHLSCookie(cookie.Value)
			if !ok || !VerifyPlaybackScopeSignature(key, scope, requestPath, expires, sig, time.Now()) {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func canonicalHLSRequestPath(r *http.Request) (string, bool) {
	if r.URL == nil || r.URL.RawPath != "" {
		return "", false
	}
	p := r.URL.Path
	if !strings.HasPrefix(p, "/") || strings.HasPrefix(p, "//") || strings.Contains(p, "\\") ||
		path.Clean(p) != p || !validObjectKey(strings.TrimPrefix(p, "/")) {
		return "", false
	}
	return p, true
}

func encodeHLSCookie(scope string, expires int64, signature string) string {
	payload := scope + "\n" + strconv.FormatInt(expires, 10) + "\n" + signature
	return base64.RawURLEncoding.EncodeToString([]byte(payload))
}

func decodeHLSCookie(value string) (string, int64, string, bool) {
	raw, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return "", 0, "", false
	}
	parts := strings.Split(string(raw), "\n")
	if len(parts) != 3 || parts[0] == "" || parts[2] == "" {
		return "", 0, "", false
	}
	expires, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return "", 0, "", false
	}
	return parts[0], expires, parts[2], true
}
