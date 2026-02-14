# Bamboo Main - 友情链接管理系统

基于 Gin + GORM + PostgreSQL + Redis 的友情链接管理系统，从 GoFrame 架构迁移而来。

## 🎯 项目特性

- **现代化架构**: Gin + GORM + PostgreSQL + Redis
- **清洁架构**: Handler → Service → Logic → Model 分层设计
- **安全认证**: 基于 xUtil.GenerateSecurityKey() 的会话管理
- **Redis 缓存**: 完善的缓存策略和会话存储
- **API 文档**: 自动生成的 Swagger 文档
- **友情链接管理**: 完整的友链申请、审核、分组、颜色管理

## 🛠 技术栈

- **Web 框架**: Gin
- **ORM**: GORM
- **数据库**: PostgreSQL
- **缓存**: Redis
- **认证**: 自定义 Token + Redis 会话
- **文档**: Swagger
- **配置**: `.env` 环境变量

## 📁 项目结构

```
bamboo-main/
├── main.go                    # 应用入口
├── go.mod                     # 依赖管理
├── .env.example               # 环境变量示例
├── internal/
│   ├── handler/              # HTTP 处理层
│   ├── service/              # 服务接口层
│   ├── logic/                # 业务逻辑层
│   ├── models/               # 配置模型与基础结构
│   ├── entity/               # 数据库实体
│   ├── repository/           # 数据访问层（含 cache）
│   ├── middleware/           # 中间件
│   └── app/                  # 路由与启动编排
├── pkg/                      # 常量与工具函数
├── scripts/
│   └── init_admin.sql        # 初始化 SQL
└── docs/                     # Swagger 文档
```

## 🚀 快速开始

### 1. 环境准备

- Go 1.24+
- PostgreSQL 12+
- Redis 6+

### 2. 配置数据库

创建 PostgreSQL 数据库：
```sql
CREATE DATABASE bamboo_main;
CREATE USER bamboo_main WITH PASSWORD 'bamboo_main';
GRANT ALL PRIVILEGES ON DATABASE bamboo_main TO bamboo_main;
```

### 3. 配置环境变量

复制配置模板并按需填写：
```bash
cp .env.example .env
```

常用配置项：
- `XLF_DEBUG` / `XLF_PORT`
- `DATABASE_HOST` / `DATABASE_PORT` / `DATABASE_USER` / `DATABASE_PASS` / `DATABASE_NAME`
- `NOSQL_HOST` / `NOSQL_PORT` / `NOSQL_DATABASE` / `NOSQL_PREFIX`
- `EMAIL_HOST` / `EMAIL_PORT` / `EMAIL_USER` / `EMAIL_PASS`

### 4. 安装依赖

```bash
go mod tidy
```

### 5. 初始化数据库

运行应用后，GORM 会自动创建表结构：
```bash
go run main.go
```

然后执行初始化 SQL 创建管理员账户：
```bash
psql -h localhost -U bamboo_main -d bamboo_main -f scripts/init_admin.sql
```

### 6. 访问系统

- **API 服务**: http://localhost:23333
- **API 文档**: http://localhost:23333/swagger/index.html
- **健康检查**: http://localhost:23333/api/v1/public/health

## 📚 API 接口

### 认证相关
- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/logout` - 用户登出
- `GET /api/v1/auth/user` - 获取用户信息
- `POST /api/v1/auth/password/change` - 修改密码
- `POST /api/v1/auth/password/reset` - 重置密码

### 友情链接管理
- `POST /api/v1/admin/links` - 添加友情链接
- `GET /api/v1/admin/links` - 获取友情链接列表
- `GET /api/v1/admin/links/{uuid}` - 获取友情链接详情
- `PUT /api/v1/admin/links/{uuid}` - 更新友情链接
- `DELETE /api/v1/admin/links/{uuid}` - 删除友情链接
- `PUT /api/v1/admin/links/{uuid}/status` - 更新链接状态
- `PUT /api/v1/admin/links/{uuid}/fail` - 更新失效状态

### 公开接口
- `GET /api/v1/public/links` - 获取公开友情链接
- `GET /api/v1/public/health` - 健康检查
- `GET /api/v1/public/ping` - Ping 测试

## 🔐 认证方式

系统使用自定义 Token 认证：

1. 登录成功后获得 token（格式：`cs_` + 64位字符串）
2. 请求头添加：`Authorization: Bearer {token}`
3. Token 存储在 Redis 中，默认有效期 24 小时
4. Redis Key 格式：`{NOSQL_PREFIX}:auth:token:{token}`（默认 `bm:auth:token:{token}`）

## 📊 Redis 常量规范

项目使用统一的 Redis Key 命名规范：

```go
type RedisKey string

const (
    RedisAuthToken   RedisKey = "auth:token:%s"
    RedisMailQueue   RedisKey = "mail:queue"
    RedisLinkFriend  RedisKey = "link:friend:%d"
    RedisSponsorChan RedisKey = "sponsor:channel:%d"
)

key := RedisAuthToken.Get(token).String() // => bm:auth:token:{token}
```

## 🔧 开发相关

### 生成 Swagger 文档
```bash
swag init -g main.go
```

### 运行测试
```bash
go test ./...
```

### 格式化代码
```bash
go fmt ./...
```

## 默认账户

- **用户名**: admin  
- **密码**: admin123456
- **邮箱**: admin@example.com
- **角色**: admin

## 🤝 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 📄 许可证

MIT License - 详见 [LICENSE](LICENSE) 文件

## 🙏 致谢

- 基于 [demo](demo/) 项目的架构设计
- 从 [old](old/) 项目迁移业务逻辑
- 使用 [bamboo-base-go](https://github.com/bamboo-services/bamboo-base-go) 基础库
