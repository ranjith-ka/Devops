package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ranjith-ka/buildpilot/internal/ai/ollama"
	"github.com/ranjith-ka/buildpilot/internal/httpapi"
	"github.com/ranjith-ka/buildpilot/internal/store"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	address := env("HTTP_ADDRESS", ":8090")
	repositoryRoot := env("REPOSITORY_ROOT", ".")
	provider := ollama.New(
		env("OLLAMA_URL", "http://localhost:11434"),
		env("OLLAMA_MODEL", "qwen2.5-coder:7b"),
		2*time.Minute,
	)
	api, err := httpapi.New(logger, provider, store.NewMemory(), repositoryRoot)
	if err != nil {
		logger.Error("configure control plane", "error", err)
		os.Exit(1)
	}

	server := &http.Server{
		Addr:              address,
		Handler:           api.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      3 * time.Minute,
		IdleTimeout:       60 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		logger.Info("control plane listening", "address", address, "repository_root", repositoryRoot)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("control plane stopped unexpectedly", "error", err)
			os.Exit(1)
		}
	}()

	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("graceful shutdown", "error", err)
	}
}

func env(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
