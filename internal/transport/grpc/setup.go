// Package transportgrpc 组装 gRPC Server，对标 HTTP 层的 router.go。
package transportgrpc

import (
	"time"

	commonlogger "github.com/luckysxx/common/logger"
	"github.com/luckysxx/common/metrics"
	notepb "github.com/luckysxx/common/proto/note"
	"github.com/luckysxx/go-note/internal/service"
	"github.com/luckysxx/go-note/internal/transport/grpc/interceptor"
	grpcserver "github.com/luckysxx/go-note/internal/transport/grpc/server"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
)

// SetupServer 组装 gRPC Server（对标 HTTP 的 SetupRouter）
func SetupServer(
	snippetSvc service.SnippetService,
	healthServer healthgrpc.HealthServer,
	log *zap.Logger,
) *grpc.Server {
	s := grpc.NewServer(
		// ── Keepalive ────────────────────────────────────────
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             5 * time.Second,
			PermitWithoutStream: true,
		}),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     5 * time.Minute,
			MaxConnectionAge:      30 * time.Minute,
			MaxConnectionAgeGrace: 10 * time.Second,
			Time:                  15 * time.Second,
			Timeout:               5 * time.Second,
		}),
		// ── Observability ────────────────────────────────────
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			metrics.GRPCMetricsInterceptor(),
			interceptor.GatewayAuthInterceptor(),
			commonlogger.GRPCUnaryServerInterceptor(log, interceptor.LogFieldsFromContext),
		),
	)

	notepb.RegisterNoteServiceServer(s, grpcserver.NewNoteServer(snippetSvc, log))
	healthgrpc.RegisterHealthServer(s, healthServer)

	return s
}
