package router

import (
	"fmt"
	"net/http"

	"github.com/dinesh/vibecoding-framework/backend/internal/http/contract"
	"github.com/dinesh/vibecoding-framework/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func New(basePath string, appEnv string, openAPISpecPath string) (*gin.Engine, error) {
	validator, err := contract.New(appEnv, openAPISpecPath)
	if err != nil {
		return nil, fmt.Errorf("setup contract validator: %w", err)
	}

	engine := gin.New()
	engine.Use(middleware.RequestID())
	engine.Use(gin.Logger())
	engine.Use(middleware.Recovery())
	engine.Use(validator.Middleware())

	engine.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := engine.Group(basePath)
	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	return engine, nil
}
