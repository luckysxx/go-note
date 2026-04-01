package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"github.com/luckysxx/go-note/internal/transport/http/response"
)

// UploadHandler 处理附件上传转换为草稿的 HTTP 请求（双轨制 HTTP 层实现）。
type UploadHandler struct {
	logger *zap.Logger
}

// NewUploadHandler 创建 UploadHandler 实例。
func NewUploadHandler(logger *zap.Logger) *UploadHandler {
	return &UploadHandler{logger: logger}
}

// Upload 接收文件上传并将其转化为草稿片段。
func (h *UploadHandler) Upload(c *gin.Context) {
	// TODO: 未完成 - 解析 multipart/form-data 并对接对象存储 (S3/OSS)
	response.Success(c, map[string]any{
		"id":       "temp-1234",
		"filename": "uploaded.go",
		"size":     1024,
	})
}
