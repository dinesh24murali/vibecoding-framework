package checkout

import (
	"errors"
	"net/http"

	apperrors "github.com/dinesh/vibecoding-framework/backend/internal/errors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

type createCheckoutRequest struct {
	ProviderID int64 `json:"provider_id" binding:"required"`
	PlanID     int64 `json:"plan_id" binding:"required"`
	Customer   struct {
		Name        string `json:"name" binding:"required"`
		PhoneNumber string `json:"phone_number" binding:"required"`
	} `json:"customer" binding:"required"`
	RecaptchaToken string `json:"recaptcha_token" binding:"required"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateServiceRequestAndPayment(c *gin.Context) {
	var req createCheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "invalid request payload")
		return
	}

	idempotencyKey := c.GetHeader("Idempotency-Key")
	result, err := h.service.CreateServiceRequestAndPayment(c.Request.Context(), CreateCheckoutInput{
		ProviderID:     req.ProviderID,
		PlanID:         req.PlanID,
		CustomerName:   req.Customer.Name,
		CustomerPhone:  req.Customer.PhoneNumber,
		RecaptchaToken: req.RecaptchaToken,
		IdempotencyKey: idempotencyKey,
	})
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidInput):
			respondError(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "invalid request payload")
		case errors.Is(err, ErrNotFound):
			respondError(c, http.StatusNotFound, "RESOURCE_NOT_FOUND", "provider or plan not found")
		case errors.Is(err, ErrConflict):
			respondError(c, http.StatusConflict, "CONFLICT_DUPLICATE", "conflicting customer record")
		case errors.Is(err, ErrCaptchaFailed):
			respondError(c, http.StatusUnprocessableEntity, "CAPTCHA_FAILED", "Please try after some time")
		default:
			respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		}

		return
	}

	c.JSON(http.StatusCreated, result)
}

func respondError(c *gin.Context, status int, code string, message string) {
	requestID, _ := c.Get("request_id")
	c.JSON(status, apperrors.ErrorResponse{
		Error: apperrors.ErrorBody{
			Code:      code,
			Message:   message,
			RequestID: toString(requestID),
		},
	})
}

func toString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}

	return ""
}
