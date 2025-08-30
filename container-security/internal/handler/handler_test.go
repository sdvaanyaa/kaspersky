package handler

import (
	"bytes"
	"encoding/json"
	"github.com/sdvaanyaa/kaspersky/container-security/internal/config"
	"github.com/sdvaanyaa/kaspersky/container-security/internal/entity"
	"github.com/sdvaanyaa/kaspersky/container-security/internal/workerpool"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEnqueueHandler(t *testing.T) {
	wp := workerpool.New(config.Config{Workers: 1, QueueSize: 1})

	tests := []struct {
		name       string
		body       entity.Task
		wantStatus int
	}{
		{
			name:       "valid",
			body:       entity.Task{ID: "test", MaxRetries: 3},
			wantStatus: http.StatusAccepted,
		},
		{
			name:       "invalid id",
			body:       entity.Task{ID: "", MaxRetries: 3},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid retries",
			body:       entity.Task{ID: "test", MaxRetries: -1},
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/enqueue", bytes.NewReader(body))
			rr := httptest.NewRecorder()
			EnqueueHandler(wp)(rr, req)
			if rr.Code != tt.wantStatus {
				t.Errorf("got %d, want %d", rr.Code, tt.wantStatus)
			}
		})
	}
}

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()
	HealthHandler(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("got %d, want %d", rr.Code, http.StatusOK)
	}
}
