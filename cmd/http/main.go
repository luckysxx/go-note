package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"

	"github.com/luckysxx/common/health"
	"github.com/luckysxx/common/logger"
	commonOtel "github.com/luckysxx/common/otel"
	commonRedis "github.com/luckysxx/common/redis"
	"github.com/luckysxx/go-note/internal/ent"
	"github.com/luckysxx/go-note/internal/platform/config"
	"github.com/luckysxx/go-note/internal/platform/database"
	platformidgen "github.com/luckysxx/go-note/internal/platform/idgen"
	"github.com/luckysxx/go-note/internal/repository"
	"github.com/luckysxx/go-note/internal/service"
	"github.com/luckysxx/go-note/internal/transport/http/handler"
	httprouter "github.com/luckysxx/go-note/internal/transport/http/router"
	"go.uber.org/zap"
)

func main() {
	// 先加载 .env 使 APP_ENV 生效（影响日志格式和颜色）
	_ = godotenv.Load()

	log := logger.NewLogger("go-note")
	defer log.Sync()

	cfg := config.LoadConfig()

	// 1. 初始化底层基础设施
	entClient, redisClient, idgenClient := initInfra(cfg, log)
	defer entClient.Close()
	defer redisClient.Close()
	defer idgenClient.Close()

	// 2. 初始化 OpenTelemetry 链路追踪
	otelShutdown, err := commonOtel.InitTracer(cfg.OTel.ServiceName, cfg.OTel.JaegerEndpoint)
	if err != nil {
		log.Fatal("初始化 OpenTelemetry 失败", zap.Error(err))
	}
	defer otelShutdown(context.Background())

	// 3. 依赖注入与组件装配
	router := buildRouter(cfg, entClient, redisClient, idgenClient, log)

	// 4. 阻塞运行与优雅停机
	runServer(router, cfg.Server.Port, log)
}

// initInfra 初始化基础设施
func initInfra(cfg *config.Config, log *zap.Logger) (*ent.Client, *redis.Client, platformidgen.Client) {
	entClient := database.InitEntClient(cfg.Database.Driver, cfg.Database.Source, cfg.Database.AutoMigrate, log)
	redisClient := commonRedis.Init(commonRedis.Config{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	}, log)
	idgenClient, err := platformidgen.New(cfg.IDGenerator.Addr)
	if err != nil {
		log.Fatal("初始化 id-generator 客户端失败", zap.Error(err))
	}

	return entClient, redisClient, idgenClient
}

// buildRouter 依赖注入装配
func buildRouter(cfg *config.Config, entClient *ent.Client, redisClient *redis.Client, idgenClient platformidgen.Client, log *zap.Logger) *gin.Engine {
	// Repository
	snippetRepo := repository.NewSnippetRepository(entClient)

	// Service
	snippetSvc := service.NewSnippetService(snippetRepo, idgenClient, log)

	// Transport
	snippetHandler := handler.NewSnippetHandler(snippetSvc, log)
	r := gin.New()

	// 健康检查（注册在业务中间件之前）
	healthChecker := health.NewChecker()
	healthChecker.AddCheck("postgres", func(ctx context.Context) error {
		// 用 count 探测数据库连接，避免 Scan 目标与列数不匹配。
		_, err := entClient.Snippet.Query().Limit(1).Count(ctx)
		return err
	})
	healthChecker.AddCheck("redis", func(ctx context.Context) error {
		return redisClient.Ping(ctx).Err()
	})
	healthChecker.Register(r)

	httprouter.SetupRouter(r, snippetHandler, log)

	return r
}

// runServer 启动 HTTP 服务器，监听停机信号后优雅退出
func runServer(router *gin.Engine, port string, log *zap.Logger) {
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		log.Info("HTTP 服务已启动", zap.String("port", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("HTTP 服务监听失败", zap.Error(err))
		}
	}()

	// 监听停机信号
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Info("收到停机信号，开始优雅退出...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("HTTP 服务强制退出", zap.Error(err))
	}

	log.Info("所有服务已安全退出")
}
