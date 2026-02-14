# STARTUP KNOWLEDGE BASE

## OVERVIEW
Bootstraps runtime dependencies before serving traffic.
Node order is fixed: config -> database -> redis -> prepare -> SSO.

## WHERE TO LOOK
| Task | Location | Notes |
|---|---|---|
| Startup registration order | `internal/app/startup/startup.go` | `Init()` appends nodes in execution order |
| Env-to-config mapping | `internal/app/startup/startup_config.go` | all runtime config comes from env/defaults |
| DB connect + migrate | `internal/app/startup/startup_database.go` | `AutoMigrate` runs on startup |
| Redis connect | `internal/app/startup/startup_redis.go` | ping check with timeout |
| Seed default data | `internal/app/startup/startup_prepare.go`, `internal/app/startup/prepare/` | idempotent default user/info setup |
| Worker lifecycle hook | `internal/app/startup/worker/worker_mail.go` | start on runner boot, stop on `ctx.Done()` |

## CONVENTIONS
- Keep startup nodes side-effect focused and deterministic.
- Read config from startup context via `ContextCustomConfig`, not ad-hoc globals.
- Keep prepare tasks idempotent; repeat startup must not duplicate seed data.
- Fail fast on dependency init errors (db/redis/config).
- Use defaults in config init only for local/dev-safe fallback.

## ANTI-PATTERNS
- Reordering startup nodes without dependency audit.
- Adding new config sources outside `startup_config.go`.
- Putting domain business rules in startup/prepare code.
- Running long-lived business loops here instead of worker/task modules.
- Treating README values as source of truth when startup code differs.

## NOTES
- `configs/` exists in repo but runtime config is env-driven.
- DB migration is automatic on boot; schema changes are startup-sensitive.
- Prepare step persists default admin marker in `system.admin.id`.
