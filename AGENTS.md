# 项目知识库

**生成日期:** 2026-07-25
**提交:** a40b4d6
**分支:** master

## 概述
Gin + GORM + PostgreSQL + Redis 单体后端 + React 19 + TanStack Router 前端，用于友情链接管理。
后端运行时架构为 `Route -> Handler -> Logic -> Repository(+Cache)`；前端为独立 Vite 子项目，
file-based routing + shadcn/ui。基于 `bamboo-base-go` v1.1.2（defined/common/major/plugins 四层架构）构建。

## 目录结构
```text
.
|- main.go                    # 进程入口；xOption 声明式配置（DB/Cache/Route）+ Runner
|- go.mod                     # bamboo-base-go v1.1.2；beacon-sso-sdk 本地 replace
|- Makefile                   # dev/swag 等快捷命令
|- frontend/                  # 独立前端子项目（pnpm + Vite + React 19 + TanStack）
|- internal/
|  |- app/startup/            # 邮件客户端 + 邮件业务配置 + SSO 节点（DB/Cache 由 xOption）；prepare 钩子用于种子数据
|  |- app/route/              # 路由注册 + 中间件链（RouteRegistrar 模式）
|  |- handler/                # 仅 HTTP bind/validate/respond
|  |- logic/                  # 业务编排、事务、异步触发
|  |- repository/             # DB 访问 + 经 xCache.Manager 的缓存失效
|  |- entity/                 # GORM 实体（Snowflake ID 钩子）
|  |- middleware/             # 鉴权与角色中间件
|  |- models/                 # 基础配置模型与 redis 缓存模型
|- api/                       # 按领域分组的请求/响应 DTO
|- pkg/constants/             # 状态/角色/context/redis 键常量
|- docs/                      # 生成的 swagger 产物（禁止手改）
|- templates/mail/            # 邮件模板内容块（框架 xEmail 外部模板目录）
|- test/                      # SMTP E2E 测试（按 env 开启）
```

## 导航指南
| 任务 | 位置 | 说明 |
|---|---|---|
| 启动流程/启动失败 | `main.go`、`internal/app/startup/` | main 使用 xOption；startup 注册 Email 客户端 + Email 业务配置 + SSO 节点 |
| 数据库/缓存配置 | `main.go`、`.env` | xOption.WithDatabase(FromEnv) + WithCache(FromEnv)；需 DATABASE_DRIVER/NOSQL_DRIVER 环境变量 |
| 种子数据（默认 admin/info） | `internal/app/startup/prepare/` | xOptionDB.WithPrepare 钩子，在 AutoMigrate 之后运行 |
| 注册端点 | `internal/app/route/` | admin 路由应用 OAuth 校验 + 本地鉴权中间件 |
| 新增 API 行为 | `internal/handler/` + `internal/logic/` | handler 保持瘦身；logic 持有规则 |
| DB 查询/事务 | `internal/repository/`、`internal/logic/*` | logic 开启 tx，repo 支持可选 `tx *gorm.DB` |
| 缓存接入 | 各 repo 构造器 + `pkg/constants/cache.go` | `xCache.KeyCacheOf[string, entity.X](m)`；键经 `RedisX.Get(id).String()` |
| Redis 键使用 | `pkg/constants/cache.go` | 键构建器自动追加 `NOSQL_PREFIX` + `:` |
| 邮件发送 | `internal/logic/mail.go` | 框架 xEmail 插件发送；调用点经 `xAsync.Async` 解耦请求上下文异步触发，`SendTemplate` 内 `xCtxUtil.GetEmailClient(ctx)` |
| 邮件模板 | `templates/mail/*.html` | `{{define "名称"}}` 内容块，经 `EMAIL_TEMPLATE_DIR` 加载，套用框架 `_base.html` 布局 |
| 鉴权/会话问题 | `internal/logic/auth.go`、`internal/middleware/` | token/user context + OAuth 集成 |
| 前端开发 | `frontend/` | pnpm + Vite，独立子项目 |
| 前端路由 | `frontend/src/routes/` | TanStack Router file-based，自动生成 routeTree |
| bamboo-base-go 包路径 | 见「备注」 | v1.1.2：xCache→major/cache、xError→common/error、xLog→common/log |

## 代码地图
| 符号/区域 | 类型 | 位置 | 作用 |
|---|---|---|---|
| `main` | func | `main.go` | 进程装配；xOption 声明式配置 + Runner 入口 |
| `startup.Init` | func | `internal/app/startup/startup.go` | 返回 ctx + Email 客户端/业务配置/SSO 自定义节点（DB/Cache 由 xOption） |
| `prepare.DefaultData` | func | `internal/app/startup/prepare/prepare.go` | xOptionDB.PrepareFunc：种子默认 admin + 站点信息 |
| `route.NewRoute` | func | `internal/app/route/route.go` | RouteRegistrar：`func(ctx, serve)` 路由 + 中间件 + 子路由 |
| `handler.NewHandler[T]` | 泛型构造器 | `internal/handler/handler.go` | 注入全部 logic 依赖 |
| `New*Repo` | 构造器 | `internal/repository/*.go` | 接受 `db *gorm.DB` + `m *xCache.Manager`，内部建 `kc` |
| `router` / `Route` | 前端入口 | `frontend/src/main.tsx` | RouterProvider + createRouter(routeTree) |

