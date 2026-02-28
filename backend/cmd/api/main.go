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

	"github.com/dinesh/vibecoding-framework/backend/internal/config"
	"github.com/dinesh/vibecoding-framework/backend/internal/router"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	engine, err := router.New(cfg.APIBasePath, cfg.AppEnv, cfg.OpenAPISpecPath)
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
