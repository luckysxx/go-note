package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	grpchealth "google.golang.org/grpc/health"

	"github.com/luckysxx/common/logger"
	commonOtel "github.com/luckysxx/common/otel"
	"github.com/luckysxx/common/probe"
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

	otelShutdown, err := commonOtel.InitTracer(cfg.OTel)
	if err != nil {
		log.Fatal("初始化 OpenTelemetry 失败", zap.Error(err))
	}
	defer otelShutdown(context.Background())

	// 探针：独立管理端口 + gRPC Health 同步
	grpcHealthServer := grpchealth.NewServer()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	probeShutdown := probe.Serve(ctx, ":"+cfg.Metrics.Port, log,
		probe.WithCheck("postgres", func(ctx context.Context) error {
			_, err := entClient.Snippet.Query().Limit(1).Count(ctx)
			return err
		}),
		probe.WithRedis(redisClient),
		probe.WithGRPCHealth(grpcHealthServer, "note.NoteService"),
	)
	defer probeShutdown()

	snippetSvc := buildServices(entClient, idgenClient, log)
	grpcServer := transportgrpc.SetupServer(snippetSvc, grpcHealthServer, log)

	runServer(grpcServer, grpcHealthServer, cfg.GRPCServer.Port, log)
}

func initInfra(cfg *config.Config, log *zap.Logger) (*ent.Client, *redis.Client, platformidgen.Client) {
	entClient := database.InitEntClient(cfg.Database.Driver, cfg.Database.Source, cfg.Database.AutoMigrate, log)
	redisClient := commonRedis.Init(cfg.Redis, log)
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

func runServer(s *grpc.Server, healthServer *grpchealth.Server, port string, log *zap.Logger) {
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

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Info("收到停机信号，开始优雅退出...")
	probe.GRPCShutdown(healthServer, "note.NoteService")
	s.GracefulStop()

	log.Info("gRPC 服务已安全退出")
}
