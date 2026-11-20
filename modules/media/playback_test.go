package media

import (
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestSignedPlaybackURL(t *testing.T) {
	key := []byte(strings.Repeat("k", 40))
	base, err := url.Parse("https://cdn.example.com")
	if err != nil { t.Fatal(err) }
	p := &Playback{base: base, key: key, ttl: 5 * time.Minute}
	now := time.Unix(1800000000, 0)
	got, expires, err := p.URL("movies/42/master.m3u8", now)
	if err != nil { t.Fatal(err) }
	u, err := url.Parse(got)
	if err != nil { t.Fatal(err) }
	if u.Scheme != "https" || u.Host != "cdn.example.com" || u.Path != "/movies/42/master.m3u8" {
		t.Fatalf("unexpected URL: %s", got)
	}
	exp, err := strconv.ParseInt(u.Query().Get("expires"), 10, 64)
	if err != nil || exp != expires.Unix() { t.Fatalf("invalid expiry: %v", err) }
	sig := u.Query().Get("signature")
	if !VerifyPlaybackSignature(key, u.Path, exp, sig, now) {
		t.Fatal("valid signature rejected")
	}
	if VerifyPlaybackSignature(key, "/movies/43/master.m3u8", exp, sig, now) {
		t.Fatal("signature must bind the object path")
	}
	if VerifyPlaybackSignature(key, u.Path, exp, sig, expires) {
		t.Fatal("expired signature accepted")
	}
	if VerifyPlaybackSignature([]byte(strings.Repeat("x", 40)), u.Path, exp, sig, now) {
		t.Fatal("wrong key accepted")
	}
	if VerifyPlaybackSignature(key, u.Path, exp, sig+"A", now) {
		t.Fatal("tampered signature accepted")
	}
	if VerifyPlaybackSignature(key, u.Path, now.Add(time.Hour).Unix(), signature(key, u.Path, now.Add(time.Hour).Unix()), now) {
		t.Fatal("excessively long expiry accepted")
	}
}

func TestPlaybackObjectKeyValidation(t *testing.T) {
	for _, key := range []string{"", "/root/master.m3u8", "../master.m3u8", "a/../master.m3u8", "a//master.m3u8", "a/./master.m3u8", "a\\master.m3u8"} {
		if validObjectKey(key) { t.Errorf("accepted unsafe key %q", key) }
	}
	for _, key := range []string{"a/master.m3u8", "a/b/c/master.m3u8"} {
		if !validObjectKey(key) { t.Errorf("rejected valid key %q", key) }
	}
}

func TestPlaybackRequiresManifest(t *testing.T) {
	base, _ := url.Parse("https://cdn.example.com")
	p := &Playback{base: base, key: []byte(strings.Repeat("k", 40)), ttl: 5*time.Minute}
	for _, key := range []string{"video.mp4", "movies/a/master.m3u8/../x.m3u8"} {
		if _, _, err := p.URL(key, time.Now()); err == nil {
			t.Fatalf("accepted invalid manifest key %q", key)
		}
	}
}
