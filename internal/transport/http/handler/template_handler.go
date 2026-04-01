package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"github.com/luckysxx/go-note/internal/transport/http/response"
)

// TemplateHandler 处理模板相关的 HTTP 请求（双轨制 HTTP 层实现）。
type TemplateHandler struct {
	logger *zap.Logger
}

// NewTemplateHandler 创建 TemplateHandler 实例。
func NewTemplateHandler(logger *zap.Logger) *TemplateHandler {
	return &TemplateHandler{logger: logger}
}

// List 获取系统或个人可用的公共代码片段模板。
func (h *TemplateHandler) List(c *gin.Context) {
	// TODO: 未完成 - 返回预置或官方公共模板列表
	response.Success(c, []map[string]any{
		{
			"id":       "tpl-1",
			"name":     "Go HTTP Server",
			"language": "go",
		},
	})
}

// Get 获取单个模板的内容和详细信息。
func (h *TemplateHandler) Get(c *gin.Context) {
	// TODO: 未完成 - 根据 ID 查询完整的模板内容
	id := c.Param("id")
	response.Success(c, map[string]any{
		"id":       id,
		"name":     "Go HTTP Server",
		"content":  "package main\n\nimport \"net/http\"\n...",
		"language": "go",
	})
}
