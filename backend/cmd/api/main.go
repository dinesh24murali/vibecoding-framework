package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dinesh/vibecoding-framework/backend/internal/auth"
	"github.com/dinesh/vibecoding-framework/backend/internal/config"
	"github.com/dinesh/vibecoding-framework/backend/internal/db"
	"github.com/dinesh/vibecoding-framework/backend/internal/middleware"
	"github.com/dinesh/vibecoding-framework/backend/internal/providers"
	"github.com/dinesh/vibecoding-framework/backend/internal/router"
	"github.com/dinesh/vibecoding-framework/backend/internal/users"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	gormDB, err := db.OpenGorm(cfg.DB)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatalf("extract sql db: %v", err)
	}
	defer func() {
		if closeErr := sqlDB.Close(); closeErr != nil {
			log.Printf("close db error: %v", closeErr)
		}
	}()

	userRepo := users.NewRepository(gormDB)
	tokenManager := auth.NewTokenManager(
		cfg.JWT.AccessSecret,
		cfg.JWT.RefreshSecret,
		cfg.JWT.AccessTTL,
		cfg.JWT.RefreshTTL,
	)
	authService := auth.NewService(userRepo, tokenManager)
	authHandler := auth.NewHandler(authService)
	providerRepo := providers.NewRepository(gormDB)
	providerService := providers.NewService(providerRepo)
	providerHandler := providers.NewHandler(providerService)
	adminAuthMiddleware := middleware.AdminAuth(func(token string) (string, string, error) {
		claims, verifyErr := tokenManager.VerifyAccessToken(token)
		if verifyErr != nil {
			return "", "", verifyErr
		}

		return claims.RegisteredClaims.Subject, claims.Role, nil
	})

	engine, err := router.New(cfg.APIBasePath, cfg.AppEnv, cfg.OpenAPISpecPath, router.Dependencies{
		AuthHandler:         authHandler,
		ProvidersHandler:    providerHandler,
		AdminAuthMiddleware: adminAuthMiddleware,
	})
	if err != nil {
		log.Fatalf("init router: %v", err)
	}

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.APIPort),
		Handler:           engine,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("api starting env=%s port=%d", cfg.AppEnv, cfg.APIPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("api server failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("api shutdown error: %v", err)
	}
}
