package servicerequests

import (
	"errors"
	"net/http"
	"strconv"

	apperrors "github.com/dinesh/vibecoding-framework/backend/internal/errors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

type statusUpdateRequest struct {
	Status string `json:"status" binding:"required"`
}

type paginatedResponse struct {
	Items    []ServiceRequestResponse `json:"items"`
	Page     int                      `json:"page"`
	PageSize int                      `json:"page_size"`
	Total    int64                    `json:"total"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ListServiceRequests(c *gin.Context) {
	page := parseInt(c.Query("page"), 1)
	pageSize := parseInt(c.Query("page_size"), 20)
	providerID := parseOptionalID(c.Query("provider_id"))

	items, total, err := h.service.ListAdmin(c.Request.Context(), ListAdminInput{
		Page:       page,
		PageSize:   pageSize,
		Status:     c.Query("status"),
		ProviderID: providerID,
		FromDate:   c.Query("from_date"),
		ToDate:     c.Query("to_date"),
		Sort:       c.Query("sort"),
	})
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, paginatedResponse{Items: items, Page: page, PageSize: pageSize, Total: total})
}

func (h *Handler) GetServiceRequestByID(c *gin.Context) {
	id, ok := parseID(c.Param("serviceRequestId"))
	if !ok {
		respondError(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "invalid service request id")
		return
	}

	item, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *Handler) UpdateServiceRequestStatus(c *gin.Context) {
	id, ok := parseID(c.Param("serviceRequestId"))
	if !ok {
		respondError(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "invalid service request id")
		return
	}

	var req statusUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "invalid request payload")
		return
	}

	item, err := h.service.UpdateStatus(c.Request.Context(), id, req.Status)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *Handler) handleServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrServiceRequestNotFound):
		respondError(c, http.StatusNotFound, "RESOURCE_NOT_FOUND", "service request not found")
	case errors.Is(err, ErrInvalidTransition):
		respondError(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "invalid state transition")
	case errors.Is(err, ErrServiceRequestInvalid):
		respondError(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "invalid request payload")
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

func requestIDString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}

	return ""
}
