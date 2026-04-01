package grpcserver

import (
	"context"

	commonlogger "github.com/luckysxx/common/logger"
	notepb "github.com/luckysxx/common/proto/note"
	"github.com/luckysxx/go-note/internal/service"
	servicecontract "github.com/luckysxx/go-note/internal/service/contract"
	grpcerrs "github.com/luckysxx/go-note/internal/transport/grpc/errs"
	"github.com/luckysxx/go-note/internal/transport/grpc/interceptor"
	"go.uber.org/zap"
)

// NoteServer 实现 NoteServiceServer 接口。
type NoteServer struct {
	notepb.UnimplementedNoteServiceServer
	snippetSvc service.SnippetService
	log        *zap.Logger
}

// NewNoteServer 创建一个 NoteService gRPC 服务实现。
func NewNoteServer(snippetSvc service.SnippetService, log *zap.Logger) *NoteServer {
	return &NoteServer{snippetSvc: snippetSvc, log: log}
}

func (s *NoteServer) CreateSnippet(ctx context.Context, req *notepb.CreateSnippetRequest) (*notepb.SnippetResponse, error) {
	userID, err := interceptor.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	result, err := s.snippetSvc.Create(ctx, userID, &servicecontract.CreateSnippetCommand{
		Type:       "code",
		Title:      req.Title,
		Content:    req.Content,
		Language:   req.Language,
		Visibility: req.Visibility,
	})
	if err != nil {
		commonlogger.Ctx(ctx, s.log).Error("gRPC CreateSnippet 失败", zap.Error(err))
		return nil, grpcerrs.ToStatusError(err)
	}
	return toProto(result), nil
}

func (s *NoteServer) GetSnippet(ctx context.Context, req *notepb.GetSnippetRequest) (*notepb.SnippetResponse, error) {
	result, err := s.snippetSvc.GetByID(ctx, req.SnippetId)
	if err != nil {
		commonlogger.Ctx(ctx, s.log).Error("gRPC GetSnippet 失败", zap.Int64("snippet_id", req.SnippetId), zap.Error(err))
		return nil, grpcerrs.ToStatusError(err)
	}
	return toProto(result), nil
}

func (s *NoteServer) ListSnippets(ctx context.Context, req *notepb.ListSnippetsRequest) (*notepb.ListSnippetsResponse, error) {
	userID, err := interceptor.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	results, err := s.snippetSvc.ListMine(ctx, userID)
	if err != nil {
		commonlogger.Ctx(ctx, s.log).Error("gRPC ListSnippets 失败", zap.Int64("user_id", userID), zap.Error(err))
		return nil, grpcerrs.ToStatusError(err)
	}

	snippets := make([]*notepb.SnippetResponse, len(results))
	for i := range results {
		snippets[i] = toProto(&results[i])
	}
	return &notepb.ListSnippetsResponse{Snippets: snippets}, nil
}

func (s *NoteServer) UpdateSnippet(ctx context.Context, req *notepb.UpdateSnippetRequest) (*notepb.SnippetResponse, error) {
	userID, err := interceptor.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	result, err := s.snippetSvc.Update(ctx, userID, req.SnippetId, &servicecontract.UpdateSnippetCommand{
		Title:      req.Title,
		Content:    req.Content,
		Language:   req.Language,
		Visibility: req.Visibility,
	})
	if err != nil {
		commonlogger.Ctx(ctx, s.log).Error("gRPC UpdateSnippet 失败", zap.Error(err))
		return nil, grpcerrs.ToStatusError(err)
	}
	return toProto(result), nil
}

// toProto 将业务层结果转为 protobuf 响应。
func toProto(r *servicecontract.SnippetResult) *notepb.SnippetResponse {
	resp := &notepb.SnippetResponse{
		Id:         r.ID,
		OwnerId:    r.OwnerID,
		Title:      r.Title,
		Type:       r.Type,
		Content:    r.Content,
		FileUrl:    r.FileURL,
		FileSize:   r.FileSize,
		MimeType:   r.MimeType,
		Language:   r.Language,
		Visibility: r.Visibility,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}
	if r.GroupID != nil {
		resp.GroupId = *r.GroupID
	}
	return resp
}
