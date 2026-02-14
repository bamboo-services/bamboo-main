# REPOSITORY KNOWLEDGE BASE

## OVERVIEW
Repository layer encapsulates GORM access and cache invalidation policies.
Write paths generally invalidate Redis cache keys explicitly.

## WHERE TO LOOK
| Task | Location | Notes |
|---|---|---|
| Friend-link persistence | `internal/repository/link.go` | CRUD, list/filter, association cleanup, cache invalidation |
| Group/color persistence | `internal/repository/link_group.go`, `internal/repository/link_color.go` | status/sort/list/delete behaviors |
| Sponsor persistence | `internal/repository/sponsor_channel.go`, `internal/repository/sponsor_record.go` | sponsor channel/record queries |
| User/system persistence | `internal/repository/system_user.go`, `internal/repository/system.go` | user auth lookups and system config storage |
| Cache adapters | `internal/repository/cache/` | entity cache `Get/Set/Delete` wrappers |

## CONVENTIONS
- Constructors follow `New*Repo(db *gorm.DB, rdb *redis.Client)`.
- Methods accept optional transaction `tx *gorm.DB` and resolve DB via `pickDB(tx)`.
- Query context uses request context path consistently from caller.
- Cache pattern is cache-aside: read-through + invalidation on writes/deletes.
- Not-found behavior should be explicit (`found bool` or `RowsAffected` checks).

## ANTI-PATTERNS
- Business rule branching in repository methods.
- Skipping cache invalidation after mutable operations.
- Bypassing `pickDB(tx)` and breaking transaction consistency.
- Returning framework-specific response payloads from repository layer.

## NOTES
- `internal/repository/cache/` is a key companion namespace for this layer.
- Prefer explicit preload/select patterns over hidden query side-effects.
- Keep cache TTL/key decisions centralized with constants + cache adapters.
