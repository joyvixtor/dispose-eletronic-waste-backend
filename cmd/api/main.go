package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/joyvixtor/dispose-eletronic-waste-backend/internal/app"
	"github.com/joyvixtor/dispose-eletronic-waste-backend/internal/config"
	"github.com/joyvixtor/dispose-eletronic-waste-backend/pkg/logger"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("warn, .env not found")
	}
	logger.Setup()

	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("failed to load config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	application, err := app.NewApp(cfg)
	if err != nil {
		slog.Error("failed to create app", slog.String("error", err.Error()))
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := application.Run(); err != nil && err != http.ErrServerClosed {
			slog.Error("failed to run app", slog.String("error", err.Error()))
			stop()
		}
	}()

	slog.Info("Aplicação Iniciada")

	<-ctx.Done()
	slog.Info("Shutting down gracefully, press Ctrl+C again to force")

	stop()
}
