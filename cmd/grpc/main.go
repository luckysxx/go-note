package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	grpchealth "google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"

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
	transportgrpc "github.com/luckysxx/go-note/internal/transport/grpc"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	_ = godotenv.Load()

	log := logger.NewLogger("go-note-grpc")
	defer log.Sync()

	cfg := config.LoadConfig()

	entClient, redisClient, idgenClient := initInfra(cfg, log)
	defer entClient.Close()
	defer redisClient.Close()
	defer idgenClient.Close()

	otelShutdown, err := commonOtel.InitTracer(cfg.OTel.ServiceName, cfg.OTel.JaegerEndpoint)
	if err != nil {
		log.Fatal("初始化 OpenTelemetry 失败", zap.Error(err))
	}
	defer otelShutdown(context.Background())

	healthChecker := buildHealthChecker(entClient, redisClient)
	grpcHealthServer := grpchealth.NewServer()
	startHealthSync(healthChecker, grpcHealthServer, log)

	snippetSvc := buildServices(entClient, idgenClient, log)
	grpcServer := transportgrpc.SetupServer(snippetSvc, grpcHealthServer, log)
	adminServer := buildAdminServer(cfg, healthChecker)

	runServer(grpcServer, grpcHealthServer, adminServer, cfg.GRPCServer.Port, log)
}

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

func buildServices(entClient *ent.Client, idgenClient platformidgen.Client, log *zap.Logger) service.SnippetService {
	snippetRepo := repository.NewSnippetRepository(entClient)
	return service.NewSnippetService(snippetRepo, idgenClient, log)
}

func buildHealthChecker(entClient *ent.Client, redisClient *redis.Client) *health.Checker {
	healthChecker := health.NewChecker()
	healthChecker.AddCheck("postgres", func(ctx context.Context) error {
		_, err := entClient.Snippet.Query().Limit(1).Count(ctx)
		return err
	})
	healthChecker.AddCheck("redis", func(ctx context.Context) error {
		return redisClient.Ping(ctx).Err()
	})
	return healthChecker
}

func buildAdminServer(cfg *config.Config, healthChecker *health.Checker) *http.Server {
	mux := http.NewServeMux()
	healthChecker.RegisterHTTP(mux)
	mux.Handle("/metrics", promhttp.Handler())

	return &http.Server{
		Addr:    ":" + cfg.Metrics.Port,
		Handler: mux,
	}
}

func startHealthSync(checker *health.Checker, healthServer *grpchealth.Server, log *zap.Logger) {
	var lastStatus healthgrpc.HealthCheckResponse_ServingStatus
	var initialized bool

	update := func() {
		allHealthy, results := checker.Evaluate(context.Background())
		statusCode := healthgrpc.HealthCheckResponse_SERVING
		if !allHealthy {
			statusCode = healthgrpc.HealthCheckResponse_NOT_SERVING
		}

		healthServer.SetServingStatus("", statusCode)
		healthServer.SetServingStatus("note.NoteService", statusCode)

		if initialized && statusCode == lastStatus {
			return
		}
		lastStatus = statusCode
		initialized = true

		if allHealthy {
			log.Debug("gRPC health 状态已更新", zap.String("status", statusCode.String()))
			return
		}

		log.Warn("gRPC health 状态已更新",
			zap.String("status", statusCode.String()),
			zap.Any("checks", results),
		)
	}

	update()

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			update()
		}
	}()
}

func runServer(s *grpc.Server, healthServer *grpchealth.Server, adminServer *http.Server, port string, log *zap.Logger) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("gRPC 端口监听失败", zap.Error(err))
	}

	go func() {
		log.Info("gRPC 服务已启动", zap.String("port", port))
		if err := s.Serve(lis); err != nil {
			log.Fatal("gRPC 服务异常终止", zap.Error(err))
		}
	}()

	go func() {
		log.Info("gRPC 管理端口已启动",
			zap.String("port", adminServer.Addr),
			zap.String("endpoints", "/metrics, /healthz, /readyz"),
		)
		if err := adminServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("gRPC 管理端口异常终止", zap.Error(err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Info("收到停机信号，开始优雅退出...")
	healthServer.SetServingStatus("", healthgrpc.HealthCheckResponse_NOT_SERVING)
	healthServer.SetServingStatus("note.NoteService", healthgrpc.HealthCheckResponse_NOT_SERVING)
	s.GracefulStop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := adminServer.Shutdown(shutdownCtx); err != nil {
		log.Fatal("gRPC 管理端口强制退出", zap.Error(err))
	}

	log.Info("gRPC 服务已安全退出")
}
