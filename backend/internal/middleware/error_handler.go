package middleware

import (
	"net/http"

	apperrors "github.com/dinesh/vibecoding-framework/backend/internal/errors"
	"github.com/gin-gonic/gin"
)

type APIError struct {
	Status  int
	Code    string
	Message string
	Details []apperrors.ErrorDetail
}

func (e *APIError) Error() string {
	if e == nil {
		return ""
	}

	return e.Message
}

func NewAPIError(status int, code string, message string, details ...apperrors.ErrorDetail) *APIError {
	return &APIError{
		Status:  status,
		Code:    code,
		Message: message,
		Details: details,
	}
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if c.Writer.Written() || len(c.Errors) == 0 {
			return
		}

		status := http.StatusInternalServerError
		code := "INTERNAL_ERROR"
		message := "internal server error"
		var details []apperrors.ErrorDetail

		if candidate := c.Errors.Last(); candidate != nil {
			if apiErr, ok := candidate.Err.(*APIError); ok {
				if apiErr.Status > 0 {
					status = apiErr.Status
				}
				if apiErr.Code != "" {
					code = apiErr.Code
				}
				if apiErr.Message != "" {
					message = apiErr.Message
				}
				details = apiErr.Details
			}
		}

		c.JSON(status, apperrors.ErrorResponse{
			Error: apperrors.ErrorBody{
				Code:      code,
				Message:   message,
				Details:   details,
				RequestID: GetRequestID(c),
			},
		})
	}
}

func AbortWithAPIError(c *gin.Context, status int, code string, message string, details ...apperrors.ErrorDetail) {
	_ = c.Error(NewAPIError(status, code, message, details...))
	c.Abort()
}
