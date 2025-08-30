package workerpool

import (
	"errors"
	"github.com/sdvaanyaa/kasp/internal/config"
	"github.com/sdvaanyaa/kasp/internal/entity"
	"testing"
)

func TestEnqueue(t *testing.T) {
	wp := New(config.Config{Workers: 1, QueueSize: 1})

	err := wp.Enqueue(entity.Task{ID: "test"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	err = wp.Enqueue(entity.Task{ID: "overflow"})
	if !errors.Is(err, ErrQueueOverflow) {
		t.Errorf("got %v, want %v", err, ErrQueueOverflow)
	}
}

func TestProcessTask(t *testing.T) {
	wp := New(config.Config{Workers: 1, QueueSize: 1})

	tests := []struct {
		name      string
		task      entity.Task
		workFunc  func() error
		wantState string
	}{
		{
			name:      "success",
			task:      entity.Task{ID: "success", MaxRetries: 0},
			workFunc:  func() error { return nil },
			wantState: entity.TaskStateDone,
		},
		{
			name:      "fail no retry",
			task:      entity.Task{ID: "fail", MaxRetries: 0},
			workFunc:  func() error { return ErrSimulatedWork },
			wantState: entity.TaskStateFailed,
		},
		{
			name:      "fail with retry",
			task:      entity.Task{ID: "retry", MaxRetries: 1},
			workFunc:  func() error { return ErrSimulatedWork },
			wantState: entity.TaskStateRunning, // после первой итерации статус будет "running" и task пойдет на вторую итерацию
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wp.workFunc = tt.workFunc
			wp.processTask(tt.task)

			wp.stateMu.Lock()
			got := wp.taskStates[tt.task.ID].State
			wp.stateMu.Unlock()

			if got != tt.wantState {
				t.Errorf("for %s got %s, want %s", tt.task.ID, got, tt.wantState)
			}
		})
	}
}
