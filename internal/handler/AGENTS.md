# HANDLER 知识库

## 概述
HTTP 处理层，仅负责请求绑定、参数校验和响应格式编排；业务决策必须留在 `internal/logic/`。

## 目录结构
```text
internal/handler/
|- handler.go                  # 泛型构造器 NewHandler[T]，注入全部 logic 依赖
|- auth.go                     # 认证相关 HTTP 端点
|- link.go                     # 友链端点
|- link_group.go               # 友链分组端点
|- link_color.go               # 友链颜色端点
|- sponsor_channel.go          # 赞助渠道端点
|- sponsor_record.go           # 赞助记录端点
|- dashboard.go                # 仪表盘统计端点
|- public.go                   # 公开/健康检查端点
|- info.go                     # 站点信息端点
```

## 导航指南
| 任务 | 位置 | 说明 |
|---|---|---|
| 共享 handler 构造 | `handler.go` | 泛型 `NewHandler[T]` 注入全部 logic 依赖，`IHandler` 接口约束 |
| 认证 HTTP 端点 | `auth.go` | 登录/注册/密码/重置/校验/个人资料端点 |
| 友链及分组端点 | `link.go`、`link_group.go`、`link_color.go` | bind + logic 调用 + 响应格式化，含用户自助子集 |
| 赞助端点 | `sponsor_channel.go`、`sponsor_record.go` | admin/public 赞助操作 |
| 仪表盘端点 | `dashboard.go` | 友链统计与最近待审核申请 |
| 公开/信息端点 | `public.go`、`info.go` | 健康/信息读写端点 |

## 约定
- 校验流程：`ShouldBind*` → `xValid.HandleValidationError` → 返回，保证错误格式统一
- logic 错误：`_ = c.Error(err)` 后直接 return，由中间件统一渲染错误载荷，避免 handler 自行拼装错误体
- 成功响应优先使用 `xResult.SuccessHasData`，消息文案保持中文以贴合现有 API 风格
- handler 不管理数据库事务、不操作缓存，事务边界归 logic 层
- Swagger 注释必须与真实路由、响应行为一致；行为变更需同步更新注释与路由映射

## 反模式
- 在 handler 中嵌入业务规则——会让逻辑分散、难以测试
- 越过 logic 层直接调用 repository——破坏分层与事务边界
- 自行返回临时 JSON 结构而不使用 `xResult` 约定——前端契约不一致
- 跳过校验助手、手动格式化 bind 错误——错误响应格式漂移

## 调试路径
1. 请求 4xx/校验失败：检查 handler 中 `ShouldBind*` 与 `xValid.HandleValidationError` 调用顺序
2. 响应结构异常：确认是否误用了自定义 JSON 而非 `xResult` 约定
3. 接口文档与实现不一致：核对 handler 顶部 Swagger 注释与 `internal/app/route/` 中路由注册
4. 依赖注入缺失：回到 `handler.go` 的 `NewHandler[T]`，确认所需 logic 已在 `service` 结构体装配
