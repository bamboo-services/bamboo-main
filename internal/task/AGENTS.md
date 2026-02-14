# TASK KNOWLEDGE BASE

## OVERVIEW
Implements async mail delivery pipeline backed by Redis queue + retry set.
Contains worker orchestration and SMTP/TLS connection pool behavior.

## WHERE TO LOOK
| Task | Location | Notes |
|---|---|---|
| Worker lifecycle and queue loop | `internal/task/mail.go` | worker goroutines, BRPop consume, retry scheduler |
| SMTP connection pooling | `internal/task/mail_pool.go` | pooled connections, TLS/STARTTLS handling |
| Startup runner integration | `internal/app/startup/worker/worker_mail.go` | starts worker on boot and stops on context cancel |
| Queue key definitions | `pkg/constants/cache.go` | `RedisMailQueue`, `RedisMailRetry`, `RedisMailFailed`, `RedisMailStats` |

## CONVENTIONS
- Queue ingestion: pop from main list, retry via sorted-set by due timestamp.
- Retry policy uses exponential backoff with jitter and max retry guard.
- Worker count/timeouts come from email config with safe defaults.
- SMTP mode supports TLS direct and STARTTLS; mode can be auto-inferred by port.
- Stop path must be graceful: cancel context, waitgroup drain, close pool.

## ANTI-PATTERNS
- Blocking request handlers on email send path.
- Writing mail queue code that bypasses shared Redis key constants.
- Changing retry semantics without keeping scheduler + queue moves in sync.
- Leaking goroutines or sockets by skipping shutdown/close flow.

## NOTES
- Queue + retry consistency depends on atomic move semantics in scheduler.
- `mail_pool.go` is concurrency-heavy; small changes require careful review.
- SMTP env flags (`EMAIL_USE_TLS`, `EMAIL_USE_STARTTLS`) influence runtime mode.
- Missing SMTP config should disable worker safely rather than fail boot.
