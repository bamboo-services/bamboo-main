# Repository Guidelines

## 项目结构与模块组织
- `main.go` 为 Gin 入口，组合 `internal/app/startup` 完成配置、数据库、Redis 与默认数据初始化，并通过 Runner 挂载常驻 worker。
- `internal/app/route` 注册路由；`internal/handler` 处理 HTTP 请求；`internal/logic` 封装业务流程；`internal/repository` 访问数据库/Redis（含 cache 生命周期管理）。
- `internal/models/base` 管理配置结构；数据库实体放在 `internal/entity`。
- 配置采用 `.env` 环境变量方案；`docs/` 生成的 Swagger 规范；`logs/` 为运行日志输出目录。仓库根下的 `bamboo-main` 是已编译二进制，不应提交新版。

## 开发、构建与运行
- 需要 Go 1.24+，并确保 `go.mod` 的 `replace github.com/bamboo-services/bamboo-base-go => ../bamboo-base` 指向可用路径。
- 安装依赖：`go mod tidy`。
- 构建：`go build ./...`（生成二进制 `bamboo-main`）。
- 运行：`go run main.go`（读取环境变量；可复制 `.env.example` 为 `.env`）。
- 测试：`go test ./...`；表驱动写法优先，必要时使用 `httptest` 针对 Gin 路由。
- 若修改 API 注解，需要 `swag init -g main.go -o docs`（需预装 `github.com/swaggo/swag/cmd/swag@latest`）。

## 编码风格与命名规范
- 提交前执行 `gofmt`/`goimports`（建议 `go fmt ./...`）。缩进使用制表符，保持 Go 默认风格。
- 包名全小写、无下划线；导出符号用 PascalCase，未导出使用 camelCase，常量使用驼峰或全大写并靠近使用处。
- 处理链条：`handler` 只负责请求校验与响应包装，业务放入 `logic`，资源访问放入 `service`。
- DTO/Request 结构体放入 `internal/model/dto|request`，数据库实体放入 `entity`，接口路径保持 RESTful，日志统一使用注入的 Sugar logger 并携带上下文。

## 测试指南
- 单元测试文件命名为 `*_test.go`，函数命名 `TestXxx`，对子场景使用表驱动 case。
- 涉及 Gin 路由的测试可通过 `httptest.NewRecorder` + `router.Init` 构建。对数据库/Redis 的用例请隔离测试库或使用 mock，避免污染本地数据。
- 新增功能优先提供基础覆盖（建议≥70% 对新代码），并在 PR 中说明测试范围与结果。

## 提交与 Pull Request
- 遵循 Conventional Commits，历史示例：`feat: 实现友链分组管理功能并重构网络工具包`、`refactor: 重构认证中间件并优化用户上下文管理`、`chore(scope): ...`。可使用中文主题，必要时添加 scope。
- 提交前自检：`go fmt ./...`、`go test ./...`、如有改动更新 `docs/swagger.*`。
- PR 描述需包含：变更摘要、关联 Issue/需求编号、主要测试步骤与结果、接口变更示例（如 curl/响应示例），以及潜在风险或回滚方式。涉及配置改动请标注新增键及默认值。

## 配置与安全提示
- 默认配置来源为环境变量（本地可使用 `.env`）；生产环境请改用私有配置文件或环境变量，勿提交真实凭证。调试标志 `XLF_DEBUG` 开启详细日志，发布前保持关闭。
- PostgreSQL/Redis 连接依赖 `database.*` 与 `nosql.*` 字段，请确保端口可达；邮件配置为空时相关功能将被跳过。
