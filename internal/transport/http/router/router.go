package router

import (
	"github.com/luckysxx/go-note/internal/auth"
	"github.com/luckysxx/go-note/internal/transport/http/handler"
	"github.com/luckysxx/go-note/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SetupRouter 配置路由
func SetupRouter(r *gin.Engine, pasteHandler *handler.PasteHandler, authHandler *handler.AuthHandler, authClient *auth.AuthClient, log *zap.Logger) {
	r.Use(middleware.GinLogger(log))
	r.Use(middleware.GinRecovery(log, true))

	// 健康检查
	healthHandler := func(c *gin.Context) {
		c.String(200, "ok")
	}
	r.GET("/health", healthHandler)
	r.HEAD("/health", healthHandler)

	v1 := r.Group("/api/v1")
	{
		// 认证代理（无需鉴权，登录转发到 user-platform gRPC）
		// 注册走 user-platform 统一注册页面（前端跳转）
		users := v1.Group("/users")
		{
			users.POST("/login", authHandler.Login)
			users.POST("/refresh", authHandler.RefreshToken)
		}

		// 需要鉴权的接口
		me := v1.Group("/me")
		me.Use(middleware.JWTAuth(authClient, log))
		{
			me.GET("/pastes", pasteHandler.ListMine)
		}

		pastes := v1.Group("/pastes")
		pastes.Use(middleware.JWTAuth(authClient, log))
		{
			pastes.POST("", pasteHandler.Create)
			pastes.GET("/:id", pasteHandler.Get)
			pastes.PUT("/:id", pasteHandler.Update)
		}
	}
}
