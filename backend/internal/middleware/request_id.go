package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	apperrors "github.com/dinesh/vibecoding-framework/backend/internal/errors"
	"github.com/gin-gonic/gin"
)

const RequestIDHeader = "X-Request-ID"

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(RequestIDHeader)
		if requestID == "" {
			requestID = randomID()
		}

		c.Set("request_id", requestID)
		c.Writer.Header().Set(RequestIDHeader, requestID)
		c.Next()
	}
}

func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, _ any) {
		c.AbortWithStatusJSON(http.StatusInternalServerError, apperrors.NewInternal(GetRequestID(c)))
	})
}

func GetRequestID(c *gin.Context) string {
	requestID, exists := c.Get("request_id")
	if !exists {
		return ""
	}

	value, ok := requestID.(string)
	if !ok {
		return ""
	}

	return value
}

func randomID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "00000000000000000000000000000000"
	}

	return hex.EncodeToString(b)
}
