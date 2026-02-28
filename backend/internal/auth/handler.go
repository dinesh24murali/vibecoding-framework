package auth

import (
	"errors"
	"net/http"

	apperrors "github.com/dinesh/vibecoding-framework/backend/internal/errors"
	"github.com/gin-gonic/gin"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type Handler struct {
	service *Service
}

type adminLoginRequest struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type adminLoginResponse struct {
	AccessToken      string    `json:"access_token"`
	TokenType        string    `json:"token_type"`
	ExpiresIn        int64     `json:"expires_in"`
	RefreshToken     string    `json:"refresh_token"`
	RefreshExpiresIn int64     `json:"refresh_expires_in"`
	User             userShape `json:"user"`
}

type userShape struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Role        string `json:"role"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) AdminLogin(c *gin.Context) {
	var req adminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "invalid request payload")
		return
	}

	result, err := h.service.LoginAdmin(c.Request.Context(), LoginInput{
		Identifier: req.Identifier,
		Password:   req.Password,
	})
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			respondError(c, http.StatusUnauthorized, "AUTH_INVALID_CREDENTIALS", "invalid credentials")
			return
		}

		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}

	c.JSON(http.StatusOK, adminLoginResponse{
		AccessToken:      result.AccessToken,
		TokenType:        result.TokenType,
		ExpiresIn:        result.ExpiresIn,
		RefreshToken:     result.RefreshToken,
		RefreshExpiresIn: result.RefreshExpiresIn,
		User: userShape{
			ID:          result.User.ID,
			Name:        result.User.Name,
			PhoneNumber: result.User.PhoneNumber,
			Role:        result.User.Role,
			Status:      result.User.Status,
			CreatedAt:   result.User.CreatedAt,
			UpdatedAt:   result.User.UpdatedAt,
		},
	})
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
