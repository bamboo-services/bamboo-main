---
name: swagger-writer
description: >-
  为 Go handler 编写和补全 godoc-swagger 注释，严格对齐本项目 swag 风格与字段顺序。
  当用户提到"写 swagger 注释"、"补 godoc"、"修复 swag 注解"、"给接口补文档"、
  "生成 API 文档"、"handler 注释不全"或任何涉及 swaggo/swag 注释编写的场景时，
  都应使用此技能。即使用户只是说"给这个接口加个文档"或"注释写一下"，也应触发。
---

# Swagger Writer

为 `internal/handler/*.go` 中的接口方法编写或修正 Godoc + Swag 注释，目标是直接可用于 `swag init` 生成文档。

注释是 `swag init` 解析的唯一来源——如果注释和代码不一致，生成的 API 文档就会误导调用方。所以这个技能的核心原则是：**先读代码定事实，再写注释补表达**。

---

## 工作流

1. **读取目标 handler 文件与对应函数**，理解接口实际行为。
2. 若路由/请求模型不明确，继续读取 `route`、`api/<domain>/`（DTO）、`entity`、`logic` 相关文件。
3. 确定接口事实信息：
   - 路径（`@Router`）— 从 `route_*.go` 注册代码确认
   - 方法（GET/POST/PUT/DELETE...）— 同上
   - 入参位置（path/query/body/header）— 从 handler 的 `ShouldBindJSON` / `Param` / `Query` 调用确认
   - 成功响应数据类型（`xBase.BaseResponse{data=...}`）— 从 `xResult.SuccessHasData` 调用的实参确认
   - 失败状态码范围 — 从 `ctx.Error(xErr)` 传递的错误类型推断
4. 按固定顺序生成注释块并写回函数前。
5. 自检：注释内容与函数行为一致，不捏造字段。

---

## 注释顺序规范

按以下顺序输出，保持与项目内已有 swagger 注释一致：

```go
// FuncName 中文一句话说明
//
// @Summary     [玩家/管理/超管] 接口名
// @Description 说明接口做什么，输入输出是什么
// @Tags        模块接口
// @Accept      json
// @Produce     json
// @Param       ...
// @Success     200   {object}  xBase.BaseResponse{data=entity.XXX}  "成功"
// @Failure     400   {object}  xBase.BaseResponse                    "请求参数错误"
// @Failure     401   {object}  xBase.BaseResponse                    "未授权"
// @Failure     403   {object}  xBase.BaseResponse                    "无权限"
// @Failure     404   {object}  xBase.BaseResponse                    "资源不存在"
// @Router      /api/v1/xxx [GET]
```

几点说明：

- 无请求参数时省略 `@Param`，不要写空参数占位。
- `@Failure` 只保留真实可能出现的状态码。写不存在的错误码会让调用方误判接口的失败语义。
- `@Summary` 推荐使用 `[玩家/管理/超管] 动作` 结构，例如 `[玩家] 用户信息`，这样在 Swagger UI 中按角色分组更清晰。

---

## @Param 写法

### Path 参数

```go
// @Param id path int true "用户ID"
```

### Query 参数

```go
// @Param page query int false "页码"
// @Param size query int false "每页数量"
```

### Body 参数

```go
// @Param request body dto.CreateUserRequest true "创建用户请求"
```

Body 类型引用 `api/<domain>/` 下的 DTO 结构体。如果 handler 使用 `ShouldBindJSON` 绑定，就从绑定目标确认类型。

### Header 参数（按需）

```go
// @Param Authorization header string true "Bearer Access Token"
```

---

## 响应模型

本项目统一使用 `xBase.BaseResponse` 作为响应外壳（定义在 `bamboo-base-go/common/base_response.go`，包名 `xBase`）。

成功响应通过 `data=` 语法标注内部数据类型：

```go
// @Success 200 {object} xBase.BaseResponse{data=entity.User} "成功"
```

列表响应按真实结构填写 data：

- `data=[]entity.User` — 返回实体切片
- `data=dto.UserListResponse` — 返回自定义列表 DTO

注释中的类型与函数实际通过 `xResult.SuccessHasData(ctx, msg, data)` 传递的 `data` 参数类型保持一致，这样 `swag init` 生成的文档才能准确反映接口行为。

---

## 文案风格

- 注释文案使用中文，简洁、可读。
- `@Description` 说明"依据什么入参，返回什么结果"，让调用方快速理解接口用途。
- `@Tags` 统一使用"中文模块 + 接口"，例如：`用户接口`、`认证接口`，与文件内其他注释保持一致。
- 不混用中英标点格式。

---

## 质量检查清单

写完注释后逐项自检：

- [ ] 注释块位于函数定义正上方。
- [ ] `@Router` 的路径与方法和 `route_*.go` 中实际注册一致。
- [ ] `@Param` 名称、位置、必填状态与 handler 中的绑定结构一致。
- [ ] `@Success` 的 `data=` 类型与 `xResult.SuccessHasData` 实参一致。
- [ ] 未改动函数逻辑、返回流程、错误处理——只动注释。

---

## 何时询问用户

在以下信息无法从代码推断时，使用 `AskUserQuestion`：

1. 同一函数被多个路由复用，无法确定主路由。
2. 返回数据模型存在多个候选（entity/dto 均可能）。
3. 业务要求的失败码文案有团队约定但代码中未体现。

优先先读代码再问。能从代码推断的场景下直接判断，不浪费用户时间。

---

## 示例

```go
// UserCurrent 获取用户的信息
//
// @Summary     [玩家] 用户信息
// @Description 根据 AT 获取用户信息，获取到本程序的用户信息
// @Tags        用户接口
// @Accept      json
// @Produce     json
// @Success     200   {object}  xBase.BaseResponse{data=entity.User}  "登录成功"
// @Failure     400   {object}  xBase.BaseResponse                    "请求体格式不正确"
// @Failure     401   {object}  xBase.BaseResponse                    "用户名或密码错误"
// @Failure     403   {object}  xBase.BaseResponse                    "用户已禁用或账户已锁定"
// @Failure     404   {object}  xBase.BaseResponse                    "用户不存在"
// @Router      /api/v1/user/info [GET]
func (h *UserHandler) UserCurrent(ctx *gin.Context) {}
```

---

一句话准则：先读代码定事实，再写注释补表达；注释必须准确，不允许"想当然"。
