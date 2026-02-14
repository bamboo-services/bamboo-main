# ROUTE KNOWLEDGE BASE

## OVERVIEW
Central HTTP wiring for Gin groups, middleware chain, and handler registration.
Routes are organized by domain (`auth`, `public`, `link`, `info`, `sponsor`, `admin`).

## WHERE TO LOOK
| Task | Location | Notes |
|---|---|---|
| Global middleware and `/api/v1` root | `internal/app/route/route.go` | response/cors/options middleware order is defined here |
| Auth endpoints | `internal/app/route/route_auth.go` | split public auth and auth-required subgroup |
| Admin endpoints + auth gate | `internal/app/route/route_admin.go` | OAuth check + local auth + role middleware |
| Public health/ping routes | `internal/app/route/route_public.go` | no auth middleware |
| Link/public link routes | `internal/app/route/route_link.go` | public link listing endpoint |
| Swagger mount behavior | `internal/app/route/route_swagger.go` | only register in debug mode |

## CONVENTIONS
- Instantiate handlers with `handler.NewHandler[T](r.context, "Name")`.
- Keep route files as wiring only; no request parsing or business branching.
- Keep protected routes under explicit middleware chain.
- Keep path naming stable (`/api/v1/...`) and domain-grouped.
- Register Swagger route only when debug env is enabled.

## ANTI-PATTERNS
- Adding business logic directly in route methods.
- Exposing admin routes without OAuth + `AuthMiddleware` + role check.
- Registering Swagger routes in production mode by default.
- Duplicating handler construction patterns inconsistently across route files.

## NOTES
- Route group root is `/api/v1` and should remain stable.
- `systemUserAdminRouter` and `systemLogRouter` are placeholders today.
