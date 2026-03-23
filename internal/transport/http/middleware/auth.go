package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/luckysxx/go-note/internal/auth"
)

// JWTAuth 鉴权中间件：通过 gRPC 调用 user-platform 验证 Token
func JWTAuth(authClient *auth.AuthClient, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 提取 Bearer Token
		token, err := auth.ExtractBearerToken(c.GetHeader("Authorization"))
		if err != nil {
			logger.Debug("请求鉴权拦截", zap.Error(err), zap.String("client_ip", c.ClientIP()))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "无效的访问凭证或已过期，请重新登录",
			})
			return
		}

		// 2. 调用 user-platform gRPC 验证 Token
		result, err := authClient.VerifyToken(c.Request.Context(), token)
		if err != nil {
			logger.Debug("Token 验证失败", zap.Error(err), zap.String("client_ip", c.ClientIP()))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "无效的访问凭证或已过期，请重新登录",
			})
			return
		}

		// 3. 将用户信息挂载到 Gin Context
		c.Set("userID", result.UserID)
		c.Set("username", result.Username)

		c.Next()
	}
}

// GetUserID 从 Context 中安全获取 userID
func GetUserID(c *gin.Context) (int64, bool) {
	val, exists := c.Get("userID")
	if !exists {
		return 0, false
	}

	userID, ok := val.(int64)
	return userID, ok
}
