package media

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
)

type S3Presigner struct {
	endpoint *url.URL
	bucket string
	region string
	accessKey string
	secretKey string
}

func NewS3PresignerFromEnv() (*S3Presigner, error) {
	endpointRaw := strings.TrimSpace(os.Getenv("S3_ENDPOINT"))
	bucket := strings.TrimSpace(os.Getenv("S3_BUCKET"))
	region := strings.TrimSpace(os.Getenv("S3_REGION"))
	access := strings.TrimSpace(os.Getenv("S3_ACCESS_KEY_ID"))
	secret := os.Getenv("S3_SECRET_ACCESS_KEY")
	if endpointRaw == "" && bucket == "" && region == "" && access == "" && secret == "" {
		return nil, nil
	}
	if endpointRaw == "" || bucket == "" || region == "" || access == "" || len(secret) < 16 {
		return nil, errors.New("S3 requires endpoint, bucket, region, access key and secret key")
	}
	u, err := url.Parse(endpointRaw)
	if err != nil || u.Scheme != "https" || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return nil, errors.New("S3_ENDPOINT must be an HTTPS endpoint without credentials, query or fragment")
	}
	if strings.Contains(bucket, "/") || strings.Contains(bucket, "..") || strings.ContainsAny(bucket, " ") {
		return nil, errors.New("invalid S3 bucket")
	}
	return &S3Presigner{endpoint:u,bucket:bucket,region:region,accessKey:access,secretKey:secret},nil
}

func hmacSHA256(key []byte, value string) []byte {
	m := hmac.New(sha256.New,key)
	m.Write([]byte(value))
	return m.Sum(nil)
}
func sha256Hex(value string) string {
	sum:=sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}
func (s *S3Presigner) signingKey(date string) []byte {
	kDate:=hmacSHA256([]byte("AWS4"+s.secretKey),date)
	kRegion:=hmacSHA256(kDate,s.region)
	kService:=hmacSHA256(kRegion,"s3")
	return hmacSHA256(kService,"aws4_request")
}

func (s *S3Presigner) PresignPut(objectKey, contentType string, now time.Time, ttl time.Duration) (string,error) {
	if s==nil || s.endpoint==nil { return "",errors.New("S3 storage is not configured") }
	if !validObjectKey(objectKey) || contentType=="" || strings.ContainsAny(contentType,"\r\n") { return "",errors.New("invalid upload object key or content type") }
	if ttl<time.Minute || ttl>15*time.Minute { return "",errors.New("presigned upload expiry must be between 1 and 15 minutes") }
	u:=*s.endpoint
	prefix:=strings.TrimRight(u.Path,"/")
	u.Path=prefix+"/"+s.bucket+"/"+objectKey
	u.RawPath=""
	now=now.UTC()
	date:=now.Format("20060102")
	amzDate:=now.Format("20060102T150405Z")
	scope:=date+"/"+s.region+"/s3/aws4_request"
	signedHeaders:="content-type;host"
	q:=u.Query()
	q.Set("X-Amz-Algorithm","AWS4-HMAC-SHA256")
	q.Set("X-Amz-Credential",s.accessKey+"/"+scope)
	q.Set("X-Amz-Date",amzDate)
	q.Set("X-Amz-Expires",strconvI(int(ttl.Seconds())))
	q.Set("X-Amz-SignedHeaders",signedHeaders)
	u.RawQuery=q.Encode()
	canonicalHeaders:="content-type:"+strings.TrimSpace(contentType)+"\n"+"host:"+u.Host+"\n"
	canonicalRequest:="PUT\n"+u.EscapedPath()+"\n"+u.RawQuery+"\n"+canonicalHeaders+"\n"+signedHeaders+"\nUNSIGNED-PAYLOAD"
	stringToSign:="AWS4-HMAC-SHA256\n"+amzDate+"\n"+scope+"\n"+sha256Hex(canonicalRequest)
	signature:=hex.EncodeToString(hmacSHA256(s.signingKey(date),stringToSign))
	q=u.Query()
	q.Set("X-Amz-Signature",signature)
	u.RawQuery=q.Encode()
	return u.String(),nil
}

func strconvI(v int) string {
	if v<0 { return "0" }
	if v==0 { return "0" }
	var b [20]byte
	i:=len(b)
	for v>0 { i--; b[i]=byte('0'+v%10); v/=10 }
	return string(b[i:])
}

// canonicalQuery is retained as a small helper for deterministic signing tests.
func canonicalQuery(values url.Values) string {
	keys:=make([]string,0,len(values))
	for k:=range values { keys=append(keys,k) }
	sort.Strings(keys)
	return values.Encode()
}
