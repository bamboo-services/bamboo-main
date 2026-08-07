# 项目知识库

**生成日期:** 2026-08-07
**提交:** cb3413a
**分支:** master

## 概述
Gin + GORM + PostgreSQL + Redis 单体后端 + React 19 + TanStack Router 前端，用于友情链接管理。
后端运行时架构为 `Route -> Handler -> Logic -> Repository(+Cache)`，并附带截图任务子系统与 Prometheus 埋点；
前端为独立 Vite 子项目，file-based routing + shadcn/ui + 竹林水墨设计语言，产物内嵌进二进制实现单文件部署。
基于 `bamboo-base-go` v1.2.0（defined/common/major/plugins 四层架构，数据库驱动插件化）构建。

## 目录结构
```text
.
|- main.go                    # 进程入口；xOption 声明式配置（DB/Cache/Route）+ Runner + cron/worker 协程
|- go.mod                     # bamboo-base-go v1.2.0；beacon-sso-sdk 本地 replace
|- Makefile                   # dev/generate（单二进制构建）/swag 等快捷命令
|- frontend/                  # 独立前端子项目（pnpm + Vite + React 19 + TanStack）
|- resources/frontend/        # 前端构建产物内嵌（go:embed，产物不入库）
|- cmd/                       # 独立命令行工具
|  `- import-old-data/        #   旧数据迁移工具（pg_restore + COPY 解析）
|- internal/
|  |- app/startup/            # 邮件客户端 + 邮件业务配置 + 截图服务管理器 + SSO 节点（DB/Cache 由 xOption）
|  |- app/route/              # 路由注册 + 中间件链 + SPA 单二进制托管（RouteRegistrar 模式）
|  |- handler/                # 仅 HTTP bind/validate/respond
|  |- logic/                  # 业务编排、事务、异步触发
|  |- repository/             # DB 访问 + 经 xCache.Manager 的缓存失效 + Redis 会话/令牌
|  |- entity/                 # GORM 实体（嵌入 xModels.BaseEntity）
|  |- middleware/             # 鉴权与角色中间件
|  |- models/                 # 基础配置模型与 redis 缓存模型
|  |- service/screenshot/     # 友链站点截图任务子系统（队列 + rod 无头浏览器）
|  `- metrics/                # Prometheus 请求埋点 + 运行时指标 + /metrics
|- api/                       # 按领域分组的请求/响应 DTO
|- pkg/constants/             # 状态/角色/context/redis 键常量
|- pkg/util/                  # ctx 上下文工具与网络工具
|- docs/                      # 生成的 swagger 产物（禁止手改）
|- templates/mail/            # 邮件模板内容块（框架 xEmail 外部模板目录）
|- scripts/                   # 数据迁移 SQL 脚本
|- deploy/                    # 部署配置（Prometheus 等）
|- test/                      # SMTP E2E 测试（按 env 开启）
```

