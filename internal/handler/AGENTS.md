# HANDLER KNOWLEDGE BASE

## OVERVIEW
Handler layer owns HTTP bind/validate/respond flow only.
Business decisions must stay in `internal/logic/`.

## WHERE TO LOOK
| Task | Location | Notes |
|---|---|---|
| Shared handler construction | `internal/handler/handler.go` | generic `NewHandler[T]` injects all logic deps |
| Auth HTTP endpoints | `internal/handler/auth.go` | login/register/password/reset/verify endpoints |
| Link and group endpoints | `internal/handler/link.go`, `internal/handler/link_group.go`, `internal/handler/link_color.go` | bind + logic call + response format |
| Sponsor endpoints | `internal/handler/sponsor_channel.go`, `internal/handler/sponsor_record.go` | admin/public sponsor operations |
| Public/Info endpoints | `internal/handler/public.go`, `internal/handler/info.go` | health/info read/update endpoints |

## CONVENTIONS
- Validation flow: `ShouldBind*` -> `xValid.HandleValidationError` -> return.
- Logic errors: `_ = c.Error(err)` then return; let middleware render error payload.
- Success responses: prefer `xResult.SuccessHasData` with user-facing Chinese message.
- Handler methods should not manage DB transactions or cache operations.
- Keep Swagger annotations aligned with actual route and response behavior.

## ANTI-PATTERNS
- Embedding business rules in handlers.
- Calling repository methods from handlers directly.
- Returning custom ad-hoc JSON shape instead of `xResult` conventions.
- Skipping validation helper and manually formatting bind errors inconsistently.

## NOTES
- Handler files are long partly due to Swagger annotations.
- Keep message text user-facing and Chinese to match current API style.
- If behavior changes, update both handler annotations and route mapping.
