package media

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHLSCookieMiddlewareAuthorizesScopedSegments(t *testing.T) {
	key := []byte(strings.Repeat("k", 32))
	now := time.Now()
	expires := now.Add(5 * time.Minute).Unix()
	scope := "/hls/asset-42/"
	sig := signature(key, scope, expires)

	handler := HLSCookieMiddleware(EdgeAuthConfig{SigningKey: key})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	master := httptest.NewRequest(http.MethodGet, "/hls/asset-42/master.m3u8?scope=%2Fhls%2Fasset-42%2F&expires="+
		strings.TrimSpace(stringInt(expires))+"&signature="+sig, nil)
	masterResult := httptest.NewRecorder()
	handler.ServeHTTP(masterResult, master)
	if masterResult.Code != http.StatusNoContent {
		t.Fatalf("master request status = %d", masterResult.Code)
	}
	cookies := masterResult.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected one playback cookie, got %d", len(cookies))
	}
	cookie := cookies[0]
	if !cookie.HttpOnly || !cookie.Secure || cookie.Path != scope {
		t.Fatalf("cookie flags/path are not secure: %+v", cookie)
	}

	segment := httptest.NewRequest(http.MethodGet, "/hls/asset-42/720p/segment_00001.ts", nil)
	segment.AddCookie(cookie)
	segmentResult := httptest.NewRecorder()
	handler.ServeHTTP(segmentResult, segment)
	if segmentResult.Code != http.StatusNoContent {
		t.Fatalf("valid segment status = %d", segmentResult.Code)
	}

	otherAsset := httptest.NewRequest(http.MethodGet, "/hls/asset-420/segment.ts", nil)
	otherAsset.AddCookie(cookie)
	otherResult := httptest.NewRecorder()
	handler.ServeHTTP(otherResult, otherAsset)
	if otherResult.Code != http.StatusForbidden {
		t.Fatalf("cross-asset request status = %d, want 403", otherResult.Code)
	}
}

func TestHLSCookieMiddlewareRejectsMissingExpiredAndNonCanonicalRequests(t *testing.T) {
	key := []byte(strings.Repeat("k", 32))
	handler := HLSCookieMiddleware(EdgeAuthConfig{SigningKey: key})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	for _, target := range []string{
		"/hls/asset-42/segment.ts",
		"/hls//asset-42/segment.ts",
		"/hls/asset-42/../asset-43/segment.ts",
		"/hls/asset-420/segment.ts?scope=%2Fhls%2Fasset-42%2F&expires=1&signature=bad",
	} {
		req := httptest.NewRequest(http.MethodGet, target, nil)
		res := httptest.NewRecorder()
		handler.ServeHTTP(res, req)
		if res.Code != http.StatusForbidden {
			t.Errorf("%s: status = %d, want 403", target, res.Code)
		}
	}
}

func stringInt(v int64) string {
	return strconv.FormatInt(v, 10)
}
