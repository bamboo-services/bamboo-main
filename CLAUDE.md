# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

Bamboo-Main 是一个基于 Gin + GORM 构建的友情链接管理系统，采用清洁架构设计。

- **技术栈**: Go 1.24.6, Gin, GORM, PostgreSQL, Redis
- **默认端口**: 23333
- **依赖**: 本地依赖 `bamboo-base-go`，路径通过 `replace` 指向 `/Users/xiaolfeng/ProgramProjects/Cooperate/bamboo-service/bamboo-base`

## 常用命令

```bash
# 开发
go run main.go                    # 启动服务
go mod tidy                       # 安装依赖

# 构建与测试
go build -o bamboo-main           # 编译二进制
go test ./...                     # 运行测试
go fmt ./...                      # 代码格式化
go vet ./...                      # 代码检查

# 文档
swag init -g main.go -o docs      # 生成 Swagger 文档 (访问 /swagger/index.html)
```

## 架构设计

### 三层架构 (Handler → Logic → Service)

```
Handler (internal/handler/)      - HTTP 请求处理、参数校验、响应封装
   ↓
Logic (internal/logic/)          - 业务逻辑编排、数据转换、规则校验
   ↓
Service (internal/service/)      - 数据库/Redis 操作、数据访问封装
```

**重要原则**:
- Handler **不得**包含业务逻辑，只负责参数绑定和响应
- Logic **不得**直接操作数据库，必须通过 Service
- Service **只**负责数据访问，不处理业务规则

### 目录结构

```
internal/
├── handler/          # HTTP 处理层
├── logic/            # 业务逻辑层
├── service/          # 数据访问层
├── middleware/       # 中间件 (认证、权限)
├── model/
│   ├── entity/       # 数据库实体 (对应表结构)
│   ├── dto/          # 数据传输对象 (API 响应)
│   └── request/      # 请求参数
└── router/           # 路由配置

pkg/
├── startup/          # 应用启动初始化
├── constants/        # 全局常量 (Redis Key、状态码等)
└── util/             # 工具函数
```

## 核心业务模块

### 认证系统

- **Token 格式**: `cs_` + 64位随机字符串
- **存储方式**: Redis, Key: `bm:auth:token:{token}`, TTL: 24小时
- **中间件**: `middleware.AuthMiddleware()` - 验证 Token 并注入用户上下文
- **权限控制**: `middleware.RequireRole("admin", "moderator")` - 基于角色的访问控制

**上下文使用**:
```go
// 获取当前用户 UUID
userUUID := ctx.GetString(constants.ContextKeyUserUUID)

// 获取当前用户信息 (需先从 service 加载)
user, _ := xCtx.GetCurrentUser(ctx)
```

### 友情链接

- **状态字段**: `0-待审核`, `1-已通过`, `2-已拒绝`
- **失效字段**: `is_failure` (0-正常, 1-失效)
- **外键关联**: `GroupUUID` (分组), `ColorUUID` (颜色)

### 友链分组

- **排序更新**: 支持批量更新排序 (PATCH `/api/v1/admin/groups/sort`)
- **级联影响**: 删除分组时,关联友链的 `group_uuid` 会被置为 NULL

## 数据库设计

### 关键约束

- **主键**: 使用 UUID v7 (时间排序)
- **表前缀**: `bm_` (通过 `database.prefix` 配置)
- **外键行为**:
  - UPDATE: CASCADE
  - DELETE: SET NULL
- **自动时间戳**: `BeforeCreate`/`BeforeUpdate` Hooks

### 迁移

数据库迁移在启动时自动执行 (`pkg/startup/register_database.go`):
```go
xDatabase.DB.AutoMigrate(
    &entity.SystemUser{},
    &entity.LinkFriend{},
    &entity.LinkGroup{},
    &entity.LinkColor{},
    &entity.Log{},
    &entity.SystemInfo{},
)
```

**注意**: 修改实体结构后,需重启服务以应用迁移。

## Redis 命名规范

**前缀**: `bm` (bamboo-main)

```go
// 常量位置: pkg/constants/redis.go
bm:auth:token:{token}      // Token 会话
bm:link:cache:{uuid}       // 链接缓存
bm:group:list:{key}        // 分组列表
bm:email:limit:{email}     // 邮件频率限制
```

## 配置管理

**配置文件**: `configs/config.yaml`

```yaml
xlf:
  debug: true              # true=开发模式(启用Swagger), false=生产模式
  server:
    port: 23333

database:                  # PostgreSQL 配置
  host: localhost
  port: 5432
  user: bamboo_main
  pass: bamboo_main
  name: bamboo_main
  prefix: bm_
  sslmode: disable
  timezone: Asia/Shanghai

nosql:                     # Redis 配置
  host: localhost
  port: 6379
  pass: ""
  database: 0
  prefix: "bm"

email:                     # 邮件配置 (TODO: 未实现)
  smtp_host: smtp.gmail.com
  smtp_port: 587
```

**加载方式**: `pkg/startup/register_config.go` → 解析为 `model.BambooConfig` → 注入上下文

## API 开发指南

### 新增接口流程

1. **定义实体** (如需新表): `internal/model/entity/xxx.go`
2. **定义请求参数**: `internal/model/request/xxx_request.go`
3. **定义响应 DTO**: `internal/model/dto/response/xxx_response.go`
4. **实现 Service**: `internal/service/xxx.go` (数据访问)
5. **实现 Logic**: `internal/logic/xxx.go` (业务逻辑)
6. **实现 Handler**: `internal/handler/xxx.go` (HTTP 处理)
7. **注册路由**: `internal/router/router_xxx.go`
8. **添加 Swagger 注释**: 在 Handler 方法上添加 `// @Summary`, `// @Tags` 等
9. **重新生成文档**: `swag init -g main.go -o docs`

