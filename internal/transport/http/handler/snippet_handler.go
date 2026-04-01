package handler

import (
	"errors"
	"strconv"

	commonlogger "github.com/luckysxx/common/logger"
	"github.com/luckysxx/go-note/internal/dberr"
	"github.com/luckysxx/go-note/internal/service"
	servicecontract "github.com/luckysxx/go-note/internal/service/contract"
	httpdto "github.com/luckysxx/go-note/internal/transport/http/dto"
	httperrs "github.com/luckysxx/go-note/internal/transport/http/errs"
	"github.com/luckysxx/go-note/internal/transport/http/middleware"
	"github.com/luckysxx/go-note/internal/transport/http/response"
	"github.com/luckysxx/go-note/internal/transport/http/validator"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SnippetHandler 处理 snippet 相关的 HTTP 请求。
type SnippetHandler struct {
	svc    service.SnippetService
	logger *zap.Logger
}

// NewSnippetHandler 创建 SnippetHandler 实例。
func NewSnippetHandler(svc service.SnippetService, logger *zap.Logger) *SnippetHandler {
	return &SnippetHandler{svc: svc, logger: logger}
}

// Create 创建代码片段。
func (h *SnippetHandler) Create(c *gin.Context) {
	var req httpdto.CreateSnippetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errMsg := validator.TranslateValidationError(err)
		commonlogger.Ctx(c.Request.Context(), h.logger).Warn("参数验证失败", zap.Error(err), zap.String("message", errMsg))
		response.BadRequest(c, errMsg)
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "未登录")
		return
	}

	snippet, err := h.svc.Create(c.Request.Context(), userID, &servicecontract.CreateSnippetCommand{
		Type:       req.Type,
		Title:      req.Title,
		Content:    req.Content,
		FileURL:    req.FileURL,
		FileSize:   req.FileSize,
		MimeType:   req.MimeType,
		Language:   req.Language,
		Visibility: req.Visibility,
		GroupID:    req.GroupID,
	})
	if err != nil {
		commonlogger.Ctx(c.Request.Context(), h.logger).Error("创建 snippet 失败", zap.Error(err))

		response.Error(c, httperrs.ConvertToCustomError(err))
		return
	}

	response.Success(c, toSnippetResponse(snippet))
}

// ListMine 获取我的代码片段列表。
func (h *SnippetHandler) ListMine(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "未登录")
		return
	}

	list, err := h.svc.ListMine(c.Request.Context(), userID)
	if err != nil {
		commonlogger.Ctx(c.Request.Context(), h.logger).Error("查询用户 snippet 列表失败", zap.Error(err))
		response.Error(c, httperrs.ConvertToCustomError(err))
		return
	}

	results := make([]httpdto.SnippetResponse, 0, len(list))
	for _, snippet := range list {
		results = append(results, *toSnippetResponse(&snippet))
	}

	response.Success(c, results)
}

// Get 获取代码片段。
func (h *SnippetHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.BadRequest(c, "无效的片段ID")
		return
	}

	snippet, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		commonlogger.Ctx(c.Request.Context(), h.logger).Error("查询 snippet 失败", zap.Error(err))

		if dberr.IsNotFoundError(err) {
			response.NotFound(c, "Snippet 不存在")
			return
		}

		response.Error(c, httperrs.ConvertToCustomError(err))
		return
	}

	response.Success(c, toSnippetResponse(snippet))
}

// Update 更新代码片段。
func (h *SnippetHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.BadRequest(c, "无效的片段ID")
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "未登录")
		return
	}

	var req httpdto.UpdateSnippetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errMsg := validator.TranslateValidationError(err)
		commonlogger.Ctx(c.Request.Context(), h.logger).Warn("参数验证失败", zap.Error(err), zap.String("message", errMsg))
		response.BadRequest(c, errMsg)
		return
	}

	snippet, err := h.svc.Update(c.Request.Context(), userID, id, &servicecontract.UpdateSnippetCommand{
		Title:      req.Title,
		Content:    req.Content,
		Language:   req.Language,
		Visibility: req.Visibility,
		GroupID:    req.GroupID,
	})
	if err != nil {
		commonlogger.Ctx(c.Request.Context(), h.logger).Error("更新 snippet 失败", zap.Error(err))

		if errors.Is(err, service.ErrForbidden) {
			response.Forbidden(c, "无权限操作该代码片段")
			return
		}

		response.Error(c, httperrs.ConvertToCustomError(err))
		return
	}

	response.Success(c, toSnippetResponse(snippet))
}

// Delete 删除代码片段。
func (h *SnippetHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.BadRequest(c, "无效的片段ID")
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "未登录")
		return
	}

	if err := h.svc.Delete(c.Request.Context(), userID, id); err != nil {
		commonlogger.Ctx(c.Request.Context(), h.logger).Error("删除 snippet 失败", zap.Error(err))

		if errors.Is(err, service.ErrForbidden) {
			response.Forbidden(c, "无权限操作该代码片段")
			return
		}

		if dberr.IsNotFoundError(err) {
			response.NotFound(c, "Snippet 不存在")
			return
		}

		response.Error(c, httperrs.ConvertToCustomError(err))
		return
	}

	response.Success(c, map[string]any{"id": id})
}

// toSnippetResponse 将 service-layer result 转换为 HTTP response DTO。
func toSnippetResponse(s *servicecontract.SnippetResult) *httpdto.SnippetResponse {
	return &httpdto.SnippetResponse{
		ID:         s.ID,
		OwnerID:    s.OwnerID,
		Type:       s.Type,
		Title:      s.Title,
		Content:    s.Content,
		FileURL:    s.FileURL,
		FileSize:   s.FileSize,
		MimeType:   s.MimeType,
		Language:   s.Language,
		Visibility: s.Visibility,
		GroupID:    s.GroupID,
		CreatedAt:  s.CreatedAt,
		UpdatedAt:  s.UpdatedAt,
	}
}
