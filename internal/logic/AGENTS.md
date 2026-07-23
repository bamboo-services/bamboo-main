# LOGIC 知识库

## 概述
业务编排层，承担领域规则、事务边界与异步副作用；项目绝大多数非平凡行为均在此实现。

## 目录结构
```text
internal/logic/
|- logic.go                    # 基础 logic 结构体（db/rdb/log）
|- auth.go                     # 认证/会话/密码流程
|- auth_oauth.go               # OAuth 同步与登录逻辑
|- link.go                     # 友链生命周期与公开列表
|- link_group.go               # 友链分组管理
|- link_color.go               # 友链颜色管理
|- sponsor_channel.go          # 赞助渠道编排
|- sponsor_record.go           # 赞助记录编排
|- public.go                   # 公开/健康检查逻辑
|- info.go                     # 站点信息读写逻辑
|- mail.go                     # 异步邮件触发入口
|- helper/                     # 跨域 helper（session/mail）
|  |- session.go
|  |- mail.go
```

## 导航指南
| 任务 | 位置 | 说明 |
|---|---|---|
| 认证/会话/密码流程 | `auth.go`、`auth_oauth.go` | OAuth 同步、token/session 行为、重置/邮件流程 |
| 友链生命周期与公开列表 | `link.go` | 增删改查/状态/失效 + 通知触发 |
| 友链分类管理 | `link_group.go`、`link_color.go` | 排序/状态/删除约束 |
| 赞助领域逻辑 | `sponsor_channel.go`、`sponsor_record.go` | 渠道/记录编排与数据塑形 |
| 公开/系统信息 | `public.go`、`info.go` | 健康/信息读写行为 |
| 跨域 helper | `helper/` | session/mail 复用逻辑 |

## 约定
- 构造器统一为 `New*Logic(ctx context.Context)`，在内部初始化 repo/helper
- 领域规则与事务边界留在 logic，不下沉到 handler/repository
- 使用项目 typed constants 表达 status/role/context 语义，避免裸字符串
- 异步触发（邮件/通知）为非阻塞副作用，失败不应回滚主流程
- 错误返回统一使用 `*xError.Error`，交由上游错误中间件渲染

## 反模式
- 把 DB 查询细节从 repository 搬到 logic——破坏数据访问封装
- 在 logic 层写响应格式化逻辑——响应编排归 handler
- 启动临时 goroutine 而无明确失败/隔离策略——易泄漏、难观测
- 向本就密集的核心文件塞入无关新逻辑而不抽取——可读性恶化

## 调试路径
1. 业务规则不符预期：从对应 domain 文件入手，`auth.go` 与 `link.go` 是当前复杂度热点
2. 事务边界异常：检查 logic 中 `tx` 开启点与 repository 的 `pickDB(tx)` 调用链
3. 副作用丢失/延迟：定位 `helper/mail.go` 与 `internal/task/` 的队列投递路径
4. 错误未渲染：确认返回类型为 `*xError.Error` 且 handler 已 `_ = c.Error(err)`
