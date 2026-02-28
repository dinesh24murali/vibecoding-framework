package middleware

import (
	"net/http"
	"strings"

	apperrors "github.com/dinesh/vibecoding-framework/backend/internal/errors"
	"github.com/gin-gonic/gin"
)

type AccessTokenVerifier interface {
	VerifyAccessToken(token string) (any, error)
}

func AdminAuth(verify func(token string) (subject string, role string, err error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			abortUnauthorized(c)
			return
		}

		token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		subject, role, err := verify(token)
		if err != nil || subject == "" || role != "admin" {
			abortUnauthorized(c)
			return
		}

		c.Set("admin_user_id", subject)
		c.Set("admin_role", role)
		c.Next()
	}
}

func abortUnauthorized(c *gin.Context) {
	requestID, _ := c.Get("request_id")
	c.AbortWithStatusJSON(http.StatusUnauthorized, apperrors.ErrorResponse{
		Error: apperrors.ErrorBody{
			Code:      "AUTH_UNAUTHORIZED",
			Message:   "unauthorized",
			RequestID: requestIDString(requestID),
		},
	})
}

func requestIDString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}

	return ""
}
