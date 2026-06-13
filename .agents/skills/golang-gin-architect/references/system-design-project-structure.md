# System Design — Go Package Layout at Scale

Package layout patterns for Go Gin APIs at different team/endpoint scales.

Companion to `system-design-dependency-graphs.md` (no-cycles rule, dependency inversion).

---

## Small: 1-3 devs, < 20 endpoints — Flat Layout

```
myapp/
├── cmd/api/main.go
├── internal/
│   ├── handler/          # All HTTP handlers
│   ├── service/          # All business logic
│   ├── repository/       # All data access
│   └── domain/           # Entities, interfaces, errors
├── pkg/
│   └── middleware/
└── go.mod
```

One package per layer. Works until the team grows or domains conflict.

---

## Medium: 3-8 devs, 20-100 endpoints — Feature Modules

```
myapp/
├── cmd/api/main.go
├── internal/
│   ├── user/
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   ├── model.go
│   │   └── ports.go        # interfaces this module needs from others
│   ├── order/
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   ├── model.go
│   │   └── ports.go
│   └── platform/           # shared infrastructure — not a business domain
│       ├── database/
│       ├── auth/
│       └── middleware/
├── pkg/
│   └── httputil/
└── go.mod
```

Each feature module owns its full vertical slice. `ports.go` defines external interface contracts.

**The `platform/` package** is infrastructure — not a domain. Business modules import from `platform/`, never the reverse.

---

## Large: 8+ devs, 100+ endpoints — Evaluate Extraction

Before extracting to microservices, all 3 must be true for the same module:

1. Independent scaling requirements (measured, not hypothetical)
2. Need to deploy independently (different release cadence, different team)
3. Clear, stable interface that won't change frequently

```
myapp/                         extracted/
├── cmd/api/main.go            ├── cmd/catalog-api/main.go
├── internal/                  ├── internal/
│   ├── user/                  │   └── catalog/
│   ├── order/                 └── go.mod
│   └── platform/
└── go.mod
```

**The extracted service talks to the monolith via HTTP or gRPC, not shared DB tables.**

---

## Shared Infrastructure in `platform/`

```go
// internal/platform/database/db.go
package database

type Config struct {
    DSN             string
    MaxOpenConns    int
    MaxIdleConns    int
    ConnMaxLifetime time.Duration
}

func New(cfg Config, logger *slog.Logger) (*sqlx.DB, error) {
    db, err := sqlx.Connect("postgres", cfg.DSN)
    if err != nil {
        return nil, fmt.Errorf("database.New: %w", err)
    }
    db.SetMaxOpenConns(cfg.MaxOpenConns)
    db.SetMaxIdleConns(cfg.MaxIdleConns)
    db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
    logger.Info("database connected")
    return db, nil
}
```

Each feature module receives `*sqlx.DB` via constructor injection from `main.go`. They never construct their own DB connections.

---

## Module Boundary Design

**Default: one module for everything.** Split only when:

| Reason to split                 | Example                              |
| ------------------------------- | ------------------------------------ |
| Independently versioned library | `github.com/myorg/myapp-sdk`         |
| Separate deployment unit        | A CLI tool distributed separately    |
| Third-party shared package      | `github.com/myorg/shared-middleware` |

**Rule:** Put everything in `internal/` by default. Only move to `pkg/` if you have a genuine external consumer.

**Anti-pattern: the "types" package**

```
// Do NOT do this
internal/
└── types/
    ├── user.go      ← everything imports this, nobody owns it
    ├── order.go
    └── product.go
```

Distribute types to their owning modules instead.

---

## Cross-Skill References

- For dependency graphs: see **[system-design-dependency-graphs.md](system-design-dependency-graphs.md)**
- For bounded context analysis: see **[system-design-bounded-contexts.md](system-design-bounded-contexts.md)**
- For C4 diagrams: see **[system-design-c4-model.md](system-design-c4-model.md)**
