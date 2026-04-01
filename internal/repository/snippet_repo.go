package repository

import (
	"context"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/luckysxx/go-note/internal/dberr"
	"github.com/luckysxx/go-note/internal/ent"
	"github.com/luckysxx/go-note/internal/ent/snippet"
	servicecontract "github.com/luckysxx/go-note/internal/service/contract"
)

// SnippetRepository 定义 snippet 数据访问接口。
type SnippetRepository interface {
	Create(ctx context.Context, id, ownerID int64, params *servicecontract.CreateSnippetCommand) (*servicecontract.SnippetResult, error)
	GetByID(ctx context.Context, id int64) (*servicecontract.SnippetResult, error)
	ListByOwner(ctx context.Context, ownerID int64) ([]servicecontract.SnippetResult, error)
	Update(ctx context.Context, ownerID, id int64, params *servicecontract.UpdateSnippetCommand) (*servicecontract.SnippetResult, error)
	Delete(ctx context.Context, ownerID, id int64) error
}

type snippetRepository struct {
	client *ent.Client
}

// NewSnippetRepository 创建 SnippetRepository 实例。
func NewSnippetRepository(client *ent.Client) SnippetRepository {
	return &snippetRepository{client: client}
}

// Create 创建一条 Snippet 记录，ID 由外部（id-generator）传入。
func (r *snippetRepository) Create(ctx context.Context, id, ownerID int64, params *servicecontract.CreateSnippetCommand) (*servicecontract.SnippetResult, error) {
	builder := r.client.Snippet.Create().
		SetID(id).
		SetOwnerID(ownerID).
		SetType(resolveType(params.Type)).
		SetTitle(params.Title).
		SetLanguage(params.Language).
		SetVisibility(resolveVisibility(params.Visibility))

	// code/note 类型设置文本内容
	if params.Content != "" {
		builder.SetContent(params.Content)
	}

	// file 类型设置文件信息
	if params.FileURL != "" {
		builder.SetFileURL(params.FileURL)
	}
	if params.FileSize > 0 {
		builder.SetFileSize(params.FileSize)
	}
	if params.MimeType != "" {
		builder.SetMimeType(params.MimeType)
	}

	// 可选分组
	if params.GroupID != nil {
		builder.SetGroupID(*params.GroupID)
	}

	s, err := builder.Save(ctx)
	if err != nil {
		return nil, dberr.ParseDBError(err)
	}

	return toSnippetResult(s), nil
}

// GetByID 根据 ID 查询单个 Snippet。
func (r *snippetRepository) GetByID(ctx context.Context, id int64) (*servicecontract.SnippetResult, error) {
	s, err := r.client.Snippet.Get(ctx, id)
	if err != nil {
		return nil, dberr.ParseDBError(err)
	}

	return toSnippetResult(s), nil
}

// ListByOwner 按 owner_id 查询用户的所有 Snippet，按更新时间倒序。
func (r *snippetRepository) ListByOwner(ctx context.Context, ownerID int64) ([]servicecontract.SnippetResult, error) {
	snippets, err := r.client.Snippet.
		Query().
		Where(snippet.OwnerIDEQ(ownerID)).
		Order(snippet.ByUpdatedAt(sql.OrderDesc())).
		All(ctx)
	if err != nil {
		return nil, dberr.ParseDBError(err)
	}

	results := make([]servicecontract.SnippetResult, 0, len(snippets))
	for _, s := range snippets {
		results = append(results, *toSnippetResult(s))
	}

	return results, nil
}

// Update 更新指定 Snippet，需要校验 ownerID 所有权。
func (r *snippetRepository) Update(ctx context.Context, ownerID, id int64, params *servicecontract.UpdateSnippetCommand) (*servicecontract.SnippetResult, error) {
	// 先确认记录存在且属于该用户
	s, err := r.client.Snippet.
		Query().
		Where(snippet.IDEQ(id), snippet.OwnerIDEQ(ownerID)).
		Only(ctx)
	if err != nil {
		return nil, dberr.ParseDBError(err)
	}

	// 执行更新
	builder := s.Update().
		SetTitle(params.Title).
		SetContent(params.Content).
		SetLanguage(params.Language).
		SetVisibility(resolveVisibility(params.Visibility))

	if params.GroupID != nil {
		builder.SetGroupID(*params.GroupID)
	}

	updated, err := builder.Save(ctx)
	if err != nil {
		return nil, dberr.ParseDBError(err)
	}

	return toSnippetResult(updated), nil
}

// Delete 删除指定 Snippet，需要校验 ownerID 所有权。
func (r *snippetRepository) Delete(ctx context.Context, ownerID, id int64) error {
	// 先确认存在且属于用户
	count, err := r.client.Snippet.
		Query().
		Where(snippet.IDEQ(id), snippet.OwnerIDEQ(ownerID)).
		Count(ctx)
	if err != nil {
		return dberr.ParseDBError(err)
	}
	if count == 0 {
		return dberr.ErrNoRows
	}

	return r.client.Snippet.DeleteOneID(id).Exec(ctx)
}

// resolveType 将字符串转为 Ent 枚举类型。
func resolveType(t string) snippet.Type {
	switch t {
	case "note":
		return snippet.TypeNote
	case "file":
		return snippet.TypeFile
	default:
		return snippet.TypeCode
	}
}

// resolveVisibility 将字符串转为 Ent 枚举类型。
func resolveVisibility(visibility string) snippet.Visibility {
	if visibility == "public" {
		return snippet.VisibilityPublic
	}
	return snippet.VisibilityPrivate
}

// toSnippetResult 将 Ent 实体转为服务层结果结构体。
func toSnippetResult(s *ent.Snippet) *servicecontract.SnippetResult {
	result := &servicecontract.SnippetResult{
		ID:         s.ID,
		OwnerID:    s.OwnerID,
		Type:       string(s.Type),
		Title:      s.Title,
		Content:    s.Content,
		Language:   s.Language,
		Visibility: string(s.Visibility),
		GroupID:    s.GroupID,
		CreatedAt:  s.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  s.UpdatedAt.Format(time.RFC3339),
	}

	if s.FileURL != "" {
		result.FileURL = s.FileURL
	}
	if s.FileSize != 0 {
		result.FileSize = s.FileSize
	}
	if s.MimeType != "" {
		result.MimeType = s.MimeType
	}

	return result
}
