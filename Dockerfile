# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /build

ENV GOPROXY=https://goproxy.cn,direct

# 拷贝 go.mod/go.sum 先缓存依赖
COPY go.mod go.sum ./
RUN go mod download

# 拷贝源码并编译
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o go-note ./cmd/http/main.go

# Run stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app
COPY --from=builder /build/go-note .
COPY --from=builder /build/configs ./configs

EXPOSE 8080

CMD ["./go-note"]
