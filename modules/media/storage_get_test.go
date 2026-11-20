package media

import (
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestPresignGet(t *testing.T) {
	endpoint,err:=url.Parse("https://objects.example.test")
	if err!=nil { t.Fatal(err) }
	s:=&S3Presigner{endpoint:endpoint,bucket:"private",region:"eu-central-1",accessKey:"TESTKEY",secretKey:strings.Repeat("s",40)}
	now:=time.Date(2026,11,20,12,30,0,0,time.UTC)
	raw,err:=s.PresignGet("movies/42/video/source.mp4",now,30*time.Minute)
	if err!=nil { t.Fatal(err) }
	u,err:=url.Parse(raw)
	if err!=nil { t.Fatal(err) }
	q:=u.Query()
	if u.Path!="/private/movies/42/video/source.mp4" { t.Fatalf("unexpected path %s",u.Path) }
	if q.Get("X-Amz-Algorithm")!="AWS4-HMAC-SHA256" || q.Get("X-Amz-Date")!="20261120T123000Z" || q.Get("X-Amz-Expires")!="1800" { t.Fatalf("unexpected query: %v",q) }
	if q.Get("X-Amz-SignedHeaders")!="host" || len(q.Get("X-Amz-Signature"))!=64 { t.Fatalf("invalid signature params: %v",q) }
	if strings.Contains(raw,"secret") { t.Fatal("secret leaked into URL") }
	if _,err:=s.PresignGet("../escape",now,30*time.Minute); err==nil { t.Fatal("accepted invalid object key") }
	if _,err:=s.PresignGet("valid/key",now,2*time.Hour); err==nil { t.Fatal("accepted excessive expiry") }
}
