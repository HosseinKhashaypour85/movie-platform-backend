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
