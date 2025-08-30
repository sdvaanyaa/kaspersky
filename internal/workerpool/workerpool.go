package workerpool

import (
	"context"
	"errors"
	"github.com/sdvaanyaa/kasp/internal/config"
	"github.com/sdvaanyaa/kasp/internal/entity"
	"log/slog"
	"math/rand"
	"sync"
	"time"
)

var (
	ErrSimulatedWork = errors.New("simulated error")
	ErrQueueOverflow = errors.New("queue overflow")
)

type WorkerPool struct {
	queue      chan entity.Task
	wg         sync.WaitGroup
	stateMu    sync.Mutex
	ctx        context.Context
	cancel     context.CancelFunc
	config     config.Config
	taskStates map[string]*entity.Task
}

func New(cfg config.Config) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	wp := &WorkerPool{
		queue:      make(chan entity.Task, cfg.QueueSize),
		ctx:        ctx,
		cancel:     cancel,
		config:     cfg,
		taskStates: make(map[string]*entity.Task),
	}

	slog.Info(
		"workerpool initialized",
		slog.Int("workers", cfg.Workers),
		slog.Int("queue_size", cfg.QueueSize),
	)

	for range cfg.Workers {
		wp.wg.Add(1)
		go wp.worker()
	}

	return wp
}

func (wp *WorkerPool) Enqueue(t entity.Task) error {
	select {
	case wp.queue <- t:
		wp.updateState(t.ID, entity.TaskStateQueued, &t)
		slog.Info("task enqueued", slog.String("task_id", t.ID))
		return nil
	default:
		slog.Warn("task queue is full", slog.String("task_id", t.ID))
		return ErrQueueOverflow
	}
}

func (wp *WorkerPool) Shutdown() {
	wp.cancel()
	close(wp.queue)
	wp.wg.Wait()
}

func (wp *WorkerPool) worker() {
	defer wp.wg.Done()

	for {
		select {
		case <-wp.ctx.Done():
			return
		case t, ok := <-wp.queue:
			if !ok {
				return
			}
			wp.processTask(t)
		}
	}
}

func (wp *WorkerPool) processTask(t entity.Task) {
	slog.Info(
		"task started",
		slog.String("task_id", t.ID),
		slog.Int("retries", t.Retries),
	)
	wp.updateState(t.ID, entity.TaskStateRunning, &t)

	err := wp.simulateWork()
	if err == nil {
		wp.updateState(t.ID, entity.TaskStateDone, &t)
		slog.Info("task completed", slog.String("task_id", t.ID))
		return
	}

	slog.Warn(
		"task processing failed, retrying",
		slog.String("task_id", t.ID),
		slog.Any("error", err),
	)

	t.Retries++
	if t.Retries > t.MaxRetries {
		wp.updateState(t.ID, entity.TaskStateFailed, &t)
		slog.Error("task failed permanently", slog.String("task_id", t.ID))
		return
	}

	time.Sleep(calcBackoff(t.Retries))
	wp.queue <- t
}

func calcBackoff(retries int) time.Duration {
	backoff := (time.Duration(1<<retries) * 100) * time.Millisecond
	jitter := time.Duration(rand.Intn(100)) * time.Millisecond
	return backoff + jitter
}

func (wp *WorkerPool) simulateWork() error {
	duration := time.Duration(100+rand.Intn(400)) * time.Millisecond
	time.Sleep(duration)

	if rand.Float64() < 0.2 {
		return ErrSimulatedWork
	}

	return nil
}

func (wp *WorkerPool) updateState(id, state string, t *entity.Task) {
	wp.stateMu.Lock()
	defer wp.stateMu.Unlock()

	t.State = state
	wp.taskStates[id] = t
}
