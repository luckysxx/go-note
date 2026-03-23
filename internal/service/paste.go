package service

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"

	pkgerrs "github.com/luckysxx/common/errs"
	"github.com/luckysxx/go-note/internal/dberr"
	"github.com/luckysxx/go-note/internal/repository"
	servicecontract "github.com/luckysxx/go-note/internal/service/contract"

	"go.uber.org/zap"
)

// 领域错误定义
var (
	ErrShortLinkGeneration = pkgerrs.NewServerErr(errors.New("无法生成唯一 shortlink"))
	ErrForbidden           = pkgerrs.New(pkgerrs.Forbidden, "无权限操作该代码片段", nil)
)

// PasteService 定义 paste 业务接口
type PasteService interface {
	Create(ctx context.Context, userID int64, req *servicecontract.CreatePasteCommand) (*servicecontract.PasteResult, error)
	GetByID(ctx context.Context, id int64) (*servicecontract.PasteResult, error)
	ListMine(ctx context.Context, userID int64) ([]servicecontract.PasteResult, error)
	Update(ctx context.Context, userID, id int64, req *servicecontract.UpdatePasteCommand) (*servicecontract.PasteResult, error)
}

type pasteService struct {
	repo   repository.PasteRepository
	logger *zap.Logger
}

// NewPasteService 创建 PasteService 实例
func NewPasteService(repo repository.PasteRepository, logger *zap.Logger) PasteService {
	return &pasteService{repo: repo, logger: logger}
}

func (s *pasteService) Create(ctx context.Context, userID int64, req *servicecontract.CreatePasteCommand) (*servicecontract.PasteResult, error) {
	const maxAttempts = 5
	for i := 0; i < maxAttempts; i++ {
		cand, err := generateShortLink(8)
		if err != nil {
			s.logger.Error("生成 shortlink 失败", zap.Error(err))
			return nil, err
		}

		paste, err := s.repo.Create(ctx, userID, req, cand)
		if err == nil {
			return paste, nil
		}

		// 检查是否是唯一键冲突，如果是则重试
		if dberr.IsDuplicateKeyError(err) {
			s.logger.Warn("shortlink 冲突，重试生成", zap.String("shortLink", cand))
			continue
		}

		// 其他数据库错误，直接返回
		s.logger.Error("创建 paste 失败",
			zap.String("shortLink", cand),
			zap.Error(err),
		)
		return nil, err
	}

	s.logger.Error("无法生成唯一 shortlink")
	return nil, ErrShortLinkGeneration
}

func generateShortLink(n int) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		b[i] = letters[idx.Int64()]
	}
	return string(b), nil
}

func (s *pasteService) GetByID(ctx context.Context, id int64) (*servicecontract.PasteResult, error) {
	paste, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("查询 paste 失败",
			zap.Int64("id", id),
			zap.Error(err),
		)
		return nil, err
	}
	return paste, nil
}

func (s *pasteService) ListMine(ctx context.Context, userID int64) ([]servicecontract.PasteResult, error) {
	list, err := s.repo.ListByOwner(ctx, userID)
	if err != nil {
		s.logger.Error("查询用户 paste 列表失败",
			zap.Int64("userID", userID),
			zap.Error(err),
		)
		return nil, err
	}
	return list, nil
}

func (s *pasteService) Update(ctx context.Context, userID, id int64, req *servicecontract.UpdatePasteCommand) (*servicecontract.PasteResult, error) {
	current, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if current.OwnerID != userID {
		return nil, ErrForbidden
	}

	paste, err := s.repo.Update(ctx, userID, id, req)
	if err != nil {
		s.logger.Error("更新 paste 失败",
			zap.Int64("id", id),
			zap.Int64("userID", userID),
			zap.Error(err),
		)
		if dberr.IsNotFoundError(err) {
			return nil, ErrForbidden
		}
		return nil, err
	}

	return paste, nil
}
