package service

import (
	"context"
	"errors"
	"strings"

	pkgerrs "github.com/luckysxx/common/errs"
	commonlogger "github.com/luckysxx/common/logger"
	"github.com/luckysxx/go-note/internal/platform/idgen"
	"github.com/luckysxx/go-note/internal/repository"
	servicecontract "github.com/luckysxx/go-note/internal/service/contract"

	"go.uber.org/zap"
)

// 领域错误定义
var (
	ErrIDGeneration = pkgerrs.NewServerErr(errors.New("生成片段 ID 失败"))
	ErrForbidden    = pkgerrs.New(pkgerrs.Forbidden, "无权限操作该代码片段", nil)
)

// SnippetService 定义 snippet 业务接口。
type SnippetService interface {
	Create(ctx context.Context, userID int64, req *servicecontract.CreateSnippetCommand) (*servicecontract.SnippetResult, error)
	GetByID(ctx context.Context, id int64) (*servicecontract.SnippetResult, error)
	ListMine(ctx context.Context, userID int64) ([]servicecontract.SnippetResult, error)
	Update(ctx context.Context, userID, id int64, req *servicecontract.UpdateSnippetCommand) (*servicecontract.SnippetResult, error)
	Delete(ctx context.Context, userID, id int64) error
}

type snippetService struct {
	repo   repository.SnippetRepository
	idgen  idgen.Client
	logger *zap.Logger
}

// NewSnippetService 创建 SnippetService 实例。
func NewSnippetService(repo repository.SnippetRepository, idgenClient idgen.Client, logger *zap.Logger) SnippetService {
	return &snippetService{repo: repo, idgen: idgenClient, logger: logger}
}

func (s *snippetService) Create(ctx context.Context, userID int64, req *servicecontract.CreateSnippetCommand) (*servicecontract.SnippetResult, error) {
	params := normalizeCreateCommand(req)
	if err := validateCreateCommand(params); err != nil {
		return nil, err
	}

	id, err := s.idgen.NextID(ctx)
	if err != nil {
		commonlogger.Ctx(ctx, s.logger).Error("生成 snippet ID 失败", zap.Error(err))
		return nil, ErrIDGeneration
	}

	snippet, err := s.repo.Create(ctx, id, userID, params)
	if err != nil {
		commonlogger.Ctx(ctx, s.logger).Error("创建 snippet 失败",
			zap.Int64("id", id),
			zap.Int64("userID", userID),
			zap.Error(err),
		)
		return nil, err
	}

	return snippet, nil
}

func (s *snippetService) GetByID(ctx context.Context, id int64) (*servicecontract.SnippetResult, error) {
	snippet, err := s.repo.GetByID(ctx, id)
	if err != nil {
		commonlogger.Ctx(ctx, s.logger).Error("查询 snippet 失败",
			zap.Int64("id", id),
			zap.Error(err),
		)
		return nil, err
	}
	return snippet, nil
}

func (s *snippetService) ListMine(ctx context.Context, userID int64) ([]servicecontract.SnippetResult, error) {
	list, err := s.repo.ListByOwner(ctx, userID)
	if err != nil {
		commonlogger.Ctx(ctx, s.logger).Error("查询用户 snippet 列表失败",
			zap.Int64("userID", userID),
			zap.Error(err),
		)
		return nil, err
	}
	return list, nil
}

func (s *snippetService) Update(ctx context.Context, userID, id int64, req *servicecontract.UpdateSnippetCommand) (*servicecontract.SnippetResult, error) {
	current, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if current.OwnerID != userID {
		return nil, ErrForbidden
	}

	params := normalizeUpdateCommand(current, req)
	if err := validateUpdateCommand(params); err != nil {
		return nil, err
	}

	snippet, err := s.repo.Update(ctx, userID, id, params)
	if err != nil {
		commonlogger.Ctx(ctx, s.logger).Error("更新 snippet 失败",
			zap.Int64("id", id),
			zap.Int64("userID", userID),
			zap.Error(err),
		)
		return nil, err
	}

	return snippet, nil
}

