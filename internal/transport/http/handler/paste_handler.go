package handler

import (
	"errors"
	"strconv"

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

// PasteHandler 处理 paste 相关的 HTTP 请求
type PasteHandler struct {
	svc    service.PasteService
	logger *zap.Logger
}

// NewPasteHandler 创建 PasteHandler 实例
func NewPasteHandler(svc service.PasteService, logger *zap.Logger) *PasteHandler {
	return &PasteHandler{svc: svc, logger: logger}
}

// Create 创建代码片段
func (h *PasteHandler) Create(c *gin.Context) {
	var req httpdto.CreatePasteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errMsg := validator.TranslateValidationError(err)
		h.logger.Warn("参数验证失败", zap.Error(err), zap.String("message", errMsg))
		response.BadRequest(c, errMsg)
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "未登录")
		return
	}

	paste, err := h.svc.Create(c.Request.Context(), userID, &servicecontract.CreatePasteCommand{
		Title:      req.Title,
		Content:    req.Content,
		Language:   req.Language,
		Visibility: req.Visibility,
	})
	if err != nil {
		h.logger.Error("创建 paste 失败", zap.Error(err))

		if errors.Is(err, service.ErrShortLinkGeneration) {
			response.Error(c, httperrs.ConvertToCustomError(err))
			return
		}

		response.Error(c, httperrs.ConvertToCustomError(err))
		return
	}

	response.Success(c, toPasteResponse(paste))
}

// ListMine 获取我的代码片段列表
func (h *PasteHandler) ListMine(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "未登录")
		return
	}

	list, err := h.svc.ListMine(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("查询用户 paste 列表失败", zap.Error(err))
		response.Error(c, httperrs.ConvertToCustomError(err))
		return
	}

	results := make([]httpdto.PasteResponse, 0, len(list))
	for _, p := range list {
		results = append(results, *toPasteResponse(&p))
	}

	response.Success(c, results)
}

// Get 获取代码片段
func (h *PasteHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.BadRequest(c, "无效的片段ID")
		return
	}

	paste, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("查询 paste 失败", zap.Error(err))

		if dberr.IsNotFoundError(err) {
			response.NotFound(c, "Paste 不存在")
			return
		}

		response.Error(c, httperrs.ConvertToCustomError(err))
		return
	}

	response.Success(c, toPasteResponse(paste))
}

// Update 更新代码片段
func (h *PasteHandler) Update(c *gin.Context) {
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

	var req httpdto.UpdatePasteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errMsg := validator.TranslateValidationError(err)
		h.logger.Warn("参数验证失败", zap.Error(err), zap.String("message", errMsg))
		response.BadRequest(c, errMsg)
		return
	}

	paste, err := h.svc.Update(c.Request.Context(), userID, id, &servicecontract.UpdatePasteCommand{
		Title:      req.Title,
		Content:    req.Content,
		Language:   req.Language,
		Visibility: req.Visibility,
	})
	if err != nil {
		h.logger.Error("更新 paste 失败", zap.Error(err))

		if errors.Is(err, service.ErrForbidden) {
			response.Forbidden(c, "无权限操作该代码片段")
			return
		}

		response.Error(c, httperrs.ConvertToCustomError(err))
		return
	}

	response.Success(c, toPasteResponse(paste))
}

// toPasteResponse 将 service-layer result 转换为 HTTP response DTO
func toPasteResponse(p *servicecontract.PasteResult) *httpdto.PasteResponse {
	return &httpdto.PasteResponse{
		ID:         p.ID,
		OwnerID:    p.OwnerID,
		Title:      p.Title,
		ShortLink:  p.ShortLink,
		Content:    p.Content,
		Language:   p.Language,
		Visibility: p.Visibility,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
	}
}