## 模块架构
```text
后端
main.go (xOption 装配 DB/Cache/Route + Runner)
   |
   +-- startup.Init()      -> 自定义注册节点（Email 客户端 + Email 业务配置 + SSO）
   +-- route.NewRoute       -> Gin 路由 + 中间件 + 领域子路由
   |       +-- handler      -> bind/validate/respond（瘦）
   |              +-- logic -> 业务规则 + 事务 + 异步触发（邮件经 xAsync 解耦 + 框架 xEmail 发送）
   |                     +-- repository (+xCache.Manager) -> GORM + cache-aside

前端（独立子项目）
frontend/ (pnpm + Vite + React 19)
   +-- src/main.tsx           -> RouterProvider + createRouter
   +-- src/routeTree.gen.ts   -> @tanstack/router-plugin 自动生成
   +-- src/routes/            -> file-based routing（_admin / _public / _authorization）
   +-- src/components/ui/     -> shadcn/ui 组件库
   +-- src/components/layout/ -> 布局壳（admin-sidebar）
```
后端领域横向切分：`auth` / `link`（含 group/color）/ `sponsor`（channel/record）/ `info` / `public`。
SSO/OAuth 能力由 `beacon-sso-sdk`（本地 replace）提供，logic 与 route 分别集成 `bSdkLogic` / `bSdkRoute`。

## 约定
- handler 错误流：bind → `xValid.HandleValidationError` → logic 调用 → `_ = c.Error(err)` → return
- 成功响应优先使用 `xResult.SuccessHasData`，面向用户文案保持中文
- repository 接受 `tx *gorm.DB`；使用局部 `pickDB(tx)` 模式解析目标 DB
- 缓存策略为 cache-aside：经 `xCache.KeyCacheOf[string, entity.X](m)` 创建 `KeyCache`；写路径显式 `kc.Delete`
- 缓存键经 `constants.RedisX.Get(id).String()` 构造（带 `NOSQL_PREFIX` 前缀）；TTL 用 `xCache.WithTTL(...)`
- 配置来源 env-first（`.env` 由基础库加载）；`configs/` 存在但为空
- `xMain.Runner(reg, log, opts, goroutine...)` 配合 `[]xOption.Option` 声明式装配
- `xCtxUtil` 自 `major/utility/context` 导入（v1.1.2 保留 major/utility 路径）
- `xValid`（HandleValidationError）自 `major/validator` 导入
- `xError` 自 `common/error` 导入；`xLog` / `xLogGorm` 自 `common/log` 导入（v1.1.2 保留 common/log）
- `xCache` / `xCache.Manager` / `KeyCacheOf` 自 `major/cache` 导入（v1.1.0 起缓存层迁移至此）

## 反模式
- 手改 `docs/` 下生成的 swagger 文件，或手改 `frontend/src/routeTree.gen.ts`
- 在 handler 方法里写业务规则
- 越过 logic 层从 handler 直接调 repository
- 不接入 `startup_config.go` 就引入新配置源
- 盲信 README 默认值；以 startup 代码为真相源
- 直接持有 `*redis.Client` 做实体缓存——应统一走 `xCache.Manager` 泛型接口
- 用 bun/npm/yarn 替代 pnpm 管理前端依赖
- 手写代码文件缺失版权横幅，或擅自改写横幅中的作者名、年份、网址、许可证声明文字

## 独特风格
- 领域 DTO 放在顶层 `api/` 包树，而非 `internal/` 下
- 代码内默认管理员种子当前为 `xiao_lfeng`（见 prepare 模块）
- 前后端同仓库但前端是完全独立的 Vite 子项目（独立 `package.json` / `pnpm-lock.yaml` / tsconfig）
- Entity 文件使用手写 `BeforeCreate` 钩子（非 `xModels.BaseEntity`）；兼容 v1.1.2，可后续改用 base entity 模式

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
# 后端开发
make dev
go run main.go

# 后端质量
go test ./...
go fmt ./...
go vet ./...

# 后端文档
make swag
swag init -g main.go --parseDependency

# 前端开发（frontend/ 目录下）
pnpm install
pnpm dev           # vite --port 3000
pnpm build         # vite build && tsc
pnpm test          # vitest run
pnpm check         # prettier --write . && eslint --fix
```

## 备注
- 依赖的本地服务：PostgreSQL（`DATABASE_*`）与 Redis（`NOSQL_*`）
- 邮件 E2E 测试按需开启：`ENABLE_SMTP_E2E_TEST=true` 并配置 SMTP 环境变量
- `go.mod` 使用本地 `replace` 将 `beacon-sso-sdk` 指向 `/Users/xiaolfeng/ProgramProjects/Cooperate/phalanx/beacon-sso-sdk`
- bamboo-base-go 升级至 v1.1.2（破坏性变更：缓存层迁移至 `major/cache` 的 `xCache.Manager` + 泛型 `KeyCacheOf`；`internal/repository/cache/` 目录已删除）
- beacon-sso-sdk 本地源已适配 v1.1.2：`xCache.Cache` 替换为本地 `Cache` 结构体，`xCtxUtil` 路径迁移
- 前端固定使用 pnpm（`pnpm-lock.yaml` 为真相源），开发服务器监听 3000 端口
- 路由树由 `@tanstack/router-plugin` 在 Vite 编译期自动生成 `frontend/src/routeTree.gen.ts`

## 引用
- [frontend](./frontend/AGENTS.md) — 前端子项目（React 19 + TanStack + shadcn/ui）
- [internal/app/startup](./internal/app/startup/AGENTS.md) — 启动引导与种子钩子
- [internal/app/route](./internal/app/route/AGENTS.md) — 路由装配与中间件链
- [internal/handler](./internal/handler/AGENTS.md) — HTTP bind/validate/respond
- [internal/logic](./internal/logic/AGENTS.md) — 业务编排与事务
- [internal/repository](./internal/repository/AGENTS.md) — DB 访问与缓存失效
