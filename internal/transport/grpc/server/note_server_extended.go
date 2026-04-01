package grpcserver

import (
	"context"

	notepb "github.com/luckysxx/common/proto/note"
	grpcerrs "github.com/luckysxx/go-note/internal/transport/grpc/errs"
	"github.com/luckysxx/go-note/internal/transport/grpc/interceptor"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// =========================================================================
// 片段扩展能力 (Snippet Extensions)
// =========================================================================

// DeleteSnippet 删除指定的代码片段
func (s *NoteServer) DeleteSnippet(ctx context.Context, req *notepb.DeleteSnippetRequest) (*notepb.DeleteSnippetResponse, error) {
	userID, err := interceptor.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if err := s.snippetSvc.Delete(ctx, userID, req.SnippetId); err != nil {
		return nil, grpcerrs.ToStatusError(err)
	}

	return &notepb.DeleteSnippetResponse{Id: req.SnippetId}, nil
}

// SearchSnippets 根据查询条件搜索代码片段
func (s *NoteServer) SearchSnippets(ctx context.Context, req *notepb.SearchSnippetsRequest) (*notepb.ListSnippetsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "SearchSnippets 尚未实现")
}

// GetPublicSnippet 匿名获取公开的代码片段
func (s *NoteServer) GetPublicSnippet(ctx context.Context, req *notepb.GetPublicSnippetRequest) (*notepb.SnippetResponse, error) {
	return nil, status.Error(codes.Unimplemented, "GetPublicSnippet 尚未实现")
}

// FavoriteSnippet 收藏指定的代码片段
func (s *NoteServer) FavoriteSnippet(ctx context.Context, req *notepb.FavoriteSnippetRequest) (*notepb.FavoriteSnippetResponse, error) {
	return nil, status.Error(codes.Unimplemented, "FavoriteSnippet 尚未实现")
}

// UnfavoriteSnippet 取消收藏指定的代码片段
func (s *NoteServer) UnfavoriteSnippet(ctx context.Context, req *notepb.UnfavoriteSnippetRequest) (*notepb.FavoriteSnippetResponse, error) {
	return nil, status.Error(codes.Unimplemented, "UnfavoriteSnippet 尚未实现")
}

// CreateSnippetFromTemplate 基于已有模板创建新代码片段
func (s *NoteServer) CreateSnippetFromTemplate(ctx context.Context, req *notepb.CreateSnippetFromTemplateRequest) (*notepb.SnippetResponse, error) {
	return nil, status.Error(codes.Unimplemented, "CreateSnippetFromTemplate 尚未实现")
}

// =========================================================================
// 工作区列表能力 (Workspace Lists)
// =========================================================================

// ListRecentSnippets 获取用户最近访问的代码片段列表
func (s *NoteServer) ListRecentSnippets(ctx context.Context, req *notepb.ListSnippetsRequest) (*notepb.ListSnippetsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "ListRecentSnippets 尚未实现")
}

// ListSharedSnippets 获取与当前用户共享的代码片段列表
func (s *NoteServer) ListSharedSnippets(ctx context.Context, req *notepb.ListSnippetsRequest) (*notepb.ListSnippetsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "ListSharedSnippets 尚未实现")
}

// ListFavoriteSnippets 获取当前用户已经收藏的代码片段列表
func (s *NoteServer) ListFavoriteSnippets(ctx context.Context, req *notepb.ListSnippetsRequest) (*notepb.ListSnippetsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "ListFavoriteSnippets 尚未实现")
}

// =========================================================================
// 资源管理：分组 (Groups)
// =========================================================================

// ListGroups 获取当前用户拥有的所有分组列表
func (s *NoteServer) ListGroups(ctx context.Context, req *notepb.ListGroupsRequest) (*notepb.ListGroupsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "ListGroups 尚未实现")
}

// CreateGroup 创建一个全新的分组
func (s *NoteServer) CreateGroup(ctx context.Context, req *notepb.CreateGroupRequest) (*notepb.GroupResponse, error) {
	return nil, status.Error(codes.Unimplemented, "CreateGroup 尚未实现")
}

// UpdateGroup 更新现有分组的信息 (如重命名)
func (s *NoteServer) UpdateGroup(ctx context.Context, req *notepb.UpdateGroupRequest) (*notepb.GroupResponse, error) {
	return nil, status.Error(codes.Unimplemented, "UpdateGroup 尚未实现")
}

// DeleteGroup 删除指定的分组
func (s *NoteServer) DeleteGroup(ctx context.Context, req *notepb.DeleteGroupRequest) (*notepb.DeleteGroupResponse, error) {
	return nil, status.Error(codes.Unimplemented, "DeleteGroup 尚未实现")
}

// =========================================================================
// 资源管理：标签 (Tags)
// =========================================================================

// ListTags 获取当前用户创建的所有标签
func (s *NoteServer) ListTags(ctx context.Context, req *notepb.ListTagsRequest) (*notepb.ListTagsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "ListTags 尚未实现")
}

// CreateTag 为某代码片段创建一个新的个人标签
func (s *NoteServer) CreateTag(ctx context.Context, req *notepb.CreateTagRequest) (*notepb.TagResponse, error) {
	return nil, status.Error(codes.Unimplemented, "CreateTag 尚未实现")
}

// DeleteTag 删除某一个标签
func (s *NoteServer) DeleteTag(ctx context.Context, req *notepb.DeleteTagRequest) (*notepb.DeleteTagResponse, error) {
	return nil, status.Error(codes.Unimplemented, "DeleteTag 尚未实现")
}

// =========================================================================
// 模板与附件上传 (Templates & Uploads)
// =========================================================================

// ListTemplates 获取系统或个人可用的公共代码片段模板
func (s *NoteServer) ListTemplates(ctx context.Context, req *notepb.ListTemplatesRequest) (*notepb.ListTemplatesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "ListTemplates 尚未实现")
}

// GetTemplate 获取单个模板的内容和详细信息
func (s *NoteServer) GetTemplate(ctx context.Context, req *notepb.GetTemplateRequest) (*notepb.TemplateResponse, error) {
	return nil, status.Error(codes.Unimplemented, "GetTemplate 尚未实现")
}

// UploadFile 基于 gRPC 实现的文件字节流上传入口，通常将附件转为草稿片段
func (s *NoteServer) UploadFile(ctx context.Context, req *notepb.UploadFileRequest) (*notepb.UploadFileResponse, error) {
	return nil, status.Error(codes.Unimplemented, "UploadFile 尚未实现")
}
