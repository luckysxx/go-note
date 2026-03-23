# go-note

go-note 是一个 Pastebin 风格的代码分享服务，已接入 [user-platform](../user-platform) 统一认证体系（SSO）。

## 架构

```
用户 → user-platform 登录 → 获取 Token
     → 携带 Token 访问 go-note API
     → go-note 通过 gRPC 调用 user-platform VerifyToken 验证身份
```

- **认证**：通过 gRPC 委托给 user-platform，go-note 不持有 JWT Secret
- **ORM**：ent（与 user-platform 一致）
- **配置**：Viper + godotenv（YAML + 环境变量覆盖）
- **共享模块**：`github.com/luckysxx/common`（logger、errs）

## 技术栈

- Go、Gin、ent、PostgreSQL、Redis
- gRPC（与 user-platform 通信）
- Viper（配置管理）
- Vue 3 + Vite（前端）

## 目录结构

```
go-note/
├── cmd/http/main.go              # 入口（initInfra → buildRouter → runServer）
├── configs/config.yaml           # Viper 配置
├── internal/
│   ├── auth/client.go            # gRPC 认证客户端 → user-platform
│   ├── cache/redis.go
│   ├── dberr/dberr.go
│   ├── ent/schema/paste.go       # Ent Schema
│   ├── platform/{config,database}
│   ├── repository/paste_repo.go
│   ├── service/{contract,paste.go}
│   └── transport/http/{dto,handler,router,middleware,response,errs,validator}
├── view/
├── k8s/
├── docker-compose.yaml
└── Makefile
```

## 快速开始

### 前置条件

- Go 1.25+
- PostgreSQL、Redis（可通过根目录 docker-compose 启动基础设施）
- user-platform gRPC 服务运行在 `localhost:9091`

### 配置

```bash
cp .env.example .env
# 编辑 .env，设置 DATABASE_SOURCE 和 REDIS_PASSWORD
```

### 运行

```bash
# 直接运行
make run

# 或编译后运行
make build
./bin/go-note
```

### 常用命令

```bash
make help           # 查看全部命令
make run            # 启动服务
make build          # 编译二进制
make test           # 运行测试
make ent-generate   # 重新生成 Ent 代码
make lint           # 代码检查
```

## API 端点

所有接口需携带 user-platform 签发的 Bearer Token：

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/health` | 健康检查 |
| POST | `/api/v1/pastes` | 创建代码片段 |
| GET | `/api/v1/pastes/:id` | 获取代码片段 |
| PUT | `/api/v1/pastes/:id` | 更新代码片段 |
| GET | `/api/v1/me/pastes` | 获取我的代码片段列表 |

## 服务端口

| 服务 | 端口 |
|------|------|
| go-note HTTP | 8080 |
| user-platform gRPC | 9091 |
