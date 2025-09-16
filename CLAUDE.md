# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Bamboo Main** is a friend links (友情链接) management system built with Go, migrated from GoFrame architecture to a modern Gin-based stack.

### Technology Stack
- **Web Framework**: Gin
- **ORM**: GORM v2
- **Database**: PostgreSQL
- **Cache**: Redis
- **Authentication**: Custom Token + Redis Sessions
- **Documentation**: Swagger
- **Configuration**: YAML-based

### Common Commands

```bash
# Build and run
go build .                    # Build the application
go run main.go                # Run the application (auto-creates DB tables)
go mod tidy                   # Install/update dependencies

# Development tools
swag init -g main.go          # Generate Swagger documentation
go test ./...                 # Run all tests
go vet ./...                  # Static code analysis
go fmt ./...                  # Format code

# Database setup
psql -h localhost -U bamboo_main -d bamboo_main -f scripts/init_admin.sql
```

## Architecture Overview

### Clean Architecture Pattern
The project follows a **Handler → Service → Logic → Model** layered architecture:

- **Handler Layer** (`internal/handler/`): HTTP request handling and response formatting
- **Service Layer** (`internal/service/`): Interface definitions with Logic implementations
- **Logic Layer** (`internal/logic/`): Core business logic and data processing
- **Model Layer** (`internal/model/`): Data structures (Entity, DTO, Request, Response)

### Key Dependencies

- **bamboo-base-go**: Custom base library providing error handling, utilities, and context management
- **Context-based DI**: Database and user sessions accessed via `gin.Context`
- **Redis Integration**: Session storage and caching with structured key patterns

### Project Structure

```
internal/
├── handler/          # HTTP request handlers
├── service/          # Service interfaces
├── logic/           # Business logic implementation
├── model/           # Data models
│   ├── entity/      # Database entities
│   ├── dto/         # Data transfer objects
│   ├── request/     # Request structures
│   └── response/    # Response structures
├── middleware/      # HTTP middleware
└── router/          # Route definitions

pkg/
├── constants/       # Application constants
├── startup/         # Application initialization
└── util/           # Utility functions
```

### API Structure
- **Base Path**: `/api/v1`
- **Admin Routes**: `/api/v1/admin/*` (requires authentication + admin role)
- **Auth Routes**: `/api/v1/auth/*` (login/logout/user management)
- **Public Routes**: `/api/v1/public/*` (health checks, public content)

### Database Design
- **Primary Keys**: All entities use PostgreSQL native `uuid` type (not char(36))
- **GORM v2**: PostgreSQL driver with auto-migration
- **Soft Deletes**: Built-in GORM soft delete support
- **Relations**: Foreign keys use `*uuid.UUID` for optional relationships
- **System Config**: `entity.System` provides key-value configuration storage

### Authentication Architecture
- **Token Storage**: `dtoRedis.TokenDTO` in Redis with 24h expiration
- **Context Flow**: Middleware validates token → stores UserUUID → handlers query DB
- **Real-time Security**: User status changes (disabled/role changes) take effect immediately
- **Context Access**: Use `ctxUtil.GetUserUUID(c)` to get authenticated user UUID

### Error Handling Standards
- **Handler Layer**: Use `xError.NewError()` for structured errors, never `errors.New()`
- **Error Transmission**: Pass errors via `_ = c.Error(err)` with `*xError.Error` type
- **Success Responses**: Use `xResult.Success()` or `xResult.SuccessHasData()` only
- **Error Codes**: Use predefined ErrorCode constants from bamboo-base-go

### Redis Patterns
Key format: `bm:{category}:{type}:{identifier}` (defined in `pkg/constants/redis.go`)
- Authentication: `bm:auth:token:{token}`
- Caching: `bm:link:cache:{uuid}`, `bm:group:cache:{uuid}`

### Environment Setup
- **Go Version**: 1.24+
- **Database**: PostgreSQL 12+
- **Cache**: Redis 6+
- **Default Port**: 23333
- **Admin Account**: `xiao_lfeng` / `xiao_lfeng` (auto-created via system config)
- **Config**: `configs/config.yaml`

### Business Domain
Friend links management system with:
- Link submission and approval workflow
- Administrative review interface
- Grouping and color categorization
- Link health monitoring (working/broken status)
- Public display API for approved links