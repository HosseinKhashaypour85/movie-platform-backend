package media

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"movie-platform/internal/httpx"
	"movie-platform/modules/catalog"
)

// Playback issues short-lived manifest URLs. The configured CDN/edge MUST validate
// the signature and expiration on every protected request (including HLS segments).
// An origin/CDN that ignores these parameters is not access controlled.
type Playback struct {
	db   *gorm.DB
	base *url.URL
	key  []byte
	ttl  time.Duration
}

func NewPlayback(db *gorm.DB) (*Playback, error) {
	raw := strings.TrimSpace(os.Getenv("PLAYBACK_CDN_BASE_URL"))
	secret := os.Getenv("PLAYBACK_SIGNING_KEY")
	if raw == "" && secret == "" {
		return &Playback{db: db}, nil
	}
	if raw == "" || len(secret) < 32 {
		return nil, errors.New("playback requires PLAYBACK_CDN_BASE_URL and PLAYBACK_SIGNING_KEY (at least 32 bytes)")
	}
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "https" || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || u.Path != "" {
		return nil, errors.New("PLAYBACK_CDN_BASE_URL must be an HTTPS origin without path, credentials or query")
	}
	return &Playback{db: db, base: u, key: []byte(secret), ttl: 5 * time.Minute}, nil
}

func (p *Playback) Available() bool {
	return p != nil && p.base != nil && len(p.key) >= 32
}

func validObjectKey(key string) bool {
	if key == "" || len(key) > 1024 || strings.HasPrefix(key, "/") || strings.Contains(key, "\\") {
		return false
	}
	for _, part := range strings.Split(key, "/") {
		if part == "" || part == "." || part == ".." {
			return false
		}
	}
	return true
}

func signature(key []byte, objectPath string, expires int64) string {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(objectPath))
	mac.Write([]byte("\n"))
	mac.Write([]byte(strconv.FormatInt(expires, 10)))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

// VerifyPlaybackSignature validates a signature bound to one exact object path.
func VerifyPlaybackSignature(key []byte, objectPath string, expires int64, supplied string, now time.Time) bool {
	if !validPlaybackExpiry(expires, now) { return false }
	return verifySignature(key, objectPath, expires, supplied)
}

// VerifyPlaybackScopeSignature authorizes a request only when its canonical path
// is inside the signed asset directory. A trusted CDN edge can exchange the
// signed master-manifest query for a short-lived cookie scoped to scopePath.
func VerifyPlaybackScopeSignature(key []byte, scopePath, requestPath string, expires int64, supplied string, now time.Time) bool {
	if !validPlaybackExpiry(expires, now) || !validPlaybackScope(scopePath) || !validObjectKey(strings.TrimPrefix(requestPath, "/")) {
		return false
	}
	if !strings.HasPrefix(requestPath, scopePath) { return false }
	return verifySignature(key, scopePath, expires, supplied)
}

func validPlaybackExpiry(expires int64, now time.Time) bool {
	return expires > now.Unix() && expires <= now.Add(5*time.Minute+30*time.Second).Unix()
}

func validPlaybackScope(scope string) bool {
	return strings.HasPrefix(scope, "/") && strings.HasSuffix(scope, "/") &&
		len(scope) > 2 && len(scope) <= 1024 &&
		path.Clean(strings.TrimSuffix(scope, "/"))+"/" == scope &&
		validObjectKey(strings.Trim(scope, "/"))
}

func verifySignature(key []byte, message string, expires int64, supplied string) bool {
	if len(key) < 32 { return false }
	actual, err := base64.RawURLEncoding.DecodeString(supplied)
	if err != nil { return false }
	expected, _ := base64.RawURLEncoding.DecodeString(signature(key, message, expires))
	return hmac.Equal(actual, expected)
}

func (p *Playback) URL(objectKey string, now time.Time) (string, time.Time, error) {
	if !p.Available() || !validObjectKey(objectKey) || !strings.HasSuffix(strings.ToLower(objectKey), ".m3u8") {
		return "", time.Time{}, errors.New("playback unavailable or invalid HLS manifest")
	}
	expires := now.UTC().Add(p.ttl)
	// Path is signed in its canonical, unescaped form. Edge verification must use
	// the same canonical path and reject alternate encodings.
	objectPath := path.Join("/", objectKey)
	u := *p.base
	u.Path = objectPath
	query := u.Query()
	query.Set("expires", strconv.FormatInt(expires.Unix(), 10))
	scopePath := path.Dir(objectPath) + "/"
	query.Set("scope", scopePath)
	query.Set("signature", signature(p.key, scopePath, expires.Unix()))
	u.RawQuery = query.Encode()
	return u.String(), expires, nil
}

func (p *Playback) Manifest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	if !p.Available() {
		httpx.Error(w, http.StatusServiceUnavailable, "playback_unavailable", "playback is not configured")
		return
	}
	kind := chi.URLParam(r, "kind")
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id == 0 || (kind != "movies" && kind != "episodes") {
		httpx.Error(w, http.StatusBadRequest, "invalid_target", "invalid playback target")
		return
	}
	var asset Asset
	q := p.db.WithContext(r.Context()).Where("kind = ? AND is_active = ?", "video", true)
	if kind == "movies" {
		var movie catalog.Movie
		if err := p.db.WithContext(r.Context()).Where("is_published = ?", true).First(&movie, id).Error; err != nil {
			p.notFoundOrError(w, err)
			return
		}
		q = q.Where("movie_id = ?", id)
	} else {
		var episode catalog.Episode
		err := p.db.WithContext(r.Context()).Table("episodes").
			Joins("JOIN seasons ON seasons.id = episodes.season_id").
			Joins("JOIN series ON series.id = seasons.series_id").
			Where("episodes.id = ? AND episodes.is_published = ? AND series.is_published = ? AND series.deleted_at IS NULL", id, true, true).
			Select("episodes.*").Take(&episode).Error
		if err != nil {
			p.notFoundOrError(w, err)
			return
		}
		q = q.Where("episode_id = ?", id)
	}
	if err := q.Order("id DESC").First(&asset).Error; err != nil {
		p.notFoundOrError(w, err)
		return
	}
	playURL, expires, err := p.URL(asset.StorageKey, time.Now())
	if err != nil {
		httpx.Error(w, http.StatusServiceUnavailable, "playback_unavailable", "HLS playback is not configured for this asset")
		return
	}
	httpx.JSON(w, http.StatusOK, map[string]any{
		"asset_id":     asset.ID,
		"format":       "hls",
		"manifest_url": playURL,
		"expires_at":   expires,
	})
}

func (p *Playback) notFoundOrError(w http.ResponseWriter, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		httpx.Error(w, http.StatusNotFound, "not_found", "published playable content not found")
		return
	}
	httpx.Error(w, http.StatusInternalServerError, "internal_error", "internal server error")
}

func RegisterPlaybackRoutes(r chi.Router, p *Playback) {
	r.Get("/playback/{kind}/{id}", p.Manifest)
}
