package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"github.com/luckysxx/go-note/internal/transport/http/response"
)

// GroupHandler 处理分组相关的 HTTP 请求（双轨制 HTTP 层实现）。
type GroupHandler struct {
	logger *zap.Logger
}

// NewGroupHandler 创建 GroupHandler 实例。
func NewGroupHandler(logger *zap.Logger) *GroupHandler {
	return &GroupHandler{logger: logger}
}

// List 获取当前用户的所有分组列表。
func (h *GroupHandler) List(c *gin.Context) {
	// TODO: 未完成 - 按 owner_id 查询实际分组数据
	response.Success(c, []map[string]any{})
}

// Create 创建一个新分组。
func (h *GroupHandler) Create(c *gin.Context) {
	// TODO: 未完成 - 解析请求体并持久化新分组
	response.Success(c, map[string]any{
		"id":   1,
		"name": "新分组",
	})
}

// Update 更新已有分组信息。
func (h *GroupHandler) Update(c *gin.Context) {
	// TODO: 未完成 - 校验所有权并更新分组名称
	response.Success(c, map[string]any{
		"id":   1,
		"name": "已更新分组",
	})
}

// Delete 删除指定分组。
func (h *GroupHandler) Delete(c *gin.Context) {
	// TODO: 未完成 - 删除分组并处理关联的片段归属
	response.Success(c, map[string]any{
		"id": 1,
	})
}
