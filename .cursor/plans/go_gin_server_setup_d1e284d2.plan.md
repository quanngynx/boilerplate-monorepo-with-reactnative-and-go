---
name: Go Gin Server Setup
overview: Scaffold a Go Gin API in apps/server with GORM (Postgres), Redis, OWASP-aligned error handling and logging, a hardened Dockerfile, and docker-compose for local dev (postgres:18.4-bookworm, redis:8.8.0-alpine).
todos:
  - id: scaffold-go
    content: Create apps/server Go module, folder layout, go.mod with gin, gorm, redis, godotenv
    status: completed
  - id: gin-core
    content: Implement config, router, CORS middleware, /health and /api/v1/hello handlers
    status: completed
  - id: gorm-redis
    content: Add GORM Postgres connection, Redis client, AutoMigrate starter model, health checks for DB/Redis
    status: completed
  - id: error-logging
    content: Implement RFC 7807 error responses, typed apperror package, slog structured logging with OWASP vocabulary helpers
    status: completed
  - id: docker
    content: Add OWASP-hardened Dockerfile and docker-compose.yml (postgres 18.4-bookworm, redis 8.8.0-alpine, server)
    status: completed
  - id: turbo-shim
    content: Add apps/server/package.json, root dev:server script, turbo.json bin/** output, gitignore
    status: completed
  - id: wire-web
    content: Add Vite /api proxy in apps/web/vite.config.ts and VITE_SERVER_URL in packages/env/src/web.ts
    status: completed
  - id: wire-native
    content: Document apps/native/.env with EXPO_PUBLIC_SERVER_URL=http://localhost:8080
    status: completed
  - id: verify-docs
    content: Smoke-test endpoints, docker-compose stack, and update README with server + Docker setup
    status: completed
isProject: false
---

# Go Gin Server Setup in `apps/server`

## Context

- [`apps/server`](apps/server) does not exist yet — no Go code in the repo.
- Clients already anticipate a backend:
  - Native: [`packages/env/src/native.ts`](packages/env/src/native.ts) requires `EXPO_PUBLIC_SERVER_URL`
  - Web deploy: [`packages/infra/alchemy.run.ts`](packages/infra/alchemy.run.ts) passes `VITE_SERVER_URL`
- Web dev runs on port **3001** ([`apps/web/vite.config.ts`](apps/web/vite.config.ts)).
- Turbo orchestrates apps via `package.json` scripts ([`turbo.json`](turbo.json), root [`package.json`](package.json)).

## Architecture

```mermaid
flowchart TB
  subgraph clients [Clients]
    Web["apps/web :3001"]
    Native["apps/native Expo"]
  end

  subgraph compose [docker-compose dev stack]
  subgraph api [Go API]
    Gin["server :8080"]
  end
  subgraph data [Data layer]
    PG["postgres:18.4-bookworm"]
    Redis["redis:8.8.0-alpine"]
  end
  end

  Web -->|"Vite proxy /api"| Gin
  Native -->|"EXPO_PUBLIC_SERVER_URL"| Gin
  Gin -->|"GORM"| PG
  Gin -->|"go-redis"| Redis
```

**Dev networking strategy:**

- **Web**: Vite proxies `/api` → `http://localhost:8080` (same-origin, no CORS in dev).
- **Native**: Calls `EXPO_PUBLIC_SERVER_URL` directly; Gin CORS middleware allows Expo origins.
- **Docker**: `docker compose up` runs Postgres, Redis, and the API together on an internal bridge network.

---

## Phase 1 — Scaffold Go project

### 1.1 Directory layout

```
apps/server/
├── cmd/server/main.go
├── Dockerfile
├── docker-compose.yml
├── .dockerignore
├── internal/
│   ├── config/config.go
│   ├── database/database.go       # GORM init + AutoMigrate
│   ├── cache/redis.go             # go-redis client
│   ├── model/                     # GORM models (starter: no tables or a minimal health metadata model)
│   ├── repository/                # thin repo layer over GORM
│   ├── handler/
│   │   ├── health.go
│   │   └── hello.go
│   ├── middleware/
│   │   ├── cors.go
│   │   ├── request_logger.go      # slog per-request logging
│   │   └── error_handler.go       # global panic/error → RFC 7807
│   ├── apperror/errors.go         # typed domain errors + HTTP status mapping
│   ├── response/problem.go        # RFC 7807 Problem Details builder
│   ├── logger/logger.go           # slog setup + OWASP vocabulary event helpers
│   └── router/router.go
├── go.mod
├── go.sum
├── .env.example
├── .gitignore
└── package.json
```

### 1.2 Dependencies

```bash
cd apps/server
go mod init github.com/GenAI-Fund/boilerplate-monorepo-with-reactnative-and-go/apps/server

# HTTP
go get github.com/gin-gonic/gin@v1.10.0
go get github.com/joho/godotenv@v1.5.1

# ORM + Postgres
go get gorm.io/gorm@v1.30.0
go get gorm.io/driver/postgres@v1.6.0

# Redis
go get github.com/redis/go-redis/v9@v9.7.0
```

Use `log/slog` (stdlib) for structured logging — no third-party logger required.

### 1.3 Config env vars

| Variable | Default | Notes |
| --- | --- | --- |
| `PORT` | `8080` | API bind port |
| `GIN_MODE` | `debug` | `release` in production |
| `ALLOWED_ORIGINS` | `http://localhost:3001` | Comma-separated CORS |
| `DATABASE_URL` | — | `postgres://user:pass@postgres:5432/app?sslmode=disable` in compose |
| `REDIS_URL` | — | `redis://redis:6379/0` in compose |
| `LOG_LEVEL` | `info` | `debug` in dev, `info` in prod |
| `APP_NAME` | `server` | Logged on `sys_startup` |

Document all vars in [`apps/server/.env.example`](apps/server/.env.example).

### 1.4 Starter routes

| Method | Path            | Purpose                                    |
| ------ | --------------- | ------------------------------------------ |
| `GET`  | `/health`       | Liveness; includes `db` and `redis` status |
| `GET`  | `/api/v1/hello` | Smoke-test JSON endpoint                   |

`/health` returns `503` if Postgres or Redis is unreachable (fail fast for orchestrators).

---

## Phase 2 — GORM + Redis

### 2.1 GORM (`internal/database/database.go`)

- Open Postgres via `gorm.io/driver/postgres` using `DATABASE_URL`.
- Configure connection pool: `SetMaxOpenConns`, `SetMaxIdleConns`, `SetConnMaxLifetime`.
- Run `AutoMigrate` for starter models in `internal/model/` (e.g. a minimal placeholder or skip until first real entity).
- Expose `*gorm.DB` via dependency injection into handlers/repos.

### 2.2 Repository pattern (`internal/repository/`)

- Thin wrapper over GORM for testability.
- Handlers depend on repository interfaces, not raw `*gorm.DB`.

### 2.3 Redis (`internal/cache/redis.go`)

- `go-redis` client from `REDIS_URL`.
- Used in `/health` ping and available for session/cache later.
- Graceful shutdown: close Redis client on SIGTERM.

### 2.4 `main.go` lifecycle

1. Load config
2. Init slog logger → stdout JSON
3. Connect Postgres (GORM) — fail startup if unreachable
4. Connect Redis — fail startup if unreachable
5. Log `sys_startup` (OWASP vocabulary)
6. Build router, start HTTP server
7. On shutdown: log `sys_shutdown`, close DB pool and Redis

---

## Phase 3 — Error handling (OWASP)

Ref: [Error Handling Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Error_Handling_Cheat_Sheet.html)

### 3.1 Principles

- **Never leak stack traces, SQL errors, or internal paths to clients** in production.
- **Log full error details server-side**; return generic messages to clients.
- Use **RFC 7807 Problem Details** (`application/problem+json`) for API errors.
- **4xx** for client mistakes; **5xx** only for unexpected server failures.
- In `GIN_MODE=debug`, optionally include more detail in responses (dev only).

### 3.2 `internal/apperror/errors.go`

Typed errors with HTTP status mapping:

```go
type AppError struct {
    Code       string // stable machine-readable code, e.g. "VALIDATION_FAILED"
    Message    string // safe client message
    HTTPStatus int
    Err        error  // wrapped cause — logged, never sent to client
}
```

Predefined errors: `ErrNotFound`, `ErrValidation`, `ErrUnauthorized`, `ErrInternal`.

### 3.3 `internal/response/problem.go`

RFC 7807 response builder:

```json
{
  "type": "about:blank",
  "title": "Bad Request",
  "status": 400,
  "detail": "Invalid request parameters",
  "instance": "/api/v1/hello"
}
```

Set header `Content-Type: application/problem+json`. Optionally set `X-ERROR: true` per OWASP examples.

### 3.4 `internal/middleware/error_handler.go`

- Gin middleware that catches panics and `*apperror.AppError`.
- Logs internal cause at `Error` level with request ID.
- Returns Problem Details JSON — no `err.Error()` from wrapped DB/driver errors in production.
- Replace default `gin.Recovery()` with custom recovery that follows the same pattern.

### 3.5 Handler pattern

```go
if err := c.ShouldBindJSON(&req); err != nil {
    // log input_validation_fail event
    _ = c.Error(apperror.Validation("Invalid request body", err))
    return
}
```

---

## Phase 4 — Logging (OWASP)

Refs:

- [Logging Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Logging_Cheat_Sheet.html)
- [Logging Vocabulary Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Logging_Vocabulary_Cheat_Sheet.html)

### 4.1 `internal/logger/logger.go`

- Initialize `slog` with **JSON handler to stdout** (container-friendly; compose captures logs).
- Standard attributes on every log: `timestamp`, `level`, `app`, `request_id`, `event_type`.
- `LOG_LEVEL` env controls verbosity; default `info` in production.

### 4.2 OWASP vocabulary event helpers

Helper functions that emit consistently named security events:

| Event                   | When                   |
| ----------------------- | ---------------------- |
| `sys_startup`           | Server starts          |
| `sys_shutdown`          | Graceful shutdown      |
| `sys_crash`             | Unhandled panic        |
| `input_validation_fail` | Bind/validation errors |
| `authn_login_fail`      | (stub for future auth) |
| `authz_fail`            | (stub for future auth) |

Format: `event_type` field matches OWASP vocabulary (e.g. `input_validation_fail`).

### 4.3 Request logging middleware

Log per request (when/where/who/what):

- `request_id` (UUID, propagated via `X-Request-ID` header)
- `method`, `path`, `status`, `duration_ms`
- `client_ip` (from `X-Forwarded-For` or `RemoteAddr`)
- `user_agent`

**Never log:** passwords, tokens, `DATABASE_URL`, `REDIS_URL`, session IDs, full request bodies with PII.

### 4.4 Sanitization

- Strip CR/LF from user-supplied values before logging (log injection prevention).
- Hash or omit sensitive fields.

---

## Phase 5 — Dockerfile (OWASP Docker Security)

Ref: [Docker Security Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Docker_Security_Cheat_Sheet.html)

### 5.1 Multi-stage build

```dockerfile
# Stage 1: build
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server ./cmd/server

# Stage 2: runtime
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /server /server
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/server"]
```

### 5.2 OWASP rules applied

| Rule | Implementation |
| --- | --- |
| RULE #2 Set a user | `USER nonroot:nonroot` (distroless) |
| RULE #4 No privilege escalation | `security_opt: no-new-privileges:true` in compose |
| RULE #5 Inter-container connectivity | Dedicated `server_internal` bridge network |
| RULE #7 Limit resources | `deploy.resources.limits` in compose (memory/CPU) |
| RULE #8 Read-only filesystem | `read_only: true` + `tmpfs: /tmp` in compose |
| RULE #9 Supply chain | Pin base image digests; `go mod verify` in CI |
| RULE #13 Supply chain | Multi-stage build; minimal runtime image |
| RULE #5a Port binding | Bind API to `127.0.0.1:8080:8080` when not using compose networking |

### 5.3 `.dockerignore`

```
bin/
.git/
.env
*.md
tmp/
```

---

## Phase 6 — docker-compose.yml

Location: [`apps/server/docker-compose.yml`](apps/server/docker-compose.yml)

```yaml
services:
  postgres:
    image: postgres:18.4-bookworm
    environment:
      POSTGRES_USER: app
      POSTGRES_PASSWORD: app
      POSTGRES_DB: app
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "127.0.0.1:5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U app -d app"]
      interval: 5s
      timeout: 5s
      retries: 5
    networks: [server_internal]

  redis:
    image: redis:8.8.0-alpine
    command: redis-server --appendonly yes
    volumes:
      - redis_data:/data
    ports:
      - "127.0.0.1:6379:6379"
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 3s
      retries: 5
    networks: [server_internal]

  server:
    build:
      context: .
      dockerfile: Dockerfile
    env_file: .env
    environment:
      DATABASE_URL: postgres://app:app@postgres:5432/app?sslmode=disable
      REDIS_URL: redis://redis:6379/0
      GIN_MODE: release
      LOG_LEVEL: info
    ports:
      - "127.0.0.1:8080:8080"
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    security_opt:
      - no-new-privileges:true
    read_only: true
    tmpfs:
      - /tmp
    deploy:
      resources:
        limits:
          cpus: "1.0"
          memory: 512M
    networks: [server_internal]

networks:
  server_internal:
    driver: bridge

volumes:
  postgres_data:
  redis_data:
```

**Dev scripts** in [`apps/server/package.json`](apps/server/package.json):

```json
{
  "scripts": {
    "dev": "go run ./cmd/server",
    "dev:docker": "docker compose up --build",
    "dev:deps": "docker compose up postgres redis -d",
    "build": "go build -o ./bin/server ./cmd/server",
    "docker:build": "docker compose build",
    "docker:down": "docker compose down -v"
  }
}
```

**Local dev without Docker for the API:** run `pnpm dev:deps` to start only Postgres + Redis, then `pnpm dev:server` with `DATABASE_URL` / `REDIS_URL` pointing at `localhost`.

---

## Phase 7 — Turborepo / pnpm integration

**[`apps/server/package.json`](apps/server/package.json)** — full scripts including docker targets.

**Root [`package.json`](package.json)** — add:

```json
"dev:server": "turbo -F server dev",
"dev:server:docker": "turbo -F server dev:docker"
```

**[`turbo.json`](turbo.json)** — extend build outputs: `["dist/**", "bin/**"]`.

**[`.gitignore`](.gitignore)** — add `apps/server/bin/`.

---

## Phase 8 — Wire clients

### Web — Vite proxy ([`apps/web/vite.config.ts`](apps/web/vite.config.ts))

```ts
proxy: {
  "/api": { target: "http://localhost:8080", changeOrigin: true },
},
```

### Web env ([`packages/env/src/web.ts`](packages/env/src/web.ts))

Add `VITE_SERVER_URL` with `clientPrefix: "VITE_"`.

### Native — [`apps/native/.env`](apps/native/.env)

```env
EXPO_PUBLIC_SERVER_URL=http://localhost:8080
```

---

## Phase 9 — Verify end-to-end

1. `docker compose up --build` → all three services healthy
2. `curl http://localhost:8080/health` → `200` with `"db":"ok","redis":"ok"`
3. `curl http://localhost:8080/api/v1/hello` → JSON
4. Trigger validation error → `application/problem+json`, no stack trace
5. Check container logs → structured JSON with `sys_startup`, request logs
6. `pnpm dev:web` → `fetch("/api/v1/hello")` via proxy
7. `pnpm dev:native` → hits API via `EXPO_PUBLIC_SERVER_URL`

---

## Success criteria

- Gin API runs locally and in Docker on port 8080
- GORM connects to Postgres 18.4; Redis 8.8.0 client works
- `/health` reports DB and Redis status
- Errors return RFC 7807 Problem Details; no internal details leak in production
- Logs are structured JSON to stdout with OWASP vocabulary event types
- Dockerfile runs as non-root with minimal distroless image
- docker-compose starts postgres:18.4-bookworm + redis:8.8.0-alpine + server
- Web and native clients can reach the API
- README documents env vars, Docker workflow, and local dev options

## Key files to create or modify

| Action | File |
| --- | --- |
| Create | `apps/server/cmd/server/main.go` |
| Create | `apps/server/internal/config/config.go` |
| Create | `apps/server/internal/database/database.go` |
| Create | `apps/server/internal/cache/redis.go` |
| Create | `apps/server/internal/apperror/errors.go` |
| Create | `apps/server/internal/response/problem.go` |
| Create | `apps/server/internal/logger/logger.go` |
| Create | `apps/server/internal/middleware/error_handler.go` |
| Create | `apps/server/internal/middleware/request_logger.go` |
| Create | `apps/server/internal/router/router.go` |
| Create | `apps/server/Dockerfile`, `docker-compose.yml`, `.dockerignore` |
| Create | `apps/server/go.mod`, `go.sum`, `package.json`, `.env.example` |
| Modify | [`package.json`](package.json) — `dev:server`, `dev:server:docker` |
| Modify | [`turbo.json`](turbo.json) — `bin/**` output |
| Modify | [`apps/web/vite.config.ts`](apps/web/vite.config.ts) — proxy |
| Modify | [`packages/env/src/web.ts`](packages/env/src/web.ts) — `VITE_SERVER_URL` |
| Modify | [`README.md`](README.md) — server, Docker, GORM, logging |

## Deferred (out of initial scope)

- Auth (Better Auth / JWT) — stub OWASP `authn_*` / `authz_*` events only
- Production deploy target (Fly.io / Railway / K8s)
- Container image scanning in CI (Trivy — recommended follow-up per OWASP RULE #9)
- Database migrations beyond GORM AutoMigrate (consider golang-migrate later)
