# 运行方式说明

## 方式 1：本地开发

### 前置条件

- PostgreSQL 和 Redis 已启动（通过 `docker-compose-infra.yaml` 启动公共基础设施）
- user-platform gRPC 服务运行在 `localhost:9091`

### 运行步骤

1. **配置环境变量**

```bash
cp .env.example .env
# 编辑 .env，设置 DATABASE_SOURCE 和 REDIS_PASSWORD
```

2. **启动服务**

```bash
make run
```

Ent 会在启动时自动执行 schema migration（`auto_migrate: true`），无需手动建表。

3. **访问服务**

- go-note API: `http://localhost:8080`
- 健康检查: `http://localhost:8080/health`

---

## 方式 2：Docker Compose

### 前置条件

- 公共基础设施已通过 `docker-compose-infra.yaml` 启动（PostgreSQL、Redis 等）
- `go-net` 网络已创建

### 运行步骤

1. **确保基础设施网络存在**

```bash
docker network inspect go-net >/dev/null 2>&1 || docker network create go-net
```

2. **在 global-postgres 中创建数据库**

```bash
docker exec global-postgres psql -U luckys -d postgres -c "CREATE DATABASE go_note;"
```

3. **启动服务**

```bash
make docker-up
```

4. **访问服务**

- go-note API: `http://localhost:8080`
- 前端: `http://localhost:8082`

### 配置说明

环境变量在 `docker-compose.yaml` 中定义，关键配置：

```yaml
environment:
  - DATABASE_SOURCE=postgres://luckys:123456@global-postgres:5432/go_note?sslmode=disable
  - REDIS_ADDR=global-redis:6379
  - USER_PLATFORM_ADDR=user-platform:9091
```

---

## 配置项

| 配置项                   | 说明               | 默认值        |
| ------------------------ | ------------------ | ------------- |
| `APP_ENV`                | 运行环境           | `development` |
| `SERVER_PORT`            | HTTP 端口          | `8080`        |
| `DATABASE_DRIVER`        | 数据库驱动         | `postgres`    |
| `DATABASE_SOURCE`        | 数据库连接串       | `.env` 中设置 |
| `DATABASE_AUTO_MIGRATE`  | 自动建表           | `true`        |
| `REDIS_ADDR`             | Redis 地址         | `localhost:6379` |
| `REDIS_PASSWORD`         | Redis 密码         | `.env` 中设置 |
| `USER_PLATFORM_ADDR`     | user-platform gRPC 地址 | `localhost:9091` |

---

## 常见问题

### Q: 启动报 "连接 user-platform gRPC 失败"

A: 确保 user-platform 的 gRPC 服务已启动并监听在配置的地址（默认 `localhost:9091`）。

### Q: 数据库表不存在

A: 检查 `configs/config.yaml` 中 `database.auto_migrate` 是否为 `true`。Ent 会在启动时自动创建表。

### Q: Token 验证失败

A: go-note 通过 gRPC 调用 user-platform 验证 Token。请确认：

1. 使用的是 user-platform 签发的有效 Token
2. user-platform gRPC 服务正常运行
3. `USER_PLATFORM_ADDR` 配置正确
