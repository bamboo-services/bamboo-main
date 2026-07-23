# TASK 知识库

## 概述
异步邮件投递管线实现层，基于 Redis 队列 + 重试 zset；包含 worker 编排与 SMTP/TLS 连接池行为。

## 目录结构
```text
internal/task/
|- mail.go                     # worker 生命周期与队列循环
|- mail_pool.go                # SMTP 连接池（TLS/STARTTLS）
```

## 导航指南
| 任务 | 位置 | 说明 |
|---|---|---|
| worker 生命周期与队列循环 | `mail.go` | worker goroutine、BRPop 消费、重试调度 |
| SMTP 连接池 | `mail_pool.go` | 连接复用、TLS/STARTTLS 处理 |
| 启动 runner 集成 | `internal/app/startup/worker/worker_mail.go` | 启动时拉起 worker，ctx 取消时停止 |
| 队列键定义 | `pkg/constants/cache.go` | `RedisMailQueue`/`RedisMailRetry`/`RedisMailFailed`/`RedisMailStats` |

## 约定
- 队列消费：主 list 用 BRPop 弹出，重试经 zset 按到期时间戳调度
- 重试策略：指数退避 + 抖动 + 最大重试次数保护
- worker 数量与超时取自邮件配置，带安全默认值
- SMTP 模式支持 TLS 直连与 STARTTLS；端口可自动推断模式
- 停止路径必须优雅：取消 context → waitgroup 收尾 → 关闭连接池

## 反模式
- 在请求 handler 中阻塞式发邮件——必须走队列
- 绕过共享 Redis 键常量自行拼 key——破坏调度器与队列一致性
- 修改重试语义而不同步调度器与队列迁移逻辑——丢消息
- 跳过关闭流程导致 goroutine/socket 泄漏

## 调试路径
1. 邮件不投递：先查 `RedisMailQueue` 是否有积压，再看 worker 是否启动（`worker_mail.go` 日志）
2. 重试异常：检查 `RedisMailRetry` zset 分数（到期时间戳）与调度器轮询节奏
3. 连接失败：核对 `mail_pool.go` 的 TLS/STARTTLS 模式与 `EMAIL_USE_TLS`/`EMAIL_USE_STARTTLS` 环境变量
4. 启动/停止不干净：确认 `worker_mail.go` 的 context 收到取消信号、waitgroup 是否 drain
