package router

import (
	"github.com/luckysxx/common/metrics"
	"github.com/luckysxx/go-note/internal/transport/http/handler"
	"github.com/luckysxx/go-note/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/zap"
)

// SetupRouter 配置路由
func SetupRouter(r *gin.Engine, pasteHandler *handler.PasteHandler, log *zap.Logger) {
	r.GET("/metrics", metrics.GinMetricsHandler())
	r.Use(metrics.GinMetrics())
	r.Use(otelgin.Middleware("go-note"))
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
		// 需要鉴权的接口（信任网关传来的 X-User-Id）
		me := v1.Group("/me")
		me.Use(middleware.GatewayAuth(log))
		{
			me.GET("/pastes", pasteHandler.ListMine)
		}

		pastes := v1.Group("/pastes")
		pastes.Use(middleware.GatewayAuth(log))
		{
			pastes.POST("", pasteHandler.Create)
			pastes.GET("/:id", pasteHandler.Get)
			pastes.PUT("/:id", pasteHandler.Update)
		}
	}
}
