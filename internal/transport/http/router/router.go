package router

import (
	"github.com/luckysxx/common/logger"
	"github.com/luckysxx/go-note/internal/transport/http/handler"
	"github.com/luckysxx/go-note/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/zap"
)

// SetupRouter 配置路由。
func SetupRouter(r *gin.Engine, snippetHandler *handler.SnippetHandler, log *zap.Logger) {
	r.Use(otelgin.Middleware("go-note"))
	r.Use(logger.GinLogger(log))
	r.Use(logger.GinRecovery(log, true))

	v1 := r.Group("/api/v1")
	{
		me := v1.Group("/me")
		me.Use(middleware.GatewayAuth(log))
		{
			me.GET("/snippets", snippetHandler.ListMine)
		}

		snippets := v1.Group("/snippets")
		snippets.Use(middleware.GatewayAuth(log))
		{
			snippets.POST("", snippetHandler.Create)
			snippets.GET("/:id", snippetHandler.Get)
			snippets.PUT("/:id", snippetHandler.Update)
			snippets.DELETE("/:id", snippetHandler.Delete)
		}
	}
}
