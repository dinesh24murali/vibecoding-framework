package providers

import (
	"context"
	"net/http"
	"strconv"

	csvimport "github.com/dinesh/vibecoding-framework/backend/internal/csv"
	apperrors "github.com/dinesh/vibecoding-framework/backend/internal/errors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service  *Service
	importer *csvimport.ProvidersImporter
}

type createProviderRequest struct {
	Name     string  `json:"name" binding:"required"`
	ImageURL *string `json:"image_url"`
}

type updateProviderRequest struct {
	Name     *string `json:"name"`
	ImageURL *string `json:"image_url"`
}

type paginatedProvidersResponse struct {
	Items    []Provider `json:"items"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
	Total    int64      `json:"total"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service:  service,
		importer: csvimport.NewProvidersImporter(providerImportJob{service: service}),
	}
}

func (h *Handler) ListAdminProviders(c *gin.Context) {
	page := parseInt(c.Query("page"), 1)
	pageSize := parseInt(c.Query("page_size"), 20)
	search := c.Query("search")
	sort := c.Query("sort")

	items, total, err := h.service.ListAdmin(c.Request.Context(), page, pageSize, search, sort)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}

	c.JSON(http.StatusOK, paginatedProvidersResponse{Items: items, Page: page, PageSize: pageSize, Total: total})
}

func (h *Handler) CreateProvider(c *gin.Context) {
	var req createProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "invalid request payload")
		return
	}

	provider, err := h.service.Create(c.Request.Context(), CreateInput{Name: req.Name, ImageURL: req.ImageURL})
	if err != nil {
		switch err {
		case ErrProviderValidation:
			respondError(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "invalid request payload")
		case ErrProviderDuplicate:
			respondError(c, http.StatusConflict, "CONFLICT_DUPLICATE", "provider already exists")
		default:
			respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		}

		return
	}

	c.JSON(http.StatusCreated, provider)
}

func (h *Handler) GetProviderByID(c *gin.Context) {
	id, ok := parseID(c.Param("providerId"))
	if !ok {
		respondError(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "invalid provider id")
		return
	}

	provider, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == ErrProviderNotFound {
			respondError(c, http.StatusNotFound, "RESOURCE_NOT_FOUND", "provider not found")
			return
		}

		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}

	c.JSON(http.StatusOK, provider)
}

func (h *Handler) UpdateProvider(c *gin.Context) {
	id, ok := parseID(c.Param("providerId"))
	if !ok {
		respondError(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "invalid provider id")
		return
	}

	var req updateProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "invalid request payload")
		return
	}

	if req.Name == nil && req.ImageURL == nil {
		respondError(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "at least one field is required")
		return
	}

	provider, err := h.service.Update(c.Request.Context(), id, UpdateInput{Name: req.Name, ImageURL: req.ImageURL})
	if err != nil {
		switch err {
		case ErrProviderValidation:
			respondError(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "invalid request payload")
		case ErrProviderDuplicate:
			respondError(c, http.StatusConflict, "CONFLICT_DUPLICATE", "provider already exists")
		case ErrProviderNotFound:
			respondError(c, http.StatusNotFound, "RESOURCE_NOT_FOUND", "provider not found")
		default:
			respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		}

		return
	}

	c.JSON(http.StatusOK, provider)
}

func (h *Handler) DeleteProvider(c *gin.Context) {
	id, ok := parseID(c.Param("providerId"))
	if !ok {
		respondError(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "invalid provider id")
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		if err == ErrProviderNotFound {
			respondError(c, http.StatusNotFound, "RESOURCE_NOT_FOUND", "provider not found")
			return
		}

		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) ListCustomerProviders(c *gin.Context) {
	items, err := h.service.ListCustomer(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *Handler) UploadProvidersCSV(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		respondError(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "invalid request payload")
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		respondError(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "invalid request payload")
		return
	}
	defer file.Close()

	result, err := h.importer.Import(c.Request.Context(), file)
	if err != nil {
		if err == csvimport.ErrInvalidCSVFile {
			respondError(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "invalid request payload")
			return
		}

		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}

	c.JSON(http.StatusOK, result)
}

type providerImportJob struct {
	service *Service
}

func (j providerImportJob) UpsertProvider(ctx context.Context, row csvimport.ProviderImportRow) error {
	return j.service.UpsertByName(ctx, row.Name, row.ImageURL)
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

func requestIDString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}

	return ""
}
