# 项目知识库

**生成日期:** 2026-07-21
**提交:** 06a649f
**分支:** master

## 概述
Gin + GORM + PostgreSQL + Redis 单体，用于友情链接管理。
运行时架构为 `Route -> Handler -> Logic -> Repository(+Cache)`。
基于 `bamboo-base-go` v1.0.4（defined/common/major/plugins 四层架构）构建。

## 目录结构
```text
.
|- main.go                    # 进程入口；xOption 声明式配置（DB/Cache/Route）+ Runner + mail worker
|- internal/
|  |- app/startup/            # 仅邮件配置 + SSO 节点（DB/Cache 由 xOption）；prepare 钩子用于种子数据
|  |- app/route/              # 路由注册 + 中间件链（RouteRegistrar 模式）
|  |- handler/                # 仅 HTTP bind/validate/respond
|  |- logic/                  # 业务编排、事务、异步触发
|  |- repository/             # DB 访问 + Redis 缓存失效
|  |- entity/                 # GORM 实体（Snowflake ID 钩子）
|  |- task/                   # 异步邮件队列 worker + TLS 连接池
|  |- middleware/             # 鉴权与角色中间件
|  |- models/                 # 基础配置模型与 redis 缓存模型
|- api/                       # 按领域分组的请求/响应 DTO
|- pkg/constants/             # 状态/角色/context/redis 键常量
|- docs/                      # 生成的 swagger 产物（禁止手改）
|- test/                      # SMTP E2E 测试（按 env 开启）
```

## 导航指南
| 任务 | 位置 | 说明 |
|---|---|---|
| 启动流程/启动失败 | `main.go`、`internal/app/startup/` | main 使用 xOption；startup 仅注册 Email 配置 + SSO 节点 |
| 数据库/缓存配置 | `main.go`、`.env` | xOption.WithDatabase(FromEnv) + WithCache(FromEnv)；需 DATABASE_DRIVER/NOSQL_DRIVER 环境变量 |
| 种子数据（默认 admin/info） | `internal/app/startup/prepare/` | xOptionDB.WithPrepare 钩子，在 AutoMigrate 之后运行 |
| 注册端点 | `internal/app/route/` | admin 路由应用 OAuth 校验 + 本地鉴权中间件 |
| 新增 API 行为 | `internal/handler/` + `internal/logic/` | handler 保持瘦身；logic 持有规则 |
| DB 查询/事务 | `internal/repository/`、`internal/logic/*` | logic 开启 tx，repo 支持可选 `tx *gorm.DB` |
| Redis 键使用 | `pkg/constants/cache.go` | 键构建器自动追加 `NOSQL_PREFIX` + `:` |
| 邮件异步管线 | `internal/task/mail.go`、`internal/task/mail_pool.go` | queue、retry zset、退避、TLS/STARTTLS |
| 鉴权/会话问题 | `internal/logic/auth.go`、`internal/middleware/` | token/user context + OAuth 集成 |
| bamboo-base-go 包路径 | 见「备注」 | v1.0.4 迁移：xCtxUtil→major、xValid→major、xLogGorm→major/log |

## 代码地图
| 符号/区域 | 类型 | 位置 | 作用 |
|---|---|---|---|
| `main` | func | `main.go` | 进程装配；xOption 声明式配置 + Runner 入口 |
| `startup.Init` | func | `internal/app/startup/startup.go` | 返回 ctx + Email/SSO 自定义节点（DB/Cache 由 xOption） |
| `prepare.DefaultData` | func | `internal/app/startup/prepare/prepare.go` | xOptionDB.PrepareFunc：种子默认 admin + 站点信息 |
| `route.NewRoute` | func | `internal/app/route/route.go` | RouteRegistrar：`func(ctx, serve)` 路由 + 中间件 + 子路由 |
| `handler.NewHandler[T]` | 泛型构造器 | `internal/handler/handler.go` | 注入全部 logic 依赖 |
| `New*Logic` | 构造器 | `internal/logic/*.go` | 各领域业务服务 |
| `New*Repo` | 构造器 | `internal/repository/*.go` | 数据库与缓存访问 |
| `MailWorkerRunner` | worker 入口 | `internal/app/startup/worker/worker_mail.go` | 启停异步邮件 worker |

