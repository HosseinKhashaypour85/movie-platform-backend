package media

import (
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestS3PresignPut(t *testing.T) {
	endpoint, _ := url.Parse("https://objects.example.test")
	s := &S3Presigner{endpoint: endpoint, bucket: "media-bucket", region: "eu-central-1", accessKey: "AKIATEST", secretKey: strings.Repeat("s", 40)}
	now := time.Date(2026, 11, 20, 12, 30, 0, 0, time.UTC)
	got, err := s.PresignPut("movies/42/video/2026/11/opaque.mp4", "video/mp4", now, 10*time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	u, err := url.Parse(got)
	if err != nil {
		t.Fatal(err)
	}
	if u.Scheme != "https" || u.Host != "objects.example.test" || u.Path != "/media-bucket/movies/42/video/2026/11/opaque.mp4" {
		t.Fatalf("unexpected upload URL: %s", got)
	}
	q := u.Query()
	if q.Get("X-Amz-Algorithm") != "AWS4-HMAC-SHA256" {
		t.Fatal("missing SigV4 algorithm")
	}
	if q.Get("X-Amz-Credential") != "AKIATEST/20261120/eu-central-1/s3/aws4_request" {
		t.Fatalf("unexpected credential scope: %s", q.Get("X-Amz-Credential"))
	}
	if q.Get("X-Amz-Date") != "20261120T123000Z" || q.Get("X-Amz-Expires") != "600" {
		t.Fatal("unexpected presign date or expiry")
	}
	if q.Get("X-Amz-SignedHeaders") != "content-type;host" {
		t.Fatal("content-type must be signed")
	}
	if len(q.Get("X-Amz-Signature")) != 64 {
		t.Fatal("invalid signature length")
	}
	if strings.Contains(got, "secret") {
		t.Fatal("secret key leaked in URL")
	}
}

func TestS3PresignRejectsInvalidConfigurationAndInputs(t *testing.T) {
	s := &S3Presigner{}
	if _, err := s.PresignPut("movies/1/a.mp4", "video/mp4", time.Now(), 10*time.Minute); err == nil {
		t.Fatal("expected unconfigured storage error")
	}
	endpoint, _ := url.Parse("https://objects.example.test")
	s = &S3Presigner{endpoint: endpoint, bucket: "bucket", region: "eu-central-1", accessKey: "key", secretKey: strings.Repeat("s", 40)}
	for _, key := range []string{"../escape.mp4", "/absolute.mp4", "movies//bad.mp4"} {
		if _, err := s.PresignPut(key, "video/mp4", time.Now(), 10*time.Minute); err == nil {
			t.Errorf("accepted invalid key %q", key)
		}
	}
	if _, err := s.PresignPut("movies/1/a.mp4", "video/mp4", time.Now(), 30*time.Minute); err == nil {
		t.Fatal("accepted excessive URL lifetime")
	}
}

func TestS3PresignReadMethods(t *testing.T) {
	endpoint, _ := url.Parse("https://objects.example.test")
	s := &S3Presigner{endpoint: endpoint, bucket: "media-bucket", region: "eu-central-1", accessKey: "AKIATEST", secretKey: strings.Repeat("s", 40)}
	now := time.Date(2026, 11, 20, 12, 30, 0, 0, time.UTC)

	getURL, err := s.PresignGet("hls/42/master.m3u8", now, time.Minute)
	if err != nil { t.Fatal(err) }
	headURL, err := s.PresignHead("hls/42/master.m3u8", now, time.Minute)
	if err != nil { t.Fatal(err) }
	get, err := url.Parse(getURL)
	if err != nil { t.Fatal(err) }
	head, err := url.Parse(headURL)
	if err != nil { t.Fatal(err) }
	for _, parsed := range []*url.URL{get, head} {
		if parsed.Path != "/media-bucket/hls/42/master.m3u8" { t.Fatalf("unexpected object path: %s", parsed.Path) }
		if parsed.Query().Get("X-Amz-SignedHeaders") != "host" { t.Fatal("read URL must sign host") }
		if len(parsed.Query().Get("X-Amz-Signature")) != 64 { t.Fatal("missing SigV4 signature") }
	}
	if get.Query().Get("X-Amz-Signature") == head.Query().Get("X-Amz-Signature") {
		t.Fatal("GET and HEAD must have method-specific signatures")
	}
}
