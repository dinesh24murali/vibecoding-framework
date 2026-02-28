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
	"github.com/dinesh/vibecoding-framework/backend/internal/db"
	"github.com/dinesh/vibecoding-framework/backend/internal/notifications"
	"github.com/dinesh/vibecoding-framework/backend/internal/queue"
	"github.com/gin-gonic/gin"
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

	jobsQueue := queue.NewAsynq(gormDB)
	smsClient := notifications.NewLogSMSClient()
	smsWorker := notifications.NewWorker(jobsQueue, smsClient)

	engine := gin.New()
	engine.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.WorkerHealthPort),
		Handler:           engine,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("worker starting env=%s health_port=%d", cfg.AppEnv, cfg.WorkerHealthPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("worker health server failed: %v", err)
		}
	}()

	workerCtx, workerCancel := context.WithCancel(context.Background())
	defer workerCancel()
	go func() {
		if runErr := smsWorker.Run(workerCtx); runErr != nil {
			log.Printf("sms worker stopped with error: %v", runErr)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	workerCancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("worker shutdown error: %v", err)
	}
}
