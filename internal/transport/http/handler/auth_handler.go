package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/luckysxx/go-note/internal/auth"
)

// AuthHandler 认证代理 Handler（转发到 user-platform gRPC）
type AuthHandler struct {
	authClient *auth.AuthClient
	logger     *zap.Logger
}

// NewAuthHandler 创建 AuthHandler
func NewAuthHandler(authClient *auth.AuthClient, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{authClient: authClient, logger: logger}
}

// LoginRequest 登录请求 DTO
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

const defaultAppCode = "go-note"

// Login 登录代理：转发到 user-platform AuthService.Login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("登录请求绑定失败", zap.Error(err), zap.String("content_type", c.ContentType()))
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请填写完整的登录信息"})
		return
	}

	result, err := h.authClient.Login(c.Request.Context(), req.Username, req.Password, defaultAppCode)
	if err != nil {
		h.handleGRPCError(c, err, "登录失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "登录成功",
		"data": gin.H{
			"token":         result.AccessToken,
			"refresh_token": result.RefreshToken,
			"user_id":       result.UserID,
			"username":      result.Username,
			"email":         "",
		},
	})
}

// RefreshTokenRequest 刷新 Token 请求 DTO
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshToken 刷新代理：转发到 user-platform AuthService.RefreshToken
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("刷新 Token 请求绑定失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	result, err := h.authClient.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		h.handleGRPCError(c, err, "刷新 Token 失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "刷新成功",
		"data": gin.H{
			"token":         result.AccessToken,
			"refresh_token": result.RefreshToken,
		},
	})
}

// handleGRPCError 将 gRPC 错误转换为 HTTP 响应
func (h *AuthHandler) handleGRPCError(c *gin.Context, err error, fallbackMsg string) {
	st, ok := status.FromError(err)
	if !ok {
		h.logger.Error("gRPC 调用失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": fallbackMsg})
		return
	}

	switch st.Code() {
	case codes.InvalidArgument:
		h.logger.Warn("gRPC InvalidArgument", zap.String("msg", st.Message()))
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": st.Message()})
	case codes.NotFound:
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": st.Message()})
	case codes.PermissionDenied, codes.Unauthenticated:
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": st.Message()})
	case codes.ResourceExhausted:
		c.JSON(http.StatusTooManyRequests, gin.H{"code": 429, "msg": st.Message()})
	default:
		h.logger.Error("gRPC 调用失败", zap.Error(err), zap.String("grpc_code", st.Code().String()))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": fallbackMsg})
	}
}
