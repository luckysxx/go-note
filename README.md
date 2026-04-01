# go-note

`go-note` 是一个提供笔记管理功能的微服务。它作为核心业务服务之一，支持 HTTP 和 gRPC 双协议通信，并无缝集成 [user-platform](../user-platform) 统一认证体系。

## 架构

```text
用户 → user-platform 登录 → 获取 Token
     → 携带 Token 访问 go-note API
     → go-note 通过 gRPC 调用 user-platform VerifyToken 验证身份
内部微服务 → go-note gRPC Server → 获取 / 操作笔记数据
```

- **协议**：支持 HTTP (REST) 与 gRPC 双协议，使用统一的 Handler 逻辑
- **认证**：通过 gRPC 委托给 user-platform，go-note 不持有 JWT Secret
- **其它依赖**：通过 gRPC 调用 id-generator 服务进行 ID 生成
- **ORM**：ent（与 user-platform 一致）
- **配置**：Viper + godotenv（YAML + 环境变量覆盖）
- **共享模块**：`github.com/luckysxx/common`（logger、errs、postgres、redis 连接池、otel 等）

## 技术栈

- Go、Gin、ent、PostgreSQL、Redis
- gRPC（双协议支持，对外提供服务，并与 user-platform / id-generator 通信）
- Viper（配置管理）
- `common` 基础设施支持（含 Redis 与 PostgreSQL 统一连接池及监控、OpenTelemetry链路追踪）

## 目录结构

```text
go-note/
├── cmd/
│   ├── http/main.go              # HTTP 服务入口
│   └── grpc/main.go              # gRPC 服务入口
├── configs/config.yaml           # Viper 配置
├── internal/
│   ├── auth/client.go            # gRPC 认证客户端 → user-platform
│   ├── cache/redis.go
│   ├── dberr/dberr.go
│   ├── ent/schema/snippet.go     # Ent Schema
│   ├── platform/{config,database}
│   ├── repository/snippet_repo.go
│   ├── service/{contract,snippet.go}
│   └── transport/                  # 传输层
│       ├── http/{dto,handler,router,middleware,response,errs,validator}
│       └── grpc/                   # gRPC 定义与 Server 实现
├── docker-compose.yaml
└── Makefile
```

## 快速开始

### 前置条件

- Go 1.25+
- PostgreSQL、Redis（可通过根目录 docker-compose 启动基础设施）
- user-platform gRPC 服务运行在 `localhost:9091`
- id-generator gRPC 服务运行在 `localhost:50059`

### 配置

```bash
cp .env.example .env
# 编辑 .env，设置 DATABASE_SOURCE 和 REDIS_PASSWORD
```

### 运行

```bash
# 启动 HTTP 服务
make run

# 启动 gRPC 服务
make run-grpc

# 编译后运行
make build
make build-grpc
./bin/go-note
./bin/go-note-grpc
```

### 常用命令

```bash
make help           # 查看全部命令
make run            # 启动 HTTP 服务
make run-grpc       # 启动 gRPC 服务
make build          # 编译 HTTP 二进制
make build-grpc     # 编译 gRPC 二进制
make test           # 运行测试
make ent-generate   # 重新生成 Ent 代码
make lint           # 代码检查
```

## API 端点

所有接口需携带 user-platform 签发的 Bearer Token：

| 方法 | 路径                 | 说明                 |
|------|----------------------|----------------------|
| GET  | `/health`            | 健康检查             |
| POST | `/api/v1/snippets`     | 创建笔记             |
| GET  | `/api/v1/snippets/:id` | 获取笔记             |
| PUT  | `/api/v1/snippets/:id` | 更新笔记             |
| GET  | `/api/v1/me/snippets`  | 获取我的笔记列表     |

## 服务端口

| 服务               | 端口  |
|--------------------|-------|
| go-note HTTP       | 8080  |
| go-note gRPC       | 9093  |
| metrics/healthz    | 9094  |
| user-platform gRPC | 9091  |
| id-generator gRPC  | 50059 |
