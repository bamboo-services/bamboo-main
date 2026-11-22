# STARTUP 知识库

## 概述
启动引导层，基于 `bamboo-base-go` `application/v2` 的三态生命周期内核装配应用核心工具插件与启动期迁移/配置能力：包含受管日志与雪花配置构造、业务配置工具（`ConfigTool`）与数据库迁移/种子工具（`MigrationTool`）。

## 目录结构
```text
internal/app/startup/
|- startup.go                  # BuildLoggerConfig / BuildGormConfig / ConfigureGinMode / ParseSnowflakeIDs 环境变量装配辅助
|- config_tool.go              # ConfigTool（ModeTool：EMAIL_ADMIN_EMAIL 业务配置）
|- migration_tool.go           # MigrationTool（ModeTool：AutoMigrate + prepare 种子）
`- prepare/                    # 种子数据准备（由 MigrationTool 在 Start 阶段按序调用）
   |- prepare_default_user.go  #   默认管理员种子
   `- prepare_default_info.go  #   默认站点信息种子
```

## 导航指南
| 任务 | 位置 | 说明 |
|---|---|---|
| 日志/雪花参数读取 | `startup.go` | 读取 `LOG_*`、`APP_DEBUG`、`SNOWFLAKE_*` 装配内置插件配置 |
| 业务配置 Tool | `config_tool.go` | `NewConfigTool()` 暴露 `constants.ContextCustomConfig` |
| 迁移与种子 Tool | `migration_tool.go` | `NewMigrationTool(models...)` 依赖 DB + Snowflake + Logger，暴露 `ContextDatabaseMigrated` |
| 种子数据准备 | `prepare/` | `prepare.DefaultUser` / `prepare.DefaultInfo`（幂等，携带含雪花节点的 `startCtx`） |

## 约定
- 所有插件构造均为纯构造（无 I/O/网络副作用），实际初始化推迟至 `Start(ctx)` 阶段
- `MigrationTool` 必须声明 `Requires: [DatabaseKey, SnowflakeNodeKey, LoggerKey]`，保证 `BaseEntity.BeforeCreate` 可读取雪花节点
- 邮件 SMTP 连接与内嵌模板加载由 `main.go` 中 `email.New(email.FromEnv(), email.WithTemplateFS(...))` 直接托管
- 友链截图常驻服务已抽离为 `screenshot.NewWorker()` 与 `screenshot.NewRefreshComponent()`，不再由 startup 节点创建

## 反模式
- 在构造函数或 `main.go` 注册阶段直接建立 DB/Redis 连接或写磁盘
- 在 `Start` 阶段读取未在 `Requires` 中声明的依赖键（会被内核白名单拦截）
- 绕过 `MigrationTool` 在无雪花节点的裸 `context.Background()` 上插入 `BaseEntity`

## 调试路径
1. 启动失败：检查 `application.App` 拓扑错误与具体插件 `Start(ctx)` 返回的 `*apperr.Error`
2. 表结构未创建：核对 `main.go` 传入 `startup.NewMigrationTool(...)` 的实体列表
3. 种子数据缺失：确认 `MigrationTool` 正常执行且未因 `SnowflakeNodeKey` 缺失中断
