# Movie Platform Backend

A Go/MySQL API foundation for the movie and series platform.

## Included
- Strict environment validation; secrets are never provided by source code.
- MySQL connection pool, ping and bounded HTTP timeouts.
- Request IDs, panic recovery, security headers, JSON errors, request size limit.
- CORS allow-list and in-process IP rate limiting (edge rate limiting is still required).
- Short-lived HS256 JWT access tokens with issuer/expiry validation.
- User registration/login with bcrypt, input validation and generic login errors.
- Authenticated profile and per-user watch-history CRUD/upsert.
- CI workflow for gofmt, vet and tests.

## Critical credential notice
The previous main branch committed a `.env` file containing a JWT secret. Treat it as exposed. Rotate the JWT secret before any deployment and invalidate tokens signed with the old value. Removing the file in this branch does not erase it from Git history; purge history if it contained a real credential.

## Setup
Requires Go 1.27.3 and MySQL 8+.

```bash
cp .env.example .env
# Generate a unique JWT secret:
openssl rand -base64 48
go mod tidy
go run ./cmd/api
```

Create database and dedicated least-privilege account. For production, use TLS to MySQL and a secret manager. This starter uses AutoMigrate only outside production; production schema changes must run as a reviewed migration job.

## Endpoints
- GET /health/live
- GET /health/ready
- POST /api/v1/auth/register
- POST /api/v1/auth/login (returns short-lived access token and rotating refresh token)\n- POST /api/v1/auth/refresh\n- POST /api/v1/auth/logout (Bearer token + refresh token)\n- POST /api/v1/auth/logout-all (Bearer token)
- GET /api/v1/profile (Bearer token)
- GET /api/v1/watch-history (Bearer token)
- PUT /api/v1/watch-history/{movieID} (Bearer token)
- GET /api/v1/watch-history/{movieID} (Bearer token)
- DELETE /api/v1/watch-history/{movieID} (Bearer token)

## Not yet a complete production product
Catalog CRUD, seasons/episodes, content rights metadata, signed playback URLs, HLS packaging, search, genres, actors, ratings/comments, ads, analytics, admin RBAC, audit trail delivery, email/OTP verification, refresh-token rotation/revocation, password reset, AI job queue/workers, billing and cloud storage/CDN integrations remain to be specified and implemented. Security requires deployment review, TLS/reverse-proxy configuration, external rate limiting, secrets rotation, database migrations, monitoring, backups and penetration testing.

## Production database migrations

The repository contains ordered, versioned SQL migrations under `migrations/`. Apply pending migrations with:

```bash
go run ./cmd/migrate
```

The runner records applied filenames in `schema_migrations` and applies only unapplied `*.up.sql` files in lexical order. Back up production before applying schema changes. MySQL DDL may implicitly commit, so migrations are written to be rerunnable with `IF NOT EXISTS` where applicable. Review rollback SQL manually; the application does not automatically run down migrations.

Do not use AutoMigrate in production. Configure database credentials and the required application environment variables through the deployment secret manager.

## Signed HLS playback (integration contract)

Authenticated clients can call `GET /api/v1/playback/movies/{id}` or
`GET /api/v1/playback/episodes/{id}` with a bearer access token. Only
published titles/episodes with an active `video` media asset can receive
a manifest URL. Assets must use an HLS `.m3u8` storage key.

Configure both `PLAYBACK_CDN_BASE_URL` (an HTTPS origin) and
`PLAYBACK_SIGNING_KEY` (a dedicated random secret of at least 32 bytes).
When neither is set, the endpoint responds 503; a partially configured
signer prevents API startup. The response includes `manifest_url`,
`expires_at`, `asset_id`, and `format: hls`. The URL expires after
five minutes and must not be cached or logged.

**Important:** This backend only issues signed manifest URLs; it does not
serve HLS playlists/segments or configure an edge CDN. The signature is
base64url(HMAC-SHA256(key, canonicalPath + "\\n" + decimalUnixExpiry)).
The query parameters are `expires` and `signature`. The provided
`VerifyPlaybackSignature` helper is a reference verifier for a trusted
edge integration. An actual deployment MUST enforce signatures and
expiry on every protected media request, including HLS variant playlists,
segments and encryption keys, with a suitable segment authorization
mechanism (e.g. edge-validated signed cookies or manifest rewriting).
Signing only the master playlist does **not** secure HLS segments.
Keep the storage origin private, reject noncanonical paths, and do not
enable public access to the bucket. CDN provider integration and end-to-end
playback tests are still outstanding.

## Direct media uploads to S3-compatible storage

Configure all five `S3_*` variables in `.env` to enable direct uploads.
The endpoint must be HTTPS and the bucket should remain private. The API
supports AWS S3 and S3-compatible endpoints using path-style bucket URLs.

An authenticated administrator can request a short-lived upload URL:

```http
POST /api/v1/admin/media/uploads
Authorization: Bearer <access-token>
Content-Type: application/json

{
  "target_type": "movie",
  "target_id": 42,
  "kind": "video",
  "content_type": "video/mp4"
}
```

The response contains a `PUT` URL, the exact `Content-Type` header to send,
a generated `storage_key`, and a 10-minute expiry. Upload the bytes directly
to object storage using that method and header. Supported content types are
MP4 video and JPEG/PNG/WebP images. The upload URL is signed with AWS Signature
Version 4 and binds the HTTP method, object path, expiry, host and content type.

After the PUT succeeds, finalize the upload with `POST /api/v1/admin/media/uploads/complete` using the same `target_type`, `target_id`, `kind`, and returned `storage_key`. The API signs a HEAD request to the private bucket, checks the stored content type and size (2 GiB maximum for MP4, 20 MiB for images), and creates an inactive media asset. Repeated finalization is rejected. The asset remains inactive until a trusted processing worker has generated and verified HLS output; the current endpoint does not scan files or transcode source video. Do not expose uploaded objects publicly. Configure bucket CORS to allow PUT from the admin web origin and
restrict IAM credentials to the intended bucket/prefix and required actions.

## FFmpeg HLS transcoding worker

A verified MP4 upload automatically creates a queued `transcode_jobs` row.
Run the separate worker process on a host/container with FFmpeg installed and
the same database and private S3-compatible storage credentials:

```bash
go run ./cmd/transcoder
```

The worker claims jobs using MySQL row locks, downloads the private source
through a short-lived signed GET URL, packages HLS VOD (6-second MPEG-TS
segments and an `index.m3u8` manifest) with FFmpeg/libx264/AAC, and uploads
the generated playlist and segments to a unique private `hls/{asset-id}/...`
prefix. Only after every upload succeeds does it create an active HLS asset
and mark the job completed. Failed jobs retry up to three attempts. Temporary
files are removed after processing; configure `TRANSCODE_TEMP_DIR` on a
volume with sufficient free disk (at least the source size plus output space).

Install FFmpeg with libx264 support and `ffprobe` in the worker image. The worker is a separate process from the API; do not run it in the public HTTP request path. It probes source dimensions and creates only renditions that do not upscale the source: 360p, 480p, 720p and 1080p. The HLS master playlist references each available variant, allowing compatible players to switch quality based on measured bandwidth and buffer health. Variant playlists and segments are uploaded under the same private, unique asset prefix; the active media asset points to `master.m3u8`. The CDN edge must still enforce signed playback authorization for the master manifest, every variant playlist/segment, and any keys (for example through signed cookies or manifest rewriting). This worker does not yet create thumbnails, scan for malware, or automatically delete source files.