## 模块架构
```text
main.go (xOption 装配 DB/Cache/Route + Runner)
   |
   +-- startup.Init()      -> 自定义注册节点（Email 配置 + SSO）
   +-- route.NewRoute       -> Gin 路由 + 中间件 + 领域子路由
   |       +-- handler      -> bind/validate/respond（瘦）
   |              +-- logic -> 业务规则 + 事务 + 异步触发
   |                     +-- repository (+cache/) -> GORM + Redis cache-aside
   +-- worker.MailWorkerRunner -> 异步邮件队列 worker
```
领域横向切分：`auth` / `link`（含 group/color）/ `sponsor`（channel/record）/ `info` / `public`。
SSO/OAuth 能力由 `beacon-sso-sdk`（本地 replace）提供，logic 与 route 分别集成 `bSdkLogic` / `bSdkRoute`。

## 约定
- handler 错误流：bind → `xValid.HandleValidationError` → logic 调用 → `_ = c.Error(err)` → return
- 成功响应优先使用 `xResult.SuccessHasData`，面向用户文案保持中文
- repository 接受 `tx *gorm.DB`；使用局部 `pickDB(tx)` 模式解析目标 DB
- 缓存策略为 cache-aside，写路径显式失效
- 配置来源 env-first（`.env` 由基础库加载）；`configs/` 存在但为空
- `xMain.Runner(reg, log, opts, goroutine...)` 配合 `[]xOption.Option` 声明式装配（v1.0.4）
- `xCtxUtil` 自 `major/utility/context` 导入（v1.0.4 后不再从 `common/utility/context`）
- `xValid`（HandleValidationError）自 `major/validator` 导入（v1.0.4 后不再从 `common/validator`）
- GORM logger 使用 `xLogGorm`，来自 `major/log`（v1.0.4 后不再用 `common/log` 的 `xLog`）

## 反模式
- 手改 `docs/` 下生成的 swagger 文件
- 在 handler 方法里写业务规则
- 越过 logic 层从 handler 直接调 repository
- 不接入 `startup_config.go` 就引入新配置源
- 盲信 README 默认值；以 startup 代码为真相源
- 使用旧 bamboo-base-go 包路径（`common/utility/context`、`common/validator`、`major/http`）——v1.0.4 已迁移

## 独特风格
- 多数手写 Go 文件头部包含项目版权横幅
- 领域 DTO 放在顶层 `api/` 包树，而非 `internal/` 下
- 代码内默认管理员种子当前为 `xiao_lfeng`（见 prepare 模块）

## 常用命令
```bash
# 开发
make dev
go run main.go

# 质量
go test ./...
go fmt ./...
go vet ./...

# 文档
make swag
swag init -g main.go --parseDependency
```

## 备注
- 依赖的本地服务：PostgreSQL（`DATABASE_*`）与 Redis（`NOSQL_*`）
- 邮件 E2E 测试按需开启：`ENABLE_SMTP_E2E_TEST=true` 并配置 SMTP 环境变量
- `go.mod` 使用本地 `replace` 将 `beacon-sso-sdk` 指向 `/Users/xiaolfeng/ProgramProjects/Cooperate/phalanx/beacon-sso-sdk`
- bamboo-base-go 升级至 v1.0.4（破坏性变更：RouteRegistrar 签名、xOption Runner、包路径迁移）
- beacon-sso-sdk 本地源已适配 v1.0.4：`xCache.Cache` 替换为本地 `Cache` 结构体，`xCtxUtil` 路径迁移
- Entity 文件使用手写 `BeforeCreate` 钩子（非 `xModels.BaseEntity`）；兼容 v1.0.4，可后续改用 base entity 模式

## 引用
- [internal/app/startup](./internal/app/startup/AGENTS.md) — 启动引导与种子/worker 钩子
- [internal/app/route](./internal/app/route/AGENTS.md) — 路由装配与中间件链
- [internal/handler](./internal/handler/AGENTS.md) — HTTP bind/validate/respond
- [internal/logic](./internal/logic/AGENTS.md) — 业务编排与事务
- [internal/repository](./internal/repository/AGENTS.md) — DB 访问与缓存失效
- [internal/task](./internal/task/AGENTS.md) — 异步邮件管线与连接池
