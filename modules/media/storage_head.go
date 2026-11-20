package media

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// ObjectMetadata is returned only after the storage origin confirms the object.
type ObjectMetadata struct {
	SizeBytes uint64
	ContentType string
}

func (s *S3Presigner) HeadObject(ctx context.Context, objectKey string) (ObjectMetadata, error) {
	if s == nil || s.endpoint == nil { return ObjectMetadata{}, errors.New("S3 storage is not configured") }
	if !validObjectKey(objectKey) { return ObjectMetadata{}, errors.New("invalid object key") }

	u := *s.endpoint
	u.Path = strings.TrimRight(u.Path, "/") + "/" + s.bucket + "/" + objectKey
	u.RawPath = ""
	now := time.Now().UTC()
	date, amzDate := now.Format("20060102"), now.Format("20060102T150405Z")
	scope := date + "/" + s.region + "/s3/aws4_request"
	payloadHash := "UNSIGNED-PAYLOAD"
	signedHeaders := "host;x-amz-content-sha256;x-amz-date"
	canonicalHeaders := "host:" + strings.ToLower(u.Host) + "\n" +
		"x-amz-content-sha256:" + payloadHash + "\n" +
		"x-amz-date:" + amzDate + "\n"
	canonicalRequest := "HEAD\n" + u.EscapedPath() + "\n\n" + canonicalHeaders + "\n" + signedHeaders + "\n" + payloadHash
	stringToSign := "AWS4-HMAC-SHA256\n" + amzDate + "\n" + scope + "\n" + sha256Hex(canonicalRequest)
	signature := hex.EncodeToString(hmacSHA256(s.signingKey(date), stringToSign))
	authorization := "AWS4-HMAC-SHA256 Credential=" + s.accessKey + "/" + scope +
		", SignedHeaders=" + signedHeaders + ", Signature=" + signature

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, u.String(), nil)
	if err != nil { return ObjectMetadata{}, err }
	req.Header.Set("x-amz-date", amzDate)
	req.Header.Set("x-amz-content-sha256", payloadHash)
	req.Header.Set("Authorization", authorization)

	client := &http.Client{Timeout: 12 * time.Second}
	resp, err := client.Do(req)
	if err != nil { return ObjectMetadata{}, fmt.Errorf("object storage HEAD failed: %w", err) }
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 4096))
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound { return ObjectMetadata{}, errors.New("uploaded object not found") }
		return ObjectMetadata{}, fmt.Errorf("object storage returned status %d", resp.StatusCode)
	}
	rawSize := strings.TrimSpace(resp.Header.Get("Content-Length"))
	size, err := strconv.ParseUint(rawSize, 10, 64)
	if err != nil { return ObjectMetadata{}, errors.New("object storage returned invalid content length") }
	return ObjectMetadata{SizeBytes:size, ContentType:strings.ToLower(strings.TrimSpace(strings.Split(resp.Header.Get("Content-Type"), ";")[0]))}, nil
}

