# PROJECT KNOWLEDGE BASE

**Generated:** 2026-02-14T12:32:39+0800
**Commit:** 6c4731c
**Branch:** master

## OVERVIEW
Gin + GORM + PostgreSQL + Redis monolith for friend-link management.
Runtime architecture in this repo is `Route -> Handler -> Logic -> Repository(+Cache)`.

## STRUCTURE
```text
.
|- main.go                    # process entry; wires startup + routes + mail worker
|- internal/
|  |- app/startup/            # config/db/redis bootstrap and default data prepare
|  |- app/route/              # route registration + middleware chain
|  |- handler/                # HTTP bind/validate/respond only
|  |- logic/                  # business orchestration, transactions, async triggers
|  |- repository/             # DB access + Redis cache invalidation
|  |- entity/                 # GORM entities (Snowflake ID hooks)
|  |- task/                   # async mail queue worker + TLS pool
|- api/                       # request/response DTOs by domain
|- pkg/constants/             # status/role/context/redis key constants
|- docs/                      # generated swagger artifacts (DO NOT hand-edit)
|- test/                      # SMTP E2E tests (opt-in by env)
```

## WHERE TO LOOK
| Task | Location | Notes |
|---|---|---|
| Start flow / boot failures | `main.go`, `internal/app/startup/` | init order: config -> db -> redis -> prepare -> SSO |
| Register endpoints | `internal/app/route/` | admin routes apply OAuth check + local auth middleware |
| Add API behavior | `internal/handler/` + `internal/logic/` | handlers stay thin; logic owns rules |
| DB query / transaction | `internal/repository/`, `internal/logic/*` | logic begins tx, repo supports optional `tx *gorm.DB` |
| Redis key usage | `pkg/constants/cache.go` | key builder auto-adds `NOSQL_PREFIX` + `:` |
| Mail async pipeline | `internal/task/mail.go`, `internal/task/mail_pool.go` | queue, retry zset, backoff, TLS/STARTTLS |
| Auth/session issues | `internal/logic/auth.go`, `internal/middleware/` | token/user context + OAuth integration |

## CODE MAP
| Symbol/Area | Type | Location | Role |
|---|---|---|---|
| `main` | func | `main.go` | process composition and runner entry |
| `startup.Init` | func | `internal/app/startup/startup.go` | startup node registry |
| `route.NewRoute` | func | `internal/app/route/route.go` | router + middleware + subrouters |
| `handler.NewHandler[T]` | generic ctor | `internal/handler/handler.go` | injects all logic dependencies |
| `New*Logic` | constructors | `internal/logic/*.go` | business services per domain |
| `New*Repo` | constructors | `internal/repository/*.go` | database and cache access |
| `MailWorkerRunner` | worker entry | `internal/app/startup/worker/worker_mail.go` | starts/stops async mail worker |

## CONVENTIONS
- Error flow in handlers: bind -> `xValid.HandleValidationError` -> logic call -> `_ = c.Error(err)` -> return.
- Success responses prefer `xResult.SuccessHasData` with Chinese user-facing message.
- Repositories accept `tx *gorm.DB`; use local `pickDB(tx)` pattern.
- Cache strategy is cache-aside with explicit invalidation on writes.
- Config source is env-first (`.env` loaded by base library); `configs/` exists but is empty.

## ANTI-PATTERNS (THIS PROJECT)
- Editing generated swagger files in `docs/`.
- Putting business rules inside handler methods.
- Bypassing logic layer and calling repository directly from handlers.
- Introducing new config source paths without wiring `startup_config.go`.
- Assuming README defaults are always current; check startup code as source of truth.

## UNIQUE STYLES
- File headers include project copyright banner on most handwritten Go files.
- Domain DTOs live in top-level `api/` package tree (not under `internal/`).
- Default admin seed in code is currently `xiao_lfeng` (see prepare module).

## COMMANDS
```bash
# dev
make dev
go run main.go

# quality
go test ./...
go fmt ./...
go vet ./...

# docs
make swag
swag init -g main.go --parseDependency
```

## NOTES
- Required local services: PostgreSQL (`DATABASE_*`) and Redis (`NOSQL_*`).
- Mail E2E tests are opt-in: `ENABLE_SMTP_E2E_TEST=true` with SMTP env vars.
- `go.mod` uses local `replace` paths for `bamboo-base-go` and `beacon-sso-sdk`.
