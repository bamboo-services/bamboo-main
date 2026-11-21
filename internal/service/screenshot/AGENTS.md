# SCREENSHOT 知识库

## 概述
友链站点截图任务子系统，负责用无头浏览器（go-rod）自动截取友链站点首页，并持久化到 `SCREENSHOT_DIR`；通过内存 FIFO 队列 + 单 worker 串行消费，随 Runner 生命周期启停。

## 目录结构
```text
internal/service/screenshot/
|- config.go                  # SCREENSHOT_* 环境变量加载（env-first 配置）
|- capture.go                 # rod 无头浏览器截图（CaptureFunc 抽象 + rodCapturer 实现）
|- manager.go                 # 任务管理器：FIFO 队列 + inflight 去重 + 单 worker
`- plugin.go                  # NewWorker() (ModeWorker) 与 NewRefreshComponent() (ModeComponent)
```

## 导航指南
| 任务 | 位置 | 说明 |
|---|---|---|
| 配置加载 | `config.go` | `LoadConfig()` 读 `SCREENSHOT_ENABLED/DIR/CHROME_PATH/CDP_URL/TIMEOUT` |
| 截图实现 | `capture.go` | `NewRodCapture(cfg)` 返回 `CaptureFunc`，内部复用单浏览器实例 |
| 任务入队 | `manager.go` | `Enqueue(id)` / `EnqueueAll(ctx)`（全量查询已通过友链） |
| worker 常驻消费 | `manager.go` | `Run(ctx)` 串行处理队列，ctx 取消或 Close 时退出（禁用态保持待机防看门狗早退） |
| 生命周期插件封装 | `plugin.go` | `NewWorker()` 暴露 `ContextScreenshotManager`，`NewRefreshComponent()` 供 `cron/v2` 激发 |
| 上下文获取 | `manager.go` | `GetManager(ctx)` 通过 `application.Get` 从容器取 manager 实例 |

## 约定
- 截图依赖抽象为 `CaptureFunc` 接口，测试可注入替身；生产默认 `NewRodCapture`
- 队列去重经 `inflight map`，避免同一友链重复排队；进程重启丢失的任务由每日全量入队（`cron/v2` 激发 `RefreshComponent`）兜底
- 截图失败保留旧截图，等待下次自然重试；文件写入采用「临时文件 + rename」原子替换
- 配置 env-first（`SCREENSHOT_*`），worker 与定时任务全部由 `main.go` `application.App` 托管生命周期
- 存储目录经 `/screenshots/*` 静态路由对外访问（route.go 注册）

## 反模式
- 在 logic/handler 中直接起 rod 截图——应统一走 Manager 队列
- 绕过 `CaptureFunc` 抽象、直接依赖 rodCapturer——破坏可测试性
- 直接写最终文件而不经临时文件 + rename——可能产生半写文件被读取
- 在 Manager 外启动不受 `application.App` 托管的协程——生命周期失控
- `SCREENSHOT_ENABLED=false` 时直接从 `Worker.Run` 返回——会被 v2 看门狗判定为崩溃并触发整机停机

## 调试路径
1. 截图缺失：核对 `SCREENSHOT_ENABLED` 与 `SCREENSHOT_CHROME_PATH`/`SCREENSHOT_CDP_URL` 配置
2. 浏览器启动失败：容器内 root 运行检查 `NoSandbox` 逻辑；外部 CDP 模式检查连接地址
3. 队列不消费：确认 `main.go` 已注册 `screenshot.NewWorker()` 且迁移工具前置完成
4. 截图 URL 路径：文件名为 `<SnowflakeID>.png`，数据库存 `/screenshots/<id>.png`
