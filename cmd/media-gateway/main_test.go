package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"movie-platform/modules/media"
)

func newTestGateway(t *testing.T, origin *httptest.Server) *gateway {
	t.Helper()
	t.Setenv("S3_ENDPOINT", origin.URL)
	t.Setenv("S3_BUCKET", "private-media")
	t.Setenv("S3_REGION", "eu-central-1")
	t.Setenv("S3_ACCESS_KEY_ID", "TESTACCESS")
	t.Setenv("S3_SECRET_ACCESS_KEY", strings.Repeat("s", 40))

	storage, err := media.NewS3PresignerFromEnv()
	if err != nil {
		t.Fatalf("create S3 presigner: %v", err)
	}
	originClient := origin.Client()
	originClient.CheckRedirect = func(*http.Request, []*http.Request) error {
		return fmt.Errorf("redirects are disabled")
	}
	return &gateway{storage: storage, client: originClient}
}

func TestGatewayStreamsPrivateObjectAndForwardsRange(t *testing.T) {
	const payload = "0123456789"
	origin := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/private-media/hls/42/master.m3u8" {
			t.Errorf("unexpected private origin path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("origin method = %s, want GET", r.Method)
		}
		if r.Header.Get("Range") != "bytes=2-5" {
			t.Errorf("origin Range = %q", r.Header.Get("Range"))
		}
		w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
		w.Header().Set("Accept-Ranges", "bytes")
		w.Header().Set("Content-Range", "bytes 2-5/10")
		w.Header().Set("ETag", "\"test-etag\"")
		w.WriteHeader(http.StatusPartialContent)
		_, _ = io.WriteString(w, payload[2:6])
	}))
	defer origin.Close()

	g := newTestGateway(t, origin)
	req := httptest.NewRequest(http.MethodGet, "/hls/42/master.m3u8", nil)
	req.Header.Set("Range", "bytes=2-5")
	rec := httptest.NewRecorder()
	g.serveObject(rec, req)

	if rec.Code != http.StatusPartialContent {
		t.Fatalf("status = %d, want 206: %s", rec.Code, rec.Body.String())
	}
	if rec.Body.String() != "2345" {
		t.Fatalf("body = %q, want ranged payload", rec.Body.String())
	}
	if rec.Header().Get("Content-Range") != "bytes 2-5/10" {
		t.Fatalf("Content-Range = %q", rec.Header().Get("Content-Range"))
	}
	if rec.Header().Get("ETag") != "\"test-etag\"" {
		t.Fatal("ETag was not forwarded")
	}
	if strings.Contains(rec.Body.String(), origin.URL) {
		t.Fatal("private origin URL leaked to the client")
	}
}

func TestGatewayUsesUpstreamHEADWithoutBody(t *testing.T) {
	bodyWrites := 0
	origin := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodHead {
			t.Errorf("origin method = %s, want HEAD", r.Method)
		}
		w.Header().Set("Content-Type", "video/mp2t")
		w.Header().Set("Content-Length", "1234")
		w.Header().Set("Accept-Ranges", "bytes")
		bodyWrites++
	}))
	defer origin.Close()

	g := newTestGateway(t, origin)
	req := httptest.NewRequest(http.MethodHead, "/hls/42/segment.ts", nil)
	rec := httptest.NewRecorder()
	g.serveObject(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if rec.Body.Len() != 0 {
		t.Fatalf("HEAD returned %d body bytes", rec.Body.Len())
	}
	if rec.Header().Get("Content-Type") != "video/mp2t" {
		t.Fatal("content type missing from HEAD response")
	}
	if bodyWrites != 1 {
		t.Fatalf("origin request count = %d, want 1", bodyWrites)
	}
}

func TestGatewayMapsMissingObjectsAndRejectsInvalidKeys(t *testing.T) {
	origin := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	defer origin.Close()

	g := newTestGateway(t, origin)
	req := httptest.NewRequest(http.MethodGet, "/hls/42/missing.ts", nil)
	rec := httptest.NewRecorder()
	g.serveObject(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("missing object status = %d, want 404", rec.Code)
	}

	invalid := httptest.NewRequest(http.MethodGet, "/../escape.ts", nil)
	invalidRec := httptest.NewRecorder()
	g.serveObject(invalidRec, invalid)
	if invalidRec.Code != http.StatusNotFound {
		t.Fatalf("invalid key status = %d, want 404", invalidRec.Code)
	}
}

func TestGatewayDoesNotFollowOriginRedirects(t *testing.T) {
	var redirect *httptest.Server
	redirect = httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://example.invalid/private-object", http.StatusFound)
	}))
	defer redirect.Close()

	g := newTestGateway(t, redirect)
	req := httptest.NewRequest(http.MethodGet, "/hls/42/master.m3u8", nil)
	rec := httptest.NewRecorder()
	g.serveObject(rec, req)
	if rec.Code != http.StatusBadGateway {
		t.Fatalf("redirect response status = %d, want 502", rec.Code)
	}
	if strings.Contains(rec.Body.String(), "example.invalid") {
		t.Fatal("origin redirect target leaked to client")
	}
	_ = time.Second // keep test time-independent; signed URLs use current UTC time.
}
