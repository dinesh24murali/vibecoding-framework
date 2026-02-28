package router

import (
	"net/http"

	"github.com/dinesh/vibecoding-framework/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func New(basePath string) *gin.Engine {
	engine := gin.New()
	engine.Use(middleware.RequestID())
	engine.Use(gin.Logger())
	engine.Use(middleware.Recovery())

	engine.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := engine.Group(basePath)
	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	return engine
}
