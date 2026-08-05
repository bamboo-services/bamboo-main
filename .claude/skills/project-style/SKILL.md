---
name: project-style
description: >-
  统一 bamboo-base-go-template 多层架构写作风格，明确 handler、logic、repository
  职责边界与协作规则。当用户说"规范代码风格"、"分层怎么写"、"这个放哪层"、
  "handler 该做什么"、"logic 和 repository 区别"、"帮我检查分层"、"这个逻辑放哪"
  或任何涉及多层架构职责划分的场景时，都应使用此技能。即使用户只是问"这段代码
  写对了吗"或"应该放在哪个文件"，也应当触发。
---

# Project Writing Style

规范本项目在多层架构下的代码写作风格与职责隔离。目标是：

1. 每层只做自己该做的事。
2. 层与层之间只通过稳定接口协作。
3. 避免"顺手做掉"导致的跨层污染。

分层隔离的价值在于：当 handler 不关心存储细节、repository 不承载业务逻辑时，每一层可以独立修改和测试。如果 handler 直接拼 SQL，改数据库就要改 handler；如果 repository 做业务判断，换缓存策略就要改业务逻辑。职责越清晰，变更影响面越小。

---

## 适用范围

- `handler`（接口/传输层）
- `logic`（业务编排层）
- `repository`（数据访问层）

---

## 架构总览

```text
HTTP / RPC Request
   -> handler  (internal/handler/)
   -> logic    (internal/logic/)
   -> repository (internal/repository/)
   -> DB / Redis
```

### 依赖方向（只能向下）

- `handler` → `logic`
- `logic` → `repository`
- `repository` → `db / rdb`

以下依赖方向是错误的，因为它们绕过了中间层，让上层直接耦合了底层实现：

- `handler` 直接调用 `repository` — handler 绕过业务编排，逻辑无处安放
- `handler` 直接操作 `cache / db` — 存储细节泄漏到传输层
- `logic` 直接解析 HTTP 参数 — logic 耦合了传输协议，无法复用
- `repository` 承载业务决策 — 数据访问层承担了不该承担的判断逻辑

---

## 分层职责

### 1) handler 层（接口适配层）

**负责：**

- 接收请求（HTTP/RPC）并做基础绑定与格式校验。
- 将 transport DTO 转换为 logic 入参。
- 调用 logic 并将结果映射为响应 DTO。
- 统一处理错误到协议状态码/错误码。

**不负责：**

- 业务规则判断（如权限矩阵、状态流转、价格计算）。
- 持久化细节或缓存策略。
- 事务控制。

**写作要点：**

- 函数要短，重点是"转换 + 调用 + 映射"。
- 不在 handler 中出现 repository 类型。
- 错误文案面向接口消费者，不泄漏底层细节。

**本项目约定：**

- 通过 `NewHandler[T]` 泛型函数构造 handler（参考 `internal/handler/handler.go`）。
- 成功响应使用 `xResult.SuccessHasData(ctx, message, data)`。
- 错误通过 `ctx.Error(xErr)` 传递，由中间件统一处理。
- 日志使用 `xLog.WithName(xLog.NamedCONT, "HandlerName")`。

参考 `internal/handler/health.go`：

```go
func (h *HealthHandler) Ping(ctx *gin.Context) {
    h.log.Info(ctx, "Ping - 健康检查")
    status, xErr := h.service.healthLogic.Ping(ctx)
    if xErr != nil {
        _ = ctx.Error(xErr)
        return
    }
    xResult.SuccessHasData(ctx, "pong", status)
}
```

---

### 2) logic 层（业务编排层）

**负责：**

- 承载核心业务用例（use case）。
- 组合多个 repository 完成业务流程。
- 执行业务校验、幂等、事务边界、领域规则。
- 定义稳定的输入输出模型（面向上层/下层）。

**不负责：**

- HTTP 参数提取、响应序列化。
- SQL/Redis 语句细节。
- 框架耦合代码（Gin Context 细节外泄到方法签名）。

**写作要点：**

