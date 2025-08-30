package pool

import (
	"errors"
	"sync"
)

var (
	ErrQueueFull   = errors.New("queue is full")
	ErrPoolStopped = errors.New("pool is stopped")
)

type Pool interface {
	Submit(task func()) error
	Stop() error
}

type pool struct {
	tasks   chan func()
	wg      sync.WaitGroup
	hook    func()
	mu      sync.Mutex
	stopped bool
}

func New(numWorkers, queueSize int, hook func()) Pool {
	if numWorkers <= 0 {
		numWorkers = 1
	}

	if queueSize < 0 {
		queueSize = 0
	}

	p := &pool{
		tasks: make(chan func(), queueSize),
		hook:  hook,
	}

	for i := 0; i < numWorkers; i++ {
		go p.worker()
	}

	return p
}

func (p *pool) worker() {
	for task := range p.tasks {
		p.executeTaskSafe(task)
	}
}

func (p *pool) executeTaskSafe(task func()) {
	defer p.wg.Done()
	defer func() {
		if r := recover(); r != nil {
			// task panicked - keep worker alive
		}
	}()

	task()

	if p.hook != nil {
		defer func() {
			if r := recover(); r != nil {
				// hook panicked - ignore to avoid crashing worker
			}
		}()

		p.hook()
	}
}

func (p *pool) Submit(task func()) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.stopped {
		return ErrPoolStopped
	}

	p.wg.Add(1)

	select {
	case p.tasks <- task:
		return nil
	default:
		p.wg.Done()
		return ErrQueueFull
	}
}

func (p *pool) Stop() error {
	p.mu.Lock()
	if !p.stopped {
		p.stopped = true
		close(p.tasks)
	}
	p.mu.Unlock()

	p.wg.Wait()
	return nil
}
