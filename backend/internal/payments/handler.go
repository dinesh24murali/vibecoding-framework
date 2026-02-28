package payments

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	apperrors "github.com/dinesh/vibecoding-framework/backend/internal/errors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RetryPaymentForServiceRequest(c *gin.Context) {
	serviceRequestID, ok := parseID(c.Param("serviceRequestId"))
	if !ok {
		respondError(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "invalid service request id")
		return
	}

	var req struct {
		RecaptchaToken string `json:"recaptcha_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "invalid request payload")
		return
	}

	idempotencyKey := c.GetHeader("Idempotency-Key")
	result, err := h.service.RetryPayment(c.Request.Context(), serviceRequestID, req.RecaptchaToken, idempotencyKey)
	if err != nil {
		switch {
		case errors.Is(err, ErrServiceNotFound):
			respondError(c, http.StatusNotFound, "RESOURCE_NOT_FOUND", "service request not found")
		case errors.Is(err, ErrCaptchaRejected):
			respondError(c, http.StatusUnprocessableEntity, "CAPTCHA_FAILED", "Please try after some time")
		default:
			respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		}
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) GetCustomerPaymentStatus(c *gin.Context) {
	serviceRequestID, ok := parseID(c.Param("serviceRequestId"))
	if !ok {
		respondError(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "invalid service request id")
		return
	}

	phone := c.Query("phone_number")
	if phone == "" {
		respondError(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "phone_number is required")
		return
	}

	result, err := h.service.GetCustomerPaymentStatus(c.Request.Context(), serviceRequestID, phone)
	if err != nil {
		if errors.Is(err, ErrServiceNotFound) {
			respondError(c, http.StatusNotFound, "RESOURCE_NOT_FOUND", "service request not found")
			return
		}
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) HandleRazorpayCallback(c *gin.Context) {
	signature := c.GetHeader("X-Razorpay-Signature")
	if signature == "" {
		respondError(c, http.StatusBadRequest, "BAD_REQUEST", "missing X-Razorpay-Signature header")
		return
	}

	var body map[string]any
	raw, err := c.GetRawData()
	if err != nil {
		respondError(c, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
		return
	}

	if err := json.Unmarshal(raw, &body); err != nil {
		respondError(c, http.StatusBadRequest, "BAD_REQUEST", "invalid request payload")
		return
	}

	duplicate, err := h.service.HandleCallback(c.Request.Context(), raw, signature, body)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidSignature):
			respondError(c, http.StatusUnauthorized, "AUTH_UNAUTHORIZED", "invalid gateway signature")
		case errors.Is(err, ErrPaymentNotFound):
			respondError(c, http.StatusConflict, "CONFLICT_DUPLICATE", "duplicate callback safely ignored")
		default:
			respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		}
		return
	}

	if duplicate {
		c.JSON(http.StatusConflict, gin.H{"acknowledged": true, "duplicate": true})
		return
	}

	c.JSON(http.StatusOK, gin.H{"acknowledged": true})
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

func parseID(value string) (int64, bool) {
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil || id <= 0 {
		return 0, false
	}
	return id, true
}

func toString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