## 导航指南
| 任务 | 位置 | 说明 |
|---|---|---|
| 启动流程/启动失败 | `main.go`、`internal/app/startup/` | main 使用 xOption；startup 注册 Email 客户端 + 业务配置 + 截图服务 + SSO 节点 |
| 数据库/缓存配置 | `main.go`、`.env` | xOption.WithDatabase(FromEnv) + WithCache(FromEnv)；需 DATABASE_DRIVER/NOSQL_DRIVER 环境变量 + 驱动插件空白导入 |
| 种子数据（默认 admin/info） | `internal/app/startup/prepare/` | xOptionDB.WithPrepare 钩子：`prepare.DefaultUser` / `prepare.DefaultInfo`，AutoMigrate 之后运行 |
| 注册端点 | `internal/app/route/` | admin 路由应用 OAuth 校验 + 本地鉴权中间件；user 路由应用 AuthMiddleware |
| 新增 API 行为 | `internal/handler/` + `internal/logic/` | handler 保持瘦身；logic 持有规则 |
| DB 查询/事务 | `internal/repository/`、`internal/logic/*` | logic 开启 tx，repo 支持可选 `tx *gorm.DB` |
| 缓存接入 | 各 repo 构造器 + `pkg/constants/cache.go` | `xCache.KeyCacheOf[string, entity.X](m)`；键经 `RedisX.Get(id).String()` |
| Redis 会话/令牌 | `internal/repository/session.go`、`token.go` | 纯 Redis 仓储（KeyCache），认证/验证码/注册码存取 |
| Redis 键使用 | `pkg/constants/cache.go` | 键构建器自动追加 `NOSQL_PREFIX` + `:` |
| 邮件发送 | `internal/logic/mail.go` | 框架 xEmail 插件发送；调用点经 `xAsync.Async` 解耦请求上下文异步触发 |
| 邮件模板 | `templates/mail/*.html` | `{{define "名称"}}` 内容块，经 `EMAIL_TEMPLATE_DIR` 加载，套用框架 `_base.html` 布局 |
| 鉴权/会话问题 | `internal/logic/auth.go`、`internal/middleware/` | token/user context + OAuth 集成 |
| 仪表盘统计 | `internal/logic/dashboard.go`、`internal/handler/dashboard.go` | 友链计数与最近待审核申请 |
| 健康检查 | `internal/repository/health.go`、`internal/logic/public.go` | 数据库连通性探测 + 运行时指标 |
| Prometheus 埋点 | `internal/metrics/` | `/metrics` 端点 + 全局请求埋点中间件 + 运行时指标采集 |
| 友链截图服务 | `internal/service/screenshot/` | 队列 + rod 无头浏览器截图；cron 每日全量 + 手动重截 |
| 单二进制部署 | `resources/frontend/embed.go`、`internal/app/route/frontend.go` | 前端产物 go:embed；NoRoute 阶段 SPA fallback |
| 数据迁移 | `cmd/import-old-data/`、`scripts/` | 旧库导入工具与 SQL 迁移脚本 |
| 前端开发 | `frontend/` | pnpm + Vite，独立子项目 |
| 前端路由 | `frontend/src/routes/` | TanStack Router file-based，自动生成 routeTree |
| bamboo-base-go 包路径 | 见「备注」 | v1.2.0：xCache→major/cache、xError→common/error、xLog→common/log；驱动插件化→plugins/database/* |

## 代码地图
| 符号/区域 | 类型 | 位置 | 作用 |
|---|---|---|---|
| `main` | func | `main.go` | 进程装配；xOption 声明式配置 + Runner + cron/worker 附加协程 |
| `startup.Init` | func | `internal/app/startup/startup.go` | 返回 ctx + Email 客户端/业务配置/截图服务/SSO 自定义节点 |
| `prepare.DefaultUser` / `prepare.DefaultInfo` | func | `internal/app/startup/prepare/` | xOptionDB.PrepareFunc：种子默认 admin + 站点信息 |
| `route.NewRoute` | func | `internal/app/route/route.go` | RouteRegistrar：`func(ctx, serve)` 路由 + 中间件 + 子路由 |
| `handler.NewHandler[T]` | 泛型构造器 | `internal/handler/handler.go` | 注入全部 logic 依赖，`IHandler` 接口约束 |
| `screenshot.Manager` | 类型 | `internal/service/screenshot/manager.go` | 截图任务管理器（FIFO 队列 + 单 worker） |
| `metrics.Middleware` / `metrics.Handler` | func | `internal/metrics/metrics.go` | 请求埋点中间件 + `/metrics` promhttp 处理器 |
| `New*Repo` | 构造器 | `internal/repository/*.go` | 接受 `db *gorm.DB` + `m *xCache.Manager`，内部建 `kc`；纯 Redis 仓储只收 `m` |
| `router` / `Route` | 前端入口 | `frontend/src/main.tsx` | RouterProvider + createRouter(routeTree) |

## 模块架构
```text
后端
main.go (xOption 装配 DB/Cache/Route + Runner + cron/worker)
   |
   +-- startup.Init()      -> 自定义注册节点（Email 客户端 + 业务配置 + 截图服务 + SSO）
   +-- route.NewRoute       -> Gin 路由 + 中间件 + 领域子路由 + SPA 托管
   |       +-- handler      -> bind/validate/respond（瘦）
   |              +-- logic -> 业务规则 + 事务 + 异步触发（邮件经 xAsync 解耦 + 框架 xEmail 发送）
   |                     +-- repository (+xCache.Manager) -> GORM + cache-aside + Redis 会话/令牌
   +-- service/screenshot   -> 独立截图子系统（Queue + rod Capture + 存储）
   +-- metrics              -> Prometheus 埋点（全局中间件 + /metrics）

前端（独立子项目，产物内嵌）
frontend/ (pnpm + Vite + React 19 + TanStack)
   +-- src/main.tsx           -> RouterProvider + createRouter + QueryClient
   +-- src/routeTree.gen.ts   -> @tanstack/router-plugin 自动生成
   +-- src/routes/            -> file-based routing（_admin / _user / _authorization / _public / about / operate）
   +-- src/api/ + src/hooks/  -> 请求层 + TanStack Query 数据层
   +-- src/components/        -> shadcn/ui + 竹林水墨原语 + 领域组件
   +-- src/lib/               -> 纯逻辑（颜色视觉/分级/会话等）
   `-- resources/frontend/dist（go:embed 内嵌）
```
后端领域横向切分：`auth` / `link`（含 group/color）/ `sponsor`（channel/record）/ `info` / `public` / `user`。
SSO/OAuth 能力由 `beacon-sso-sdk`（本地 replace）提供，logic 与 route 分别集成 `bSdkLogic` / `bSdkRoute`。

## 约定
- handler 错误流：bind → `xValid.HandleValidationError` → logic 调用 → `_ = c.Error(err)` → return
- 成功响应优先使用 `xResult.SuccessHasData`，面向用户文案保持中文
- repository 接受 `tx *gorm.DB`；使用局部 `pickDB(tx)` 模式解析目标 DB
- 缓存策略为 cache-aside：经 `xCache.KeyCacheOf[string, entity.X](m)` 创建 `KeyCache`；写路径显式 `kc.Delete`
- 缓存键经 `constants.RedisX.Get(id).String()` 构造（带 `NOSQL_PREFIX` 前缀）；TTL 用 `xCache.WithTTL(...)`
- 配置来源 env-first（`.env` 由基础库加载）；`configs/` 存在但为空
- `xMain.Runner(reg, log, opts, goroutine...)` 配合 `[]xOption.Option` 声明式装配；附加协程（cron/截图 worker）随 Runner 生命周期启停
- 数据库装配需按 `DATABASE_DRIVER` 空白导入 `plugins/database/*` 驱动插件（v1.2.0）
- `xCtxUtil` 自 `major/utility/context` 导入（v1.1.2 保留 major/utility 路径，v1.2.0 沿用）
- `xValid`（HandleValidationError）自 `major/validator` 导入
- `xError` 自 `common/error` 导入；`xLog` / `xLogGorm` 自 `common/log` 导入（v1.1.2 保留 common/log，v1.2.0 沿用）
- `xCache` / `xCache.Manager` / `KeyCacheOf` 自 `major/cache` 导入（v1.1.0 起缓存层迁移至此）

## 反模式
- 手改 `docs/` 下生成的 swagger 文件，或手改 `frontend/src/routeTree.gen.ts`
- 在 handler 方法里写业务规则
- 越过 logic 层从 handler 直接调 repository
- 不接入 `startup_config.go` 就引入新配置源
- 盲信 README 默认值；以 startup 代码为真相源
- 直接持有 `*redis.Client` 做实体缓存——应统一走 `xCache.Manager` 泛型接口
- 用 bun/npm/yarn 替代 pnpm 管理前端依赖
- 使用声明式数据库装配却遗漏 `plugins/database/*` 驱动插件的空白导入（v1.2.0 起启动时会在 DatabaseInit 阶段报「不支持的数据库驱动」）
- 手改 `resources/frontend/dist` 构建产物——产物不入库，由 `make build-frontend` 生成后 go:embed 内嵌
- 手写代码文件缺失版权横幅，或擅自改写横幅中的作者名、年份、网址、许可证声明文字

## 独特风格
- 领域 DTO 放在顶层 `api/` 包树，而非 `internal/` 下
- 代码内默认管理员种子当前为 `xiao_lfeng`（见 prepare 模块）
- 前后端同仓库但前端是完全独立的 Vite 子项目（独立 `package.json` / `pnpm-lock.yaml` / tsconfig），产物经 go:embed 内嵌实现单二进制部署
- Entity 统一嵌入 `xModels.BaseEntity`（SnowflakeID 主键 + 时间戳），通过 `GetGene()` 声明业务基因
- 前端统一「竹林水墨」设计语言（淡绿宣纸 + 墨色衬线），完整 spec 沉淀于 `frontend/ARCHITECTURE.md`
- 截图服务采用「接口抽象 + 队列 + 单 worker」模式，`CaptureFunc` 可注入测试替身

## 版权注入
**所有手写的代码文件（`.go` / `.ts` / `.tsx` / `.js` / `.jsx` / `.css` / `.html` 等）必须在文件头部加入版权横幅**，置于 package 声明 / import / 其他任何代码之前。生成产物（`docs/` 下 swagger、`*.gen.ts`、`routeTree.gen.ts` 等自动生成文件）与配置文件（`*.json` / `*.toml` / `.env*`）豁免。

### 模板（逐字复制，禁止改写）
```text
// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------
```

### 适配规则
- **Go 文件**：使用 `//` 行注释，置于 `package` 声明之前；空行分隔后接 `package xxx`
- **TS/JS/CSS/HTML 文件**：将每行前缀 `//` 改为对应语言注释符（`/* */` 或 `<!-- -->`），正文逐字保留；置于文件首行
- **年份字段**：`2016-NOW(至今)` 与 `2016-2025` 为固定文案，不得替换为具体年份或当前年份
- **禁止变体**：不得改写作者名、网址、许可证声明文字；不得省略分隔线
- **完整性**：横幅必须完整出现，禁止只保留一两行摘要

### 反例
```text
// ❌ 缺失横幅
package handler

