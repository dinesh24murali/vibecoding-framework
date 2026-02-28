package router

import (
	"fmt"
	"net/http"

	"github.com/dinesh/vibecoding-framework/backend/internal/auth"
	"github.com/dinesh/vibecoding-framework/backend/internal/http/contract"
	"github.com/dinesh/vibecoding-framework/backend/internal/middleware"
	"github.com/dinesh/vibecoding-framework/backend/internal/providers"
	"github.com/gin-gonic/gin"
)

type Dependencies struct {
	AuthHandler         *auth.Handler
	ProvidersHandler    *providers.Handler
	AdminAuthMiddleware gin.HandlerFunc
}

func New(basePath string, appEnv string, openAPISpecPath string, deps Dependencies) (*gin.Engine, error) {
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

	if deps.AuthHandler != nil {
		api.POST("/auth/admin/login", deps.AuthHandler.AdminLogin)
	}

	if deps.ProvidersHandler != nil {
		api.GET("/customer/providers", deps.ProvidersHandler.ListCustomerProviders)
	}

	admin := api.Group("/admin")
	if deps.AdminAuthMiddleware != nil {
		admin.Use(deps.AdminAuthMiddleware)
	}
	admin.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	if deps.ProvidersHandler != nil {
		admin.GET("/providers", deps.ProvidersHandler.ListAdminProviders)
		admin.POST("/providers", deps.ProvidersHandler.CreateProvider)
		admin.GET("/providers/:providerId", deps.ProvidersHandler.GetProviderByID)
		admin.PATCH("/providers/:providerId", deps.ProvidersHandler.UpdateProvider)
		admin.DELETE("/providers/:providerId", deps.ProvidersHandler.DeleteProvider)
	}

	return engine, nil
}