func (s *snippetService) Delete(ctx context.Context, userID, id int64) error {
	current, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if current.OwnerID != userID {
		return ErrForbidden
	}

	if err := s.repo.Delete(ctx, userID, id); err != nil {
		commonlogger.Ctx(ctx, s.logger).Error("删除 snippet 失败",
			zap.Int64("id", id),
			zap.Int64("userID", userID),
			zap.Error(err),
		)
		return err
	}

	return nil
}

func normalizeCreateCommand(req *servicecontract.CreateSnippetCommand) *servicecontract.CreateSnippetCommand {
	if req == nil {
		return &servicecontract.CreateSnippetCommand{
			Type:       "code",
			Language:   "text",
			Visibility: "private",
		}
	}

	params := *req
	params.Type = normalizeType(params.Type)
	params.Title = strings.TrimSpace(params.Title)
	params.FileURL = strings.TrimSpace(params.FileURL)
	params.MimeType = strings.TrimSpace(params.MimeType)
	params.Visibility = normalizeVisibility(params.Visibility)
	params.Language = normalizeLanguage(params.Language)
	return &params
}

func normalizeUpdateCommand(current *servicecontract.SnippetResult, req *servicecontract.UpdateSnippetCommand) *servicecontract.UpdateSnippetCommand {
	params := &servicecontract.UpdateSnippetCommand{
		Title:      current.Title,
		Content:    current.Content,
		Language:   current.Language,
		Visibility: current.Visibility,
		GroupID:    current.GroupID,
	}
	if req == nil {
		return params
	}

	if trimmed := strings.TrimSpace(req.Title); trimmed != "" {
		params.Title = trimmed
	}
	if req.Content != "" {
		params.Content = req.Content
	}
	if req.Language != "" {
		params.Language = normalizeLanguage(req.Language)
	}
	if req.Visibility != "" {
		params.Visibility = normalizeVisibility(req.Visibility)
	}
	if req.GroupID != nil {
		params.GroupID = req.GroupID
	}

	return params
}

func validateCreateCommand(req *servicecontract.CreateSnippetCommand) error {
	if strings.TrimSpace(req.Title) == "" {
		return pkgerrs.NewParamErr("标题不能为空", nil)
	}
	if !isSupportedType(req.Type) {
		return pkgerrs.NewParamErr("type 仅支持 code、note、file", nil)
	}
	if !isSupportedVisibility(req.Visibility) {
		return pkgerrs.NewParamErr("visibility 仅支持 public 或 private", nil)
	}
	if req.Type == "file" && req.FileURL == "" {
		return pkgerrs.NewParamErr("文件片段必须提供 file_url", nil)
	}
	if req.Type != "file" && req.Content == "" {
		return pkgerrs.NewParamErr("文本片段必须提供 content", nil)
	}

	return nil
}

func validateUpdateCommand(req *servicecontract.UpdateSnippetCommand) error {
	if strings.TrimSpace(req.Title) == "" {
		return pkgerrs.NewParamErr("标题不能为空", nil)
	}
	if !isSupportedVisibility(req.Visibility) {
		return pkgerrs.NewParamErr("visibility 仅支持 public 或 private", nil)
	}
	return nil
}

func normalizeType(t string) string {
	switch strings.TrimSpace(t) {
	case "note":
		return "note"
	case "file":
		return "file"
	default:
		return "code"
	}
}

func normalizeVisibility(v string) string {
	if strings.TrimSpace(v) == "public" {
		return "public"
	}
	return "private"
}

func normalizeLanguage(v string) string {
	if strings.TrimSpace(v) == "" {
		return "text"
	}
	return strings.TrimSpace(v)
}

func isSupportedType(v string) bool {
	return v == "code" || v == "note" || v == "file"
}

func isSupportedVisibility(v string) bool {
	return v == "public" || v == "private"
}