- 方法名用业务语义（如 `CreateUser`, `BanPlayer`, `RefreshProfile`）。
- 返回错误使用 `*xError.Error`（业务错误 vs 系统错误可区分）。
- 业务流程可读性优先：按"校验 → 执行 → 收敛结果"组织。

**本项目约定：**

- 构造函数 `NewXxxLogic(ctx context.Context)`，内部通过 `xCtxUtil.MustGetDB(ctx)` 和 `xCtxUtil.MustGetRDB(ctx)` 获取依赖。
- 返回值统一为 `(result, *xError.Error)` 风格。
- 日志使用 `xLog.WithName(xLog.NamedLOGC, "LogicName")`。
- DTO 引用 `api/<domain>/` 下的结构体。

参考 `internal/logic/health.go`：

```go
func NewHealthLogic(ctx context.Context) *HealthLogic {
    db := xCtxUtil.MustGetDB(ctx)
    rdb := xCtxUtil.MustGetRDB(ctx)
    return &HealthLogic{
        logic: logic{db: db, rdb: rdb, log: xLog.WithName(xLog.NamedLOGC, "HealthLogic")},
        repo: healthRepo{health: repository.NewHealthRepo(db, rdb)},
    }
}

func (l *HealthLogic) Ping(ctx *gin.Context) (*apiHealth.PingResponse, *xError.Error) {
    // ... 校验 → 执行 → 收敛结果
}
```

---

### 3) repository 层（数据访问层）

**负责：**

- 面向实体/聚合提供稳定 CRUD 接口。
- 隔离数据库实现细节（GORM/SQL/索引策略）。
- 提供查询条件、分页、排序等数据访问能力。

**不负责：**

- 业务流程编排。
- 业务策略分支（例如"是否允许创建"）。
- transport DTO 适配。

**写作要点：**

- 输入尽量是明确 query/filter 结构，而非无序参数堆。
- repository 返回领域可用的数据模型，不返回 HTTP 语义。
- 错误包装要保留可观测信息（操作、主键、关键参数）。

**本项目约定：**

- 构造函数 `NewXxxRepo(db *gorm.DB, rdb *redis.Client)`，依赖在构造时注入。
- 返回值统一为 `(data, *xError.Error)` 风格，不返回裸 `error`。
- 错误使用 `xError.NewError(nil, xError.DatabaseError, "描述", false, err)` 包装。
- 日志使用 `xLog.WithName(xLog.NamedREPO, "RepoName")`。
- 使用 `r.db.WithContext(ctx)` 传递请求上下文。

参考 `internal/repository/health.go`：

```go
func NewHealthRepo(db *gorm.DB, rdb *redis.Client) *HealthRepo {
    return &HealthRepo{
        db:  db,
        rdb: rdb,
        log: xLog.WithName(xLog.NamedREPO, "HealthRepo"),
    }
}

func (r *HealthRepo) DatabaseReady(ctx context.Context) (bool, *xError.Error) {
    sqlDB, err := r.db.WithContext(ctx).DB()
    if err != nil {
        return false, xError.NewError(nil, xError.DatabaseError, "获取数据库连接失败", false, err)
    }
    // ...
    return true, nil
}
```

---

## 越界禁止清单

以下规则的意义是保护层间隔离——违反任何一条都会导致变更影响面扩大：

1. handler 只依赖 logic，不直接依赖 repository。handler 直连 repository 意味着跳过了业务校验和事务编排。
2. logic 通过 repository 接触数据源，不拼接 SQL/Redis 命令。logic 拼 SQL 会让换库或加缓存时改动面扩散到业务层。
3. repository 不承载业务分支，不做"是否允许"的业务判断。数据访问层做业务判断会让同一业务规则散落多处。
4. 下层不返回上层协议对象。repository 返回 HTTP 状态码会让数据访问层耦合传输协议。
5. 不跨层复用"顺手函数"破坏边界。handler 调 util 直连 DB 会让存储细节绕过 repository 的隔离。

---

## 分层目录结构

