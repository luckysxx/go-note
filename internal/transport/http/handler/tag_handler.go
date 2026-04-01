package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"github.com/luckysxx/go-note/internal/transport/http/response"
)

// TagHandler 处理标签相关的 HTTP 请求（双轨制 HTTP 层实现）。
type TagHandler struct {
	logger *zap.Logger
}

// NewTagHandler 创建 TagHandler 实例。
func NewTagHandler(logger *zap.Logger) *TagHandler {
	return &TagHandler{logger: logger}
}

// List 获取当前用户创建的所有标签。
func (h *TagHandler) List(c *gin.Context) {
	// TODO: 未完成 - 读取用户维度的标签字典
	response.Success(c, []map[string]any{})
}

// Create 创建一个新标签。
func (h *TagHandler) Create(c *gin.Context) {
	// TODO: 未完成 - 解析请求体并添加新标签
	response.Success(c, map[string]any{
		"id":   1,
		"name": "新标签",
	})
}

// Delete 删除指定标签。
func (h *TagHandler) Delete(c *gin.Context) {
	// TODO: 未完成 - 删除标签并清理桥接表的关联关系
	response.Success(c, map[string]any{
		"id": 1,
	})
}
