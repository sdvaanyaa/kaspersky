package workerpool

import (
	"sync"
)

const defaultQueueSize = 128

type WorkerPool struct {
	tasks    chan func()
	wg       sync.WaitGroup
	onceStop sync.Once
	mu       sync.Mutex
	stopped  bool
}

func New(numWorkers int) *WorkerPool {
	if numWorkers <= 0 {
		numWorkers = 1
	}

	wp := &WorkerPool{
		tasks: make(chan func(), defaultQueueSize),
	}

	for i := 0; i < numWorkers; i++ {
		go wp.worker()
	}

	return wp
}

func (wp *WorkerPool) worker() {
	for task := range wp.tasks {
		func() {
			defer func() {
				if r := recover(); r != nil {
					// task panicked - keep worker alive
				}
				wp.wg.Done()
			}()
			task()
		}()
	}
}

func (wp *WorkerPool) Submit(task func()) {
	wp.mu.Lock()
	if wp.stopped {
		wp.mu.Unlock()
		return
	}

	wp.wg.Add(1)
	wp.mu.Unlock()

	wp.tasks <- task
}

func (wp *WorkerPool) SubmitWait(task func()) {
	done := make(chan struct{})
	wp.Submit(func() {
		defer close(done)
		task()
	})
	<-done
}

func (wp *WorkerPool) Stop() {
	wp.onceStop.Do(func() {
		wp.mu.Lock()
		wp.stopped = true
		wp.mu.Unlock()
		close(wp.tasks)

		// clear queue
		for range wp.tasks {
			wp.wg.Done()
		}
	})
	wp.wg.Wait()
}

func (wp *WorkerPool) StopWait() {
	wp.onceStop.Do(func() {
		wp.mu.Lock()
		wp.stopped = true
		wp.mu.Unlock()
		close(wp.tasks)
	})
	wp.wg.Wait()
}
