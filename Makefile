SHELL := /bin/bash

.PHONY: help run run-grpc build build-grpc test lint ent-generate docker-up docker-down docker-build clean

help:
	@echo "Available targets:"
	@echo "  make run              # 启动 go-note HTTP 服务"
	@echo "  make run-grpc         # 启动 go-note gRPC 服务"
	@echo "  make build            # 编译二进制"
	@echo "  make build-grpc       # 编译 gRPC 二进制"
	@echo "  make test             # 运行所有测试"
	@echo "  make ent-generate     # 重新生成 Ent 代码"
	@echo "  make lint             # Go vet + build 检查"
	@echo "  make docker-up        # Docker Compose 启动"
	@echo "  make docker-down      # Docker Compose 停止"
	@echo "  make docker-build     # Docker Compose 构建并启动"
	@echo "  make clean            # 清理编译产物"

# ==========================================
# 开发
# ==========================================

run:
	@go run ./cmd/http/main.go

run-grpc:
	@go run ./cmd/grpc/main.go

build:
	@go build -o bin/go-note ./cmd/http/main.go
	@echo "Binary built: bin/go-note"

build-grpc:
	@go build -o bin/go-note-grpc ./cmd/grpc/main.go
	@echo "Binary built: bin/go-note-grpc"

test:
	@go test ./internal/...

lint:
	@go vet ./...
	@go build ./...
	@echo "Lint passed"

# ==========================================
# 代码生成
# ==========================================

ent-generate:
	@GOWORK=off go run -mod=mod entgo.io/ent/cmd/ent generate ./internal/ent/schema
	@echo "Ent code generated"

# ==========================================
# Docker
# ==========================================

docker-up:
	@docker compose up -d
	@echo "Services started"

docker-down:
	@docker compose down
	@echo "Services stopped"

docker-build:
	@docker compose up -d --build
	@echo "Services built and started"

# ==========================================
# 清理
# ==========================================

clean:
	@rm -rf bin/
	@echo "Cleaned"
