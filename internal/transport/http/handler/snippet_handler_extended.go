package handler

import (
	"github.com/gin-gonic/gin"
	pkgerrs "github.com/luckysxx/common/errs"
	"github.com/luckysxx/go-note/internal/transport/http/response"
)

func extensionNotImplemented(c *gin.Context) {
	response.Error(c, pkgerrs.New(pkgerrs.ServerErr, "扩展接口暂未开放", nil))
}

// Search 搜索片段。
func (h *SnippetHandler) Search(c *gin.Context) {
	extensionNotImplemented(c)
}

// GetPublic 获取公开片段（不需要鉴权）。
func (h *SnippetHandler) GetPublic(c *gin.Context) {
	extensionNotImplemented(c)
}

// Favorite 收藏片段。
func (h *SnippetHandler) Favorite(c *gin.Context) {
	extensionNotImplemented(c)
}

// Unfavorite 取消收藏片段。
func (h *SnippetHandler) Unfavorite(c *gin.Context) {
	extensionNotImplemented(c)
}

// CreateFromTemplate 从模板创建片段。
func (h *SnippetHandler) CreateFromTemplate(c *gin.Context) {
	extensionNotImplemented(c)
}

// ListRecent 获取最近访问的片段列表。
func (h *SnippetHandler) ListRecent(c *gin.Context) {
	extensionNotImplemented(c)
}

// ListShared 获取与我共享的片段列表。
func (h *SnippetHandler) ListShared(c *gin.Context) {
	extensionNotImplemented(c)
}

// ListFavorites 获取我的收藏列表。
func (h *SnippetHandler) ListFavorites(c *gin.Context) {
	extensionNotImplemented(c)
}
