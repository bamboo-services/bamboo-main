# STARTUP 知识库

## 概述
启动引导层，负责在服务对外前装配业务依赖；自 v1.0.4 起，DB/Cache 由 `main.go` 的 xOption 声明式配置接管，本层注册框架邮件客户端、邮件业务配置与 SSO 节点。

## 目录结构
```text
internal/app/startup/
|- startup.go                  # Init() 返回根 context 与自定义注册节点
|- startup_config.go          # 邮件业务配置初始化（管理员邮箱）
`- prepare/                    # 种子数据准备（xOptionDB.WithPrepare 钩子）
   |- prepare.go               #   DefaultData 入口
   |- prepare_default_user.go  #   默认管理员种子
   `- prepare_default_info.go   #   默认站点信息种子
```

## 导航指南
| 任务 | 位置 | 说明 |
|---|---|---|
| 启动注册 | `startup.go` | `Init()` 追加 Email 客户端 + Email 业务配置 + SSO 节点 |
| 邮件客户端 | `startup.go` | 框架 `xEmail.InitClient` 经 `xCtx.EmailClientKey` 注入，SMTP 配置读 `EMAIL_*` |
| 邮件业务配置 | `startup_config.go` | 仅管理员邮箱（`EMAIL_ADMIN_EMAIL`），SMTP 连接配置由框架接管 |
| 种子数据准备 | `prepare/` | `prepare.DefaultData` 作为 xOptionDB.PrepareFunc |

## 约定
- 启动节点保持副作用化、确定性，避免隐式顺序依赖
- DB/Cache 配置走 `main.go` 的 xOption，不在本层注册
- 邮件 SMTP 连接配置走框架 xEmail 插件（`EMAIL_*` 环境变量），本层仅注册客户端节点与业务配置
- prepare 任务作为 `xOptionDB.WithPrepare` 钩子，在 AutoMigrate 之后执行
- 依赖初始化错误应 fail-fast（邮件客户端/SSO 客户端）
- 配置默认值仅用于本地/开发安全兜底

## 反模式
- 在此重新引入 DB/Redis init 节点——应使用 `main.go` 的 xOption
- 在 `startup_config.go` 之外新增配置源
- 自研 SMTP 连接/发送逻辑——应使用框架 xEmail 插件
- 在 startup/prepare 中写领域业务规则

## 调试路径
1. 启动失败：检查 `Init()` 节点是否 fail-fast，邮件客户端或 SSO 客户端是否就绪
2. 表结构未创建：回到 `main.go` 的 `xOptionDB.WithAutoMigrate` 实体列表
3. 种子数据缺失：确认 `prepare.DefaultData` 是否作为 `WithPrepare` 钩子注册
4. 邮件客户端缺失：核对 `EMAIL_HOST/USER/PASS/FROM` 是否配置，`xCtx.EmailClientKey` 节点是否注册
