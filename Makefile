# 变量定义，方便后续维护
MAIN_FILE = main.go
OUTPUT_BIN = bamboo-main
SWAG_CMD = swag
SWAG_FLAGS = --parseDependency

.DEFAULT_GOAL := help

.PHONY: help install swag run dev dev-backend dev-frontend build-frontend tidy fmt test vet lint build generate

# 显示帮助信息
help:
	@echo "Bamboo Main · 友链管理后端 - 可用命令"
	@echo ""
	@echo "初始化命令:"
	@echo "  make install             - 安装前端依赖 (frontend/ 下执行 pnpm install)"
	@echo ""
	@echo "开发命令:"
	@echo "  make swag                - 生成 Swagger 文档"
	@echo "  make run                 - 运行后端程序"
	@echo "  make dev                 - 生成文档并运行后端 (跳过前端构建)"
	@echo "  make dev-backend         - 一键构建并运行后端 (包含前端构建)"
	@echo "  make dev-frontend        - 运行前端开发服务器 (端口 3000)"
	@echo ""
	@echo "构建命令:"
	@echo "  make generate            - 一键构建：前端打包 → Swagger → Go 编译（运行前请先 make install 确保依赖最新）"
	@echo "  make build               - 同 generate"
	@echo "  make build-frontend      - 仅构建前端 (产出 resources/frontend/dist)"
	@echo ""
	@echo "质量命令:"
	@echo "  make tidy                - 整理 Go 模块"
	@echo "  make fmt                 - 格式化 Go 代码"
	@echo "  make test                - 运行 Go 测试"
	@echo "  make vet                 - 运行 go vet 静态检查"
	@echo "  make lint                - 运行 golangci-lint (未安装则跳过)"
	@echo ""
	@echo "示例:"
	@echo "  make dev"
	@echo "  make dev-backend"
	@echo "  make build"
	@echo ""

# 安装前端依赖
install:
	cd frontend && pnpm install

# 提取出的 Swagger 生成目标
swag:
	$(SWAG_CMD) init -g $(MAIN_FILE) $(SWAG_FLAGS)

# 运行后端二进制
run:
	chmod +x $(OUTPUT_BIN) && ./$(OUTPUT_BIN)

# 整理 Go 模块
tidy:
	go mod tidy

# 格式化 Go 代码
fmt:
	go fmt ./...

# 运行 Go 测试
test:
	go test ./...

# 组合目标：先生成文档与前端，再运行后端程序
dev-backend: generate run

# 后端开发：生成 Swagger 文档后运行（跳过前端构建）
dev: swag run

# 前端开发服务器
dev-frontend:
	cd frontend && pnpm dev

# 一键构建：前端打包 → Swagger 文档 → Go 编译（产物内嵌进二进制，单文件即可启动）
generate: build-frontend swag
	go build -o $(OUTPUT_BIN)

# build 是 generate 的别名
build: generate

# 仅构建前端（产出 resources/frontend/dist）
build-frontend:
	cd frontend && pnpm install && pnpm build

# 静态检查
vet:
	go vet ./...

# 代码 lint（若 golangci-lint 未安装则优雅跳过）
lint:
	@which golangci-lint >/dev/null 2>&1 || { echo "golangci-lint not installed, skipping"; exit 0; }
	@golangci-lint run ./...