// ❌ 自行改写
// Author: 竹子
// 2024

// ✅ 标准
// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// ...
```

## 常用命令
```bash
# 单二进制一键构建（前端打包 → Swagger → Go 编译，产物内嵌前端）
make generate        # 或 make build
make install         # 首次：安装前端依赖（frontend/pnpm install）

# 后端开发
make dev             # 生成 Swagger 后运行（跳过前端构建）
make dev-backend     # 构建（含前端）并运行
make run             # 运行已构建的二进制
make dev-frontend    # 仅前端开发服务器（端口 3000）
go run main.go

# 后端质量
go test ./...
go fmt ./...
go vet ./...
make lint            # golangci-lint（未安装则跳过）

# 后端文档
make swag
swag init -g main.go --parseDependency

# 前端（frontend/ 目录下）
pnpm install
pnpm dev           # vite --port 3000
pnpm build         # vite build && tsc
pnpm test          # vitest run
pnpm check         # prettier --write . && eslint --fix
```

## 备注
- 依赖的本地服务：PostgreSQL（`DATABASE_*`）与 Redis（`NOSQL_*`）；截图服务按需 Chrome/Chromium（`SCREENSHOT_*`）
- 邮件 E2E 测试按需开启：`ENABLE_SMTP_E2E_TEST=true` 并配置 SMTP 环境变量
- `go.mod` 使用本地 `replace` 将 `beacon-sso-sdk` 指向 `/Users/xiaolfeng/ProgramProjects/Cooperate/phalanx/beacon-sso-sdk`
- bamboo-base-go 升级至 v1.2.0（破坏性变更：数据库驱动插件化，框架不再内置 GORM 驱动，需按 `DATABASE_DRIVER` 空白导入 `plugins/database/*` 对应驱动插件；驱动枚举迁移至 `common/database`（xDB），`xOptDatabase.Driver*` 保留别名重导出，`WithDialector` 支持注入自定义驱动）
- bamboo-base-go v1.1.2 历史（破坏性变更：缓存层迁移至 `major/cache` 的 `xCache.Manager` + 泛型 `KeyCacheOf`；`internal/repository/cache/` 目录已删除）
- beacon-sso-sdk 本地源已适配 v1.2.0：`xCache.Cache` 替换为本地 `Cache` 结构体，`xCtxUtil` 路径迁移
- 前端构建产物（`resources/frontend/dist`）**不入库**，仅保留 `.gitkeep` 占位；未执行 `make build-frontend` 时静态资源不可用（FileServer 返回 404）
- 截图服务 worker 由 `main.go` Runner 附加协程启动，cron 每日 0 点（Asia/Shanghai）全量刷新截图
- 前端固定使用 pnpm（`pnpm-lock.yaml` 为真相源），开发服务器监听 3000 端口，API 经 Vite proxy 转发到 5555
- 路由树由 `@tanstack/router-plugin` 在 Vite 编译期自动生成 `frontend/src/routeTree.gen.ts`

## 引用
- [frontend](./frontend/AGENTS.md) — 前端子项目（React 19 + TanStack + shadcn/ui + 竹林水墨）
- [internal/app/startup](./internal/app/startup/AGENTS.md) — 启动引导与种子钩子
- [internal/app/route](./internal/app/route/AGENTS.md) — 路由装配、中间件链与 SPA 托管
- [internal/handler](./internal/handler/AGENTS.md) — HTTP bind/validate/respond
- [internal/logic](./internal/logic/AGENTS.md) — 业务编排与事务
- [internal/repository](./internal/repository/AGENTS.md) — DB 访问、缓存失效与 Redis 会话/令牌
- [internal/service/screenshot](./internal/service/screenshot/AGENTS.md) — 友链站点截图子系统
