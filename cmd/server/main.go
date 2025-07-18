package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pfg/internal/app"
	"pfg/internal/config"
	"pfg/internal/logger"

	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()

	logger, err := logger.Init(cfg.Production)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	appInstance, err := app.New(cfg, logger)
	if err != nil {
		logger.Fatal("Failed to initialize app", zap.Error(err))
	}

	go func() {
		if err := appInstance.Start(); err != nil && err.Error() != "http: Server closed" {
			logger.Fatal("Server failed", zap.Error(err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := appInstance.Shutdown(ctx); err != nil {
		logger.Fatal("Graceful shutdown failed", zap.Error(err))
	}

	logger.Info("Server gracefully stopped.")
}