### Handler 编写规范

```go
// AddLink godoc
// @Summary      添加友情链接
// @Description  管理员添加新的友情链接
// @Tags         admin-links
// @Accept       json
// @Produce      json
// @Param        request  body      request.LinkAddRequest  true  "添加友链请求"
// @Success      200      {object}  result.Result
// @Failure      400      {object}  result.Result
// @Router       /api/v1/admin/links [post]
// @Security     Bearer
func (h *LinkHandler) AddLink(ctx *gin.Context) {
    var req request.LinkAddRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        xResult.Error(ctx, xError.ParamError)
        return
    }

    if err := logic.LinkLogic{}.Add(ctx, &req); err != nil {
        xResult.Error(ctx, err)
        return
    }

    xResult.Success(ctx, "添加成功", nil)
}
```

## 开发规范 (来自 AGENTS.md)

### 编码风格

- **格式化**: 使用 `gofmt`/`goimports`,制表符缩进
- **命名**:
  - 包名: 全小写无下划线
  - 导出: PascalCase
  - 未导出: camelCase
- **导入**: 标准库 → 第三方库 → 本地包,分组留空行

### 数据库 Context 使用规范

**重要**：从 `xCtxUtil.GetDB(ctx)` 获取的 db 实例已包含 Snowflake 节点，**不要**再次调用 `WithContext` 覆盖！

❌ **错误**：
```go
db := xCtxUtil.GetDB(ctx)
err := db.WithContext(ctx.Request.Context()).Create(&user).Error  // 会丢失 Snowflake 节点
```

✅ **正确**：
```go
db := xCtxUtil.GetDB(ctx)
err := db.Create(&user).Error  // 直接使用即可
```

**原理**：
- 数据库在初始化时（`pkg/startup/register_database.go:85-88`）已设置包含 Snowflake 节点的 context
- 中间件会将这个 db 实例注入到每个请求的 Gin Context 中
- 使用 `WithContext()` 会替换原有的 context，导致 BeforeCreate Hook 无法获取 Snowflake 节点生成主键 ID
- 直接使用 `xCtxUtil.GetDB(ctx)` 返回的 db 即可，无需额外设置 context

### 提交规范

遵循 Conventional Commits:
```
feat: 实现友链分组管理功能并重构网络工具包
refactor: 重构认证中间件并优化用户上下文管理
fix: 修复友链排序异常问题
chore: 更新依赖版本
```

### 测试规范

- 测试文件: `*_test.go`
- 函数命名: `TestXxx`
- 推荐表驱动测试
- 覆盖率目标: ≥70%

## bamboo-base-go 使用

这是本地共享基础库,提供通用功能:

```go
// 初始化
xInit "github.com/bamboo-services/bamboo-base-go/init"

// 响应封装
xResult "github.com/bamboo-services/bamboo-base-go/result"
xResult.Success(ctx, "成功消息", data)
xResult.Error(ctx, xError.ParamError)

// 错误处理
xError "github.com/bamboo-services/bamboo-base-go/error"
xError.ParamError       // 参数错误
xError.UnauthorizedError // 未授权

// 数据库
xDatabase "github.com/bamboo-services/bamboo-base-go/database"
xDatabase.DB.Model(&entity.LinkFriend{}).Find(&links)

// Redis
xRedis "github.com/bamboo-services/bamboo-base-go/nosql"
xRedis.Redis.Set(ctx, key, value, expiration)

// 上下文工具
xCtx "github.com/bamboo-services/bamboo-base-go/utility/ctx"
user, _ := xCtx.GetCurrentUser(ctx)  // 获取当前登录用户
```

**注意**: 修改 `bamboo-base-go` 后,需在两个项目中同时运行 `go mod tidy`。

## 常见问题

### 1. Swagger 文档未更新

**解决**: 运行 `swag init -g main.go -o docs` 重新生成。

### 2. 数据库连接失败

**检查**:
- PostgreSQL 是否运行在 5432 端口
- `configs/config.yaml` 中的数据库配置是否正确
- 数据库 `bamboo_main` 是否已创建

### 3. Redis 连接失败

**检查**:
- Redis 是否运行在 6379 端口
- `configs/config.yaml` 中的 Redis 配置是否正确

### 4. 编译错误: bamboo-base-go 找不到

**解决**: 检查 `go.mod` 中的 `replace` 指令路径是否正确:
```
replace github.com/bamboo-services/bamboo-base-go v1.0.0-202508212147
    => /Users/xiaolfeng/ProgramProjects/Cooperate/bamboo-service/bamboo-base
```

### 5. 默认管理员账户

- **用户名**: `admin`
- **密码**: `admin123456`
- **邮箱**: `admin@example.com`

**首次运行时自动创建** (`pkg/startup/register_default_user.go`)。

## 待实现功能 (TODO)

1. **友链颜色管理**: Handler 和 Logic 层实现
2. **系统用户管理**: 完整的用户 CRUD 接口
3. **系统日志管理**: 日志查询、清理功能
4. **邮件功能**: 友链申请通知、密码重置邮件
5. **频率限制**: 邮件和申请的限流逻辑

## 技术亮点

- **清洁架构**: Handler → Logic → Service 严格分层
- **UUID v7**: 时间排序的主键设计,兼顾唯一性与索引性能
- **自定义 Token**: `cs_` 前缀 + Redis 会话,24小时过期
- **RBAC 权限**: 基于角色的访问控制 (admin/moderator)
- **统一 Redis Key**: `bm:` 前缀,规范化命名
- **Swagger 自动文档**: 注释生成 API 文档,开发调试便捷
