# syntax=docker/dockerfile:1
# --------------------------------------------------------------------------------
# Copyright (c) 2016-NOW(至今) 筱锋
# Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
# --------------------------------------------------------------------------------
# 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
# 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
# https://opensource.org/licenses/MIT
# --------------------------------------------------------------------------------

# =============================================================================
# Bamboo-Main 多阶段构建
#   Stage 1  frontend-builder : Node 22 + pnpm → resources/frontend/dist
#   Stage 2  go-builder       : golang 1.25 alpine → 单二进制（CGO_ENABLED=0）
#   Stage 3  runtime          : alpine + chromium（无头浏览器截图）
# =============================================================================

# ---------- Stage 1: 前端构建 ----------
FROM node:22-alpine AS frontend-builder
WORKDIR /build

# 仅复制依赖清单，最大化层缓存；pnpm 版本与本地开发对齐（11.8.0）
COPY frontend/package.json frontend/pnpm-lock.yaml frontend/pnpm-workspace.yaml ./frontend/
RUN npm install -g pnpm@11.8.0
RUN cd frontend && pnpm install --frozen-lockfile

# 复制前端源码并构建（vite outDir=../resources/frontend/dist，输出到 /build/resources/frontend/dist）
COPY frontend/ frontend/
RUN cd frontend && pnpm build

# ---------- Stage 2: Go 编译 ----------
FROM golang:1.25-alpine AS go-builder
WORKDIR /src

# 复制依赖清单；若含本地路径 replace（CI/容器内无该目录），先剥离再预下载（层缓存友好）
COPY go.mod go.sum ./
# 剥离本地路径 replace 后预下载依赖（层缓存友好；正式环境无 replace 时此段自动跳过）
RUN if grep -q 'phalanx-labs/beacon-sso-sdk => /' go.mod; then \
      grep -v 'phalanx-labs/beacon-sso-sdk => /' go.mod > go.mod.tmp && mv go.mod.tmp go.mod; \
    fi && go mod download

# 复制源码（.dockerignore 已排除本地 dist / node_modules / data 等；会覆盖上一步的 go.mod/go.sum）
COPY . .

# 【SSO SDK 兜底】COPY 恢复的 go.mod 若仍含本地路径 replace，再次剥离并补齐 go.sum。
# 正式做法是发布 beacon-sso-sdk 新版本并移除 replace，此处保证 CI 在任何状态下都能过 go mod download。
RUN if grep -q 'phalanx-labs/beacon-sso-sdk => /' go.mod; then \
      grep -v 'phalanx-labs/beacon-sso-sdk => /' go.mod > go.mod.tmp && mv go.mod.tmp go.mod; \
    fi && go mod download

# 用前端阶段产物覆盖，确保 go:embed all:dist 内嵌最新产物
COPY --from=frontend-builder /build/resources/frontend/dist /src/resources/frontend/dist

# CGO_ENABLED=0 静态编译（pgx / sonic 回退 / gopsutil 均纯 Go，不依赖 cgo）
# GOOS/GOARCH 不写死，随 golang 镜像目标架构（amd64/arm64）自动匹配，配合 buildx 多平台构建
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/bamboo-main .

# ---------- Stage 3: 运行镜像 ----------
FROM alpine:3.21 AS runtime

# 无头浏览器 + 中文字体 + 字体配置 + CA 证书 + 优雅退出 init + 时区数据
RUN apk add --no-cache \
    chromium \
    font-noto-cjk \
    fontconfig \
    ca-certificates \
    tini \
    tzdata

# 截图服务显式指定 chromium 路径；邮件模板/截图目录为相对进程工作目录的路径，固定 WORKDIR=/app
ENV SCREENSHOT_CHROME_PATH=/usr/bin/chromium \
    EMAIL_TEMPLATE_DIR=templates/mail \
    SCREENSHOT_DIR=data/screenshots \
    TZ=Asia/Shanghai

WORKDIR /app

COPY --from=go-builder /out/bamboo-main /app/bamboo-main
COPY templates/mail /app/templates/mail
RUN mkdir -p /app/data/screenshots

EXPOSE 23333

# tini 作为 PID 1 转发 SIGTERM/SIGINT，保证框架 Runner 优雅退出
ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/app/bamboo-main"]
