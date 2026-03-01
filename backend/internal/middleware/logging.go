package middleware

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Logging() gin.HandlerFunc {
	return LoggingWithWriter(os.Stdout)
}

func LoggingWithWriter(w io.Writer) gin.HandlerFunc {
	logger := log.New(w, "", log.LstdFlags)

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		if query != "" {
			path = fmt.Sprintf("%s?%s", path, query)
		}

		errMsg := ""
		if len(c.Errors) > 0 {
			errMsg = strings.ReplaceAll(c.Errors.String(), "\n", " | ")
		}

		logger.Printf(
			"request_id=%s method=%s path=%q status=%d latency_ms=%d ip=%s errors=%q",
			GetRequestID(c),
			c.Request.Method,
			path,
			c.Writer.Status(),
			latency.Milliseconds(),
			c.ClientIP(),
			errMsg,
		)
	}
}
