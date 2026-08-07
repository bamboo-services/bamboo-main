# ROUTE 知识库

## 概述
HTTP 集中装配层，负责 Gin 分组、中间件链与 handler 注册；路由按领域（`auth`/`public`/`link`/`user`/`info`/`sponsor`/`admin`）拆分文件，并以内嵌前端产物实现单二进制部署。

## 目录结构
```text
internal/app/route/
|- route.go                    # 全局中间件、/api/v1 根分组、/metrics 与 /screenshots 静态路由
|- route.go NoRoute 分流        # frontend.go 的 noRouteHandler（SPA 单二进制部署）
|- route_auth.go               # 认证端点（公开子组 + 鉴权子组）
|- route_admin.go              # admin 端点 + OAuth 校验 + 本地鉴权 + 角色中间件
|- route_user.go               # 用户自助端点（links/sponsors/profile，AuthMiddleware 保护）
|- route_public.go             # 公开健康/ping 路由（无鉴权）
|- route_link.go              # 友链及公开友链路由
|- route_sponsor.go            # 赞助相关路由
|- route_info.go               # 站点信息路由
|- route_swagger.go            # Swagger 挂载（仅 debug 模式）
|- frontend.go                 # noRouteHandler：/api 与 /swagger 走 JSON 404，其余 SPA fallback
```

## 导航指南
| 任务 | 位置 | 说明 |
|---|---|---|
| 全局中间件与 `/api/v1` 根 | `route.go` | response/cors/options/metrics 中间件顺序在此定义 |
| 认证端点 | `route_auth.go` | 拆分公开认证与鉴权子组 |
| admin 端点 + 鉴权门 | `route_admin.go` | OAuth 校验 + 本地鉴权 + 角色中间件 |
| 用户自助端点 | `route_user.go` | 友链/赞助/资料自助管理，AuthMiddleware 保护 |
| 公开健康/ping 路由 | `route_public.go` | 无鉴权中间件 |
| 友链/公开友链路由 | `route_link.go` | 公开友链列表端点 |
| Prometheus 埋点与端点 | `route.go` | `metrics.Middleware()` 全局埋点 + `GET /metrics`（promhttp 直写） |
| 截图静态资源 | `route.go` | `GET /screenshots/*` 指向 `SCREENSHOT_DIR` |
| SPA 单二进制托管 | `frontend.go` | NoRoute 阶段按路径前缀分流，未命中静态资源回退 index.html |
| Swagger 挂载 | `route_swagger.go` | 仅 debug 模式注册 |

## 约定
- handler 实例化统一用 `handler.NewHandler[T](r.context, "Name")`
- 路由文件只做装配，不含请求解析或业务分支
- 受保护路由必须显式挂在 OAuth + `AuthMiddleware` + 角色中间件链之下
- 路径命名稳定在 `/api/v1/...`，按领域分组
- `frontend.go` 的 noRouteHandler 保持 `/api/*`、`/swagger/*` 走 JSON 404，避免污染 API 错误语义
- Swagger 路由仅在 debug 环境开启

## 反模式
- 在路由方法中直接写业务逻辑
- admin 路由未过 OAuth + `AuthMiddleware` + 角色校验
- 在 noRouteHandler 中把 `/api/*` 兜底到 SPA index.html——API 路径必须保持 JSON 404
- 生产模式默认注册 Swagger 路由
- 各路由文件构造 handler 的方式不一致

## 调试路径
1. 404/405：回到 `route.go` 确认 `NoRoute`/`NoMethod` 与领域子路由是否注册
2. 鉴权漏放：核对 `route_admin.go`/`route_user.go` 的中间件链顺序（OAuth → 本地鉴权 → 角色）
3. CORS 异常：检查 `xMiddle.ResponseMiddleware`/`ReleaseAllCors`/`AllowOptionRequest` 顺序
4. 前端页面 404：确认 `frontend.go` 的 noRouteHandler 已注册，内嵌产物 `resources/frontend/embed.go` 是否刷新
5. Swagger 不显示：确认 `xEnv.Debug` 为真且 `swaggerRegister` 被调用