```text
bamboo-base-go-template/
├── api/                        # 请求/响应 DTO（按业务域分包）
│   └── health/
│       └── health.go           # PingResponse 等 DTO
├── internal/
│   ├── app/
│   │   ├── route/              # 路由注册与中间件绑定
│   │   └── startup/            # 基础设施初始化与种子数据
│   ├── constant/               # 共享业务常量（基因编号等）
│   ├── entity/                 # GORM 实体
│   │   ├── user.go
│   │   └── role.go
│   ├── handler/                # HTTP 处理器（薄控制器层）
│   │   ├── handler.go          # NewHandler[T] 泛型构造
│   │   └── health.go
│   ├── logic/                  # 业务编排层
│   │   ├── logic.go            # 基础 logic 结构体
│   │   └── health.go
│   └── repository/             # 数据访问层
│       └── health.go
```

说明：

- `api/<domain>/`：请求/响应 DTO，按业务域分包。
- `handler`：协议适配与响应格式。
- `logic`：业务用例。
- `repository`：持久化与查询。
- `entity`：GORM 实体，嵌入 `xModels.BaseEntity` 并实现 `GetGene()`。
- `constant`：业务常量（如基因编号），统一收口，不散落在各层。

---

## 示例：正确分层

```go
// handler (internal/handler/user.go)
func (h *UserHandler) Create(ctx *gin.Context) {
    var req dto.CreateUserRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        _ = ctx.Error(xError.NewError(nil, xError.ParamError, "请求参数错误", false, err))
        return
    }
    out, xErr := h.service.userLogic.CreateUser(ctx, req.ToLogicInput())
    if xErr != nil {
        _ = ctx.Error(xErr)
        return
    }
    xResult.SuccessHasData(ctx, "创建成功", out)
}

// logic (internal/logic/user.go)
func (l *UserLogic) CreateUser(ctx *gin.Context, in CreateUserInput) (*dto.CreateUserResponse, *xError.Error) {
    if exists, xErr := l.repo.user.ExistsByUsername(ctx.Request.Context(), in.Username); xErr != nil {
        return nil, xErr
    } else if exists {
        return nil, xError.NewError(nil, xError.BusinessError, "用户名已存在", false, nil)
    }
    user, xErr := l.repo.user.Create(ctx.Request.Context(), in.ToEntity())
    if xErr != nil {
        return nil, xErr
    }
    return ToCreateUserResponse(user), nil
}

// repository (internal/repository/user.go)
func (r *UserRepo) Create(ctx context.Context, user *entity.User) (*entity.User, *xError.Error) {
    if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
        return nil, xError.NewError(nil, xError.DatabaseError, "创建用户失败", false, err)
    }
    return user, nil
}
```

---

## 示例：错误分层（反例）

```go
// 错误：handler 直接依赖 repository 并做业务判断
func (h *UserHandler) Ban(ctx *gin.Context) {
    user := h.userRepo.GetByID(ctx, id)  // 越界：handler -> repository，跳过了业务编排
    if user.Role != "ADMIN" {            // 越界：业务规则放在 handler，logic 层被架空
        ctx.JSON(403, ...)
        return
    }
}
```

这段代码的问题：handler 绕过 logic 直接操作 repository，导致业务规则散落在传输层。当权限判断逻辑变更时，需要改 handler 而不是 logic——违反了"业务变更集中在 logic 层"的分层目标。

---

## 评审检查清单

- [ ] handler 仅做协议适配，无业务规则分支。
- [ ] logic 完整表达业务用例，不含 transport/db 框架泄漏。
- [ ] repository 仅做数据访问，不承载业务决策。
- [ ] 依赖方向单向向下，无跨层旁路调用。
- [ ] handler 使用 `xResult.SuccessHasData` 返回成功，`ctx.Error(xErr)` 传递错误。
- [ ] repository 返回 `(data, *xError.Error)`，不返回裸 `error`。
- [ ] 日志按层命名：`NamedCONT` / `NamedLOGC` / `NamedREPO`。

---

## 何时需要 AskUserQuestion

当以下信息不清晰时询问：

1. 该需求属于"新增用例"还是"扩展现有用例"？
2. 事务边界放在 logic 还是 repository（项目约定在 logic）？
3. 错误语义是否有统一业务错误码体系？

---

## 一句话原则

> handler 负责"说人话（协议）"，logic 负责"做决策（业务）"，repository 负责"拿数据（存储）"。
