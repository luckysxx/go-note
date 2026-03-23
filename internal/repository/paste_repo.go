package repository

import (
	"context"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/luckysxx/go-note/internal/dberr"
	"github.com/luckysxx/go-note/internal/ent"
	"github.com/luckysxx/go-note/internal/ent/paste"
	servicecontract "github.com/luckysxx/go-note/internal/service/contract"
)

// PasteRepository 定义 paste 数据访问接口
type PasteRepository interface {
	Create(ctx context.Context, ownerID int64, params *servicecontract.CreatePasteCommand, shortLink string) (*servicecontract.PasteResult, error)
	GetByID(ctx context.Context, id int64) (*servicecontract.PasteResult, error)
	ListByOwner(ctx context.Context, ownerID int64) ([]servicecontract.PasteResult, error)
	Update(ctx context.Context, ownerID, id int64, params *servicecontract.UpdatePasteCommand) (*servicecontract.PasteResult, error)
}

type pasteRepository struct {
	client *ent.Client
}

// NewPasteRepository 创建 PasteRepository 实例
func NewPasteRepository(client *ent.Client) PasteRepository {
	return &pasteRepository{client: client}
}

func (r *pasteRepository) Create(ctx context.Context, ownerID int64, params *servicecontract.CreatePasteCommand, shortLink string) (*servicecontract.PasteResult, error) {
	builder := r.client.Paste.Create().
		SetOwnerID(ownerID).
		SetTitle(params.Title).
		SetContent(params.Content).
		SetLanguage(params.Language).
		SetVisibility(resolveVisibility(params.Visibility))

	if shortLink != "" {
		builder.SetShortLink(shortLink)
	}

	p, err := builder.Save(ctx)
	if err != nil {
		return nil, dberr.ParseDBError(err)
	}

	return toResult(p), nil
}

func (r *pasteRepository) GetByID(ctx context.Context, id int64) (*servicecontract.PasteResult, error) {
	p, err := r.client.Paste.Get(ctx, id)
	if err != nil {
		return nil, dberr.ParseDBError(err)
	}

	return toResult(p), nil
}

func (r *pasteRepository) ListByOwner(ctx context.Context, ownerID int64) ([]servicecontract.PasteResult, error) {
	pastes, err := r.client.Paste.
		Query().
		Where(paste.OwnerIDEQ(ownerID)).
		Order(paste.ByUpdatedAt(sql.OrderDesc())).
		All(ctx)
	if err != nil {
		return nil, dberr.ParseDBError(err)
	}

	results := make([]servicecontract.PasteResult, 0, len(pastes))
	for _, p := range pastes {
		results = append(results, *toResult(p))
	}

	return results, nil
}

func (r *pasteRepository) Update(ctx context.Context, ownerID, id int64, params *servicecontract.UpdatePasteCommand) (*servicecontract.PasteResult, error) {
	// 先确认记录存在且属于该用户
	p, err := r.client.Paste.
		Query().
		Where(paste.IDEQ(id), paste.OwnerIDEQ(ownerID)).
		Only(ctx)
	if err != nil {
		return nil, dberr.ParseDBError(err)
	}

	// 执行更新
	updated, err := p.Update().
		SetTitle(params.Title).
		SetContent(params.Content).
		SetLanguage(params.Language).
		SetVisibility(resolveVisibility(params.Visibility)).
		Save(ctx)
	if err != nil {
		return nil, dberr.ParseDBError(err)
	}

	return toResult(updated), nil
}

func resolveVisibility(visibility string) paste.Visibility {
	if visibility == "public" {
		return paste.VisibilityPublic
	}
	return paste.VisibilityPrivate
}

func toResult(p *ent.Paste) *servicecontract.PasteResult {
	result := &servicecontract.PasteResult{
		ID:         p.ID,
		OwnerID:    p.OwnerID,
		Title:      p.Title,
		Content:    p.Content,
		Language:   p.Language,
		Visibility: string(p.Visibility),
		CreatedAt:  p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  p.UpdatedAt.Format(time.RFC3339),
	}

	if p.ShortLink != "" {
		result.ShortLink = p.ShortLink
	}

	return result
}
