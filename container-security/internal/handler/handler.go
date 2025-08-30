package handler

import (
	"encoding/json"
	"github.com/sdvaanyaa/kasp/internal/entity"
	"log/slog"
	"net/http"

	"github.com/sdvaanyaa/kasp/internal/workerpool"
)

func EnqueueHandler(wp *workerpool.WorkerPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var t entity.Task
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			slog.Error("failed to decode request body", slog.Any("error", err))
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		if t.ID == "" || t.MaxRetries < 0 {
			slog.Warn("invalid task data", slog.String("id", t.ID), slog.Int("max_retries", t.MaxRetries))
			http.Error(w, "invalid task: id required, max_retries >=0", http.StatusBadRequest)
			return
		}

		if err := wp.Enqueue(t); err != nil {
			slog.Error("failed to enqueue task", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("healthcheck requested")
	w.WriteHeader(http.StatusOK)
}
