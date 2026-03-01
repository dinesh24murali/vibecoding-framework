package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	allowMethods = "GET,POST,PATCH,DELETE,OPTIONS"
	allowHeaders = "Authorization,Content-Type,Idempotency-Key,X-Request-ID"
)

func CORS(allowedOrigins []string) gin.HandlerFunc {
	originSet := make(map[string]struct{}, len(allowedOrigins))
	for _, origin := range allowedOrigins {
		if origin == "" {
			continue
		}
		originSet[origin] = struct{}{}
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" {
			if _, ok := originSet[origin]; ok {
				c.Header("Vary", "Origin")
				c.Header("Access-Control-Allow-Origin", origin)
				c.Header("Access-Control-Allow-Methods", allowMethods)
				c.Header("Access-Control-Allow-Headers", allowHeaders)
			}
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
