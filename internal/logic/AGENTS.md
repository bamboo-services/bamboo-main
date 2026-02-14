# LOGIC KNOWLEDGE BASE

## OVERVIEW
Logic layer orchestrates domain rules, transactions, and async side-effects.
Most non-trivial project behavior is implemented here.

## WHERE TO LOOK
| Task | Location | Notes |
|---|---|---|
| Auth/session/password flow | `internal/logic/auth.go`, `internal/logic/auth_oauth.go` | OAuth sync, token/session behavior, reset/email flows |
| Link lifecycle and public listing | `internal/logic/link.go` | add/update/delete/status + notification trigger |
| Link taxonomy management | `internal/logic/link_group.go`, `internal/logic/link_color.go` | ordering/status/delete constraints |
| Sponsor domain logic | `internal/logic/sponsor_channel.go`, `internal/logic/sponsor_record.go` | channel/record orchestration and shaping |
| Public/system info | `internal/logic/public.go`, `internal/logic/info.go` | health/info read-write behavior |
| Cross-domain helpers | `internal/logic/helper/` | session/mail helper logic reused by domains |

## CONVENTIONS
- Constructors follow `New*Logic(ctx context.Context)` and initialize repos/helpers.
- Domain rules and transaction boundaries stay in logic, not in handlers/repositories.
- Use typed project constants for status/role/context semantics.
- Async triggers (mail/notifications) are non-blocking side-effects.
- Return `*xError.Error` consistently for upstream error middleware.

## ANTI-PATTERNS
- Moving DB query details from repository into logic.
- Writing response formatting concerns in logic layer.
- Starting ad-hoc goroutines without clear failure/isolation behavior.
- Mixing unrelated concerns into already dense core files without extraction.

## NOTES
- `auth.go` and `link.go` are current complexity hotspots in this layer.
- Keep side-effects isolated and observable (mail/session/cache updates).
