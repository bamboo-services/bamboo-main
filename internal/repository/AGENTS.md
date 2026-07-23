# REPOSITORY 知识库

## 概述
数据访问层，封装 GORM 查询与缓存失效策略；写路径通常显式失效对应缓存键。
缓存经 bamboo-base-go v1.1.0 的 `xCache.Manager` 统一管理，repository 通过泛型工厂
`xCache.KeyCacheOf[string, entity.X]` 创建 `KeyCache` 实例完成 cache-aside 读写。

## 目录结构
```text
internal/repository/
|- link.go                     # 友链持久化
|- link_group.go               # 分组持久化
|- link_color.go               # 颜色持久化
|- sponsor_channel.go          # 赞助渠道持久化
|- sponsor_record.go           # 赞助记录持久化
|- system_user.go              # 用户认证查询
|- system.go                   # 系统配置存储
```

## 导航指南
| 任务 | 位置 | 说明 |
|---|---|---|
| 友链持久化 | `link.go` | CRUD、列表/过滤、关联清理、缓存失效 |
| 分组/颜色持久化 | `link_group.go`、`link_color.go` | 状态/排序/列表/删除 |
| 赞助持久化 | `sponsor_channel.go`、`sponsor_record.go` | 渠道/记录查询 |
| 用户/系统持久化 | `system_user.go`、`system.go` | 用户认证查询与系统配置存储 |
| 缓存接入 | 各 repo 构造器 | `xCache.KeyCacheOf[string, entity.X](m)` 创建实例存为 `kc` 字段 |

## 约定
- 构造器统一为 `New*Repo(db *gorm.DB, m *xCache.Manager)`，内部用 `xCache.KeyCacheOf[string, entity.X](m)` 创建 `kc`
- 方法接受可选事务 `tx *gorm.DB`，内部经 `pickDB(tx)` 解析目标 DB，保证事务一致性
- 查询上下文路径与调用方一致，使用请求 context
- 缓存策略为 cache-aside：读穿透 + 写/删时显式失效
- 缓存键经 `constants.RedisX.Get(id).String()` 构造（带 `NOSQL_PREFIX` 前缀），`K=string`、`V=entity.X`
- 写入 TTL 通过 `xCache.WithTTL(...)` 显式传入（link/group/color/user 15min，sponsor 10min）
- not-found 行为必须显式（返回 `found bool` 或检查 `RowsAffected`），不依赖 nil 推断

## 反模式
- 在 repository 方法中分支业务规则——破坏分层
- 可变操作后跳过缓存失效——导致脏读
- 绕过 `pickDB(tx)` 直接用 `db`——破坏事务一致性
- 返回框架级响应载荷——数据层不应感知 HTTP 响应结构
- 直接持有 `*redis.Client` 做实体缓存——应统一走 `xCache.Manager` 泛型接口

## 调试路径
1. 数据不一致/脏读：检查写路径是否调用 `kc.Delete` 失效，Redis key 是否经 `RedisKey.Get()` 加前缀
2. 事务未生效：排查 logic 是否传 `tx`，repository 是否走 `pickDB(tx)`
3. not-found 误判：确认返回值约定（`found bool` vs `RowsAffected`）
4. 缓存键错误：对照 `pkg/constants/cache.go` 的 `RedisKey` 常量与参数占位符
5. 后端切换：`m.Type()` 区分 Redis/Memory，工厂方法自动分发，无需 repo 感知
