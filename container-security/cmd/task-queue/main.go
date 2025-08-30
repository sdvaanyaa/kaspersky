package main

import (
	"context"
	"errors"
	"github.com/sdvaanyaa/kaspersky/container-security/internal/config"
	"github.com/sdvaanyaa/kaspersky/container-security/internal/handler"
	"github.com/sdvaanyaa/kaspersky/container-security/internal/workerpool"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(log)

	cfg := config.LoadConfig()
	wp := workerpool.New(cfg)

	srv := newServer(wp)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", slog.Any("error", err))
		}
	}()

	waitForShutdown(srv, wp)
}

func newServer(wp *workerpool.WorkerPool) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/enqueue", handler.EnqueueHandler(wp))
	mux.HandleFunc("/healthz", handler.HealthHandler)

	return &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
}

func waitForShutdown(srv *http.Server, wp *workerpool.WorkerPool) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	slog.Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server shutdown error", slog.Any("error", err))
	}
	wp.Shutdown()

	slog.Info("shutdown complete")
}
