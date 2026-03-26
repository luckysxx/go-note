package auth

import (
	"context"
	"fmt"
	"strings"

	authpb "github.com/luckysxx/common/proto/auth"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// AuthClient 封装与 user-platform 的 gRPC 通信
type AuthClient struct {
	authClient authpb.AuthServiceClient
	conn       *grpc.ClientConn
}

// VerifyResult Token 验证结果
type VerifyResult struct {
	UserID   int64
	Username string
}

// LoginResult 登录结果
type LoginResult struct {
	AccessToken  string
	RefreshToken string
	UserID       int64
	Username     string
}

// RefreshTokenResult 刷新 Token 结果
type RefreshTokenResult struct {
	AccessToken  string
	RefreshToken string
}

// NewAuthClient 创建 gRPC 认证客户端
func NewAuthClient(addr string) (*AuthClient, error) {
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()), // OTel: 自动把 SpanContext 传给 user-platform gRPC
	)
	if err != nil {
		return nil, fmt.Errorf("连接 user-platform gRPC 失败: %w", err)
	}

	return &AuthClient{
		authClient: authpb.NewAuthServiceClient(conn),
		conn:       conn,
	}, nil
}

// Close 关闭 gRPC 连接
func (c *AuthClient) Close() error {
	return c.conn.Close()
}

// VerifyToken 调用 user-platform 验证 Token
func (c *AuthClient) VerifyToken(ctx context.Context, token string) (*VerifyResult, error) {
	resp, err := c.authClient.VerifyToken(ctx, &authpb.VerifyTokenRequest{Token: token})
	if err != nil {
		return nil, err
	}

	return &VerifyResult{
		UserID:   resp.GetUserId(),
		Username: resp.GetUsername(),
	}, nil
}

// Login 调用 user-platform 登录
func (c *AuthClient) Login(ctx context.Context, username, password, appCode string) (*LoginResult, error) {
	resp, err := c.authClient.Login(ctx, &authpb.LoginRequest{
		Username: username,
		Password: password,
		AppCode:  appCode,
	})
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		AccessToken:  resp.GetAccessToken(),
		RefreshToken: resp.GetRefreshToken(),
		UserID:       resp.GetUserId(),
		Username:     resp.GetUsername(),
	}, nil
}

// RefreshToken 调用 user-platform 刷新 Token
func (c *AuthClient) RefreshToken(ctx context.Context, token string) (*RefreshTokenResult, error) {
	resp, err := c.authClient.RefreshToken(ctx, &authpb.RefreshTokenRequest{Token: token})
	if err != nil {
		return nil, err
	}

	return &RefreshTokenResult{
		AccessToken:  resp.GetAccessToken(),
		RefreshToken: resp.GetRefreshToken(),
	}, nil
}

// ExtractBearerToken 从 Authorization Header 提取 Token
func ExtractBearerToken(authHeader string) (string, error) {
	authHeader = strings.TrimSpace(authHeader)
	if authHeader == "" {
		return "", fmt.Errorf("missing authorization header")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", fmt.Errorf("invalid authorization header format")
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", fmt.Errorf("empty token")
	}

	return token, nil
}
