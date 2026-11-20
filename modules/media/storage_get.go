package media

import (
	"errors"
	"strings"
	"time"
)

// PresignGet creates a short-lived read URL for a private object.
func (s *S3Presigner) PresignGet(objectKey string, now time.Time, ttl time.Duration) (string, error) {
	if s == nil || s.endpoint == nil { return "", errors.New("S3 storage is not configured") }
	if !validObjectKey(objectKey) { return "", errors.New("invalid object key") }
	if ttl < time.Minute || ttl > 60*time.Minute { return "", errors.New("presigned download expiry must be between 1 and 60 minutes") }
	u := *s.endpoint
	u.Path = strings.TrimRight(u.Path, "/") + "/" + s.bucket + "/" + objectKey
	u.RawPath = ""
	now = now.UTC()
	date, amzDate := now.Format("20060102"), now.Format("20060102T150405Z")
	scope := date + "/" + s.region + "/s3/aws4_request"
	signedHeaders := "host"
	q := u.Query()
	q.Set("X-Amz-Algorithm", "AWS4-HMAC-SHA256")
	q.Set("X-Amz-Credential", s.accessKey+"/"+scope)
	q.Set("X-Amz-Date", amzDate)
	q.Set("X-Amz-Expires", strconvI(int(ttl.Seconds())))
	q.Set("X-Amz-SignedHeaders", signedHeaders)
	u.RawQuery = q.Encode()
	canonicalHeaders := "host:" + strings.ToLower(u.Host) + "\n"
	canonicalRequest := "GET\n" + u.EscapedPath() + "\n" + u.RawQuery + "\n" + canonicalHeaders + "\n" + signedHeaders + "\nUNSIGNED-PAYLOAD"
	stringToSign := "AWS4-HMAC-SHA256\n" + amzDate + "\n" + scope + "\n" + sha256Hex(canonicalRequest)
	q = u.Query()
	q.Set("X-Amz-Signature", hex.EncodeToString(hmacSHA256(s.signingKey(date), stringToSign)))
	u.RawQuery = q.Encode()
	return u.String(), nil
}

