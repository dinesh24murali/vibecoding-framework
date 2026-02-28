package plans

import (
	"net/http"
	"strconv"

	apperrors "github.com/dinesh/vibecoding-framework/backend/internal/errors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

type createPlanRequest struct {
	ProviderID  int64   `json:"provider_id" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
	Discount    float64 `json:"discount" binding:"required"`
	IsActive    bool    `json:"is_active"`
}

type updatePlanRequest struct {
	ProviderID  *int64   `json:"provider_id"`
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price"`
	Discount    *float64 `json:"discount"`
	IsActive    *bool    `json:"is_active"`
}

type paginatedPlansResponse struct {
	Items    []Plan `json:"items"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Total    int64  `json:"total"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ListAdminPlans(c *gin.Context) {
	page := parseInt(c.Query("page"), 1)
	pageSize := parseInt(c.Query("page_size"), 20)
	sort := c.Query("sort")

	providerID := parseOptionalID(c.Query("provider_id"))
	isActive := parseOptionalBool(c.Query("is_active"))

	items, total, err := h.service.ListAdmin(c.Request.Context(), page, pageSize, providerID, isActive, sort)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}

	c.JSON(http.StatusOK, paginatedPlansResponse{Items: items, Page: page, PageSize: pageSize, Total: total})
}

func (h *Handler) CreatePlan(c *gin.Context) {
	var req createPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "invalid request payload")
		return
	}

	plan, err := h.service.Create(c.Request.Context(), CreateInput{
		ProviderID:  req.ProviderID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Discount:    req.Discount,
		IsActive:    req.IsActive,
	})
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, plan)
}

func (h *Handler) GetPlanByID(c *gin.Context) {
	id, ok := parseID(c.Param("planId"))
	if !ok {
		respondError(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "invalid plan id")
		return
	}

	plan, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, plan)
}

func (h *Handler) UpdatePlan(c *gin.Context) {
	id, ok := parseID(c.Param("planId"))
	if !ok {
		respondError(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "invalid plan id")
		return
	}

	var req updatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "invalid request payload")
		return
	}

	if req.ProviderID == nil && req.Name == nil && req.Description == nil && req.Price == nil && req.Discount == nil && req.IsActive == nil {
		respondError(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "at least one field is required")
		return
	}

	plan, err := h.service.Update(c.Request.Context(), id, UpdateInput{
		ProviderID:  req.ProviderID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Discount:    req.Discount,
		IsActive:    req.IsActive,
	})
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, plan)
}

func (h *Handler) DeletePlan(c *gin.Context) {
	id, ok := parseID(c.Param("planId"))
	if !ok {
		respondError(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "invalid plan id")
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) ListCustomerPlansByProvider(c *gin.Context) {
	providerID, ok := parseID(c.Param("providerId"))
	if !ok {
		respondError(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "invalid provider id")
		return
	}

	items, err := h.service.ListCustomerByProvider(c.Request.Context(), providerID)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *Handler) handleServiceError(c *gin.Context, err error) {
	switch err {
	case ErrPlanValidation:
		respondError(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "invalid request payload")
	case ErrPlanDuplicate:
		respondError(c, http.StatusConflict, "CONFLICT_DUPLICATE", "plan already exists")
	case ErrProviderNotFound, ErrPlanNotFound:
		respondError(c, http.StatusNotFound, "RESOURCE_NOT_FOUND", "resource not found")
	case ErrPlanLimitExceeded:
		respondError(c, http.StatusConflict, "CONFLICT_DUPLICATE", "plan limit exceeded for provider")
	default:
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
	}
}

func respondError(c *gin.Context, status int, code string, message string) {
	requestID, _ := c.Get("request_id")
	c.JSON(status, apperrors.ErrorResponse{
		Error: apperrors.ErrorBody{
			Code:      code,
			Message:   message,
			RequestID: requestIDString(requestID),
		},
	})
}

func parseInt(value string, fallback int) int {
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func parseID(value string) (int64, bool) {
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil || id <= 0 {
		return 0, false
	}

	return id, true
}

func parseOptionalID(value string) *int64 {
	if value == "" {
		return nil
	}
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil || id <= 0 {
		return nil
	}

	return &id
}

func parseOptionalBool(value string) *bool {
	if value == "" {
		return nil
	}

	b, err := strconv.ParseBool(value)
	if err != nil {
		return nil
	}

	return &b
}

func requestIDString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}

	return ""
}
