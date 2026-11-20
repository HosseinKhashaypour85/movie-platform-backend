package media

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestValidObjectKey(t *testing.T) {
	valid := []string{"movies/12/video/2026/11/file.mp4", "episodes/9/poster/image.webp"}
	for _, key := range valid { if !validObjectKey(key) { t.Errorf("rejected valid key %q", key) } }
	invalid := []string{"", "../secret", "movies/../secret", "/absolute", "movies//bad", "movies\\bad", "movies/./bad"}
	for _, key := range invalid { if validObjectKey(key) { t.Errorf("accepted invalid key %q", key) } }
}

func TestHeadObjectSignsRequestAndReadsMetadata(t *testing.T) {
	var gotMethod, gotPath, gotAuthorization string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod, gotPath, gotAuthorization = r.Method, r.URL.Path, r.Header.Get("Authorization")
		if r.Header.Get("x-amz-date") == "" || r.Header.Get("x-amz-content-sha256") != "UNSIGNED-PAYLOAD" {
			t.Errorf("missing SigV4 headers")
		}
		w.Header().Set("Content-Length", "1234")
		w.Header().Set("Content-Type", "video/mp4")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	endpoint, err := url.Parse(server.URL)
	if err != nil { t.Fatal(err) }
	s := &S3Presigner{endpoint:endpoint, bucket:"bucket", region:"eu-central-1", accessKey:"test-access", secretKey:strings.Repeat("s", 40)}
	meta, err := s.HeadObject(context.Background(), "movies/1/video/2026/11/a.mp4")
	if err != nil { t.Fatal(err) }
	if gotMethod != http.MethodHead || gotPath != "/bucket/movies/1/video/2026/11/a.mp4" { t.Fatalf("unexpected request %s %s", gotMethod, gotPath) }
	if !strings.HasPrefix(gotAuthorization, "AWS4-HMAC-SHA256 Credential=test-access/") { t.Fatalf("missing signed authorization: %q", gotAuthorization) }
	if meta.SizeBytes != 1234 || meta.ContentType != "video/mp4" { t.Fatalf("unexpected metadata: %+v", meta) }
}
