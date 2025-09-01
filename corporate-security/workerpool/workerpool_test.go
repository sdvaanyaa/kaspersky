package workerpool

import (
	"sync/atomic"
	"testing"
	"time"
)

// TestStopWaitExecutesAllTasks проверяет, что StopWait выполнит все задачи
func TestStopWaitExecutesAllTasks(t *testing.T) {
	wp := New(3)

	var counter int32
	for i := 0; i < 10; i++ {
		wp.Submit(func() {
			time.Sleep(5 * time.Millisecond)
			atomic.AddInt32(&counter, 1)
		})
	}

	wp.StopWait()

	if counter != 10 {
		t.Fatalf("expected 10 tasks executed, got %d", counter)
	}
}

// TestStopSkipsPendingTasks проверяет, что Stop выполняет только уже взятые задачи,
// а задачи, оставшиеся в очереди, пропускаются
func TestStopSkipsPendingTasks(t *testing.T) {
	wp := New(2)

	var counter int32
	for i := 0; i < 10; i++ {
		wp.Submit(func() {
			time.Sleep(20 * time.Millisecond)
			atomic.AddInt32(&counter, 1)
		})
	}

	wp.Stop()

	if counter >= 10 {
		t.Fatalf("expected some tasks skipped, got %d", counter)
	}
}

// TestSubmitWaitBlocksUntilDone проверяет, что SubmitWait блокируется, пока задача не завершится
func TestSubmitWaitBlocksUntilDone(t *testing.T) {
	wp := New(2)

	var x int32
	wp.SubmitWait(func() {
		time.Sleep(20 * time.Millisecond)
		atomic.StoreInt32(&x, 42)
	})

	if x != 42 {
		t.Fatalf("expected x=42, got %d", x)
	}

	wp.StopWait()
}

// TestPanicInTaskDoesNotKillWorker проверяет, что паника в задаче
// не останавливает воркер, и следующие задачи выполняются
func TestPanicInTaskDoesNotKillWorker(t *testing.T) {
	wp := New(2)

	var counter int32

	wp.Submit(func() {
		panic("boom")
	})

	wp.Submit(func() {
		atomic.AddInt32(&counter, 1)
	})

	wp.StopWait()

	if counter != 1 {
		t.Fatalf("expected surviving task to run, got %d", counter)
	}
}

// TestSubmitAfterStopIgnored проверяет, что Submit игнорируется,
// если пул уже остановлен
func TestSubmitAfterStopIgnored(t *testing.T) {
	wp := New(1)
	wp.Stop()

	var counter int32
	wp.Submit(func() {
		atomic.AddInt32(&counter, 1)
	})

	time.Sleep(50 * time.Millisecond)

	if counter != 0 {
		t.Fatalf("expected task ignored after stop, got %d", counter)
	}
}

// TestManyTasksConcurrent проверяет, что большое количество задач
// корректно выполняется пулом
func TestManyTasksConcurrent(t *testing.T) {
	wp := New(8)

	var counter int32
	numTasks := 1000

	for i := 0; i < numTasks; i++ {
		wp.Submit(func() {
			atomic.AddInt32(&counter, 1)
		})
	}

	wp.StopWait()

	if counter != int32(numTasks) {
		t.Fatalf("expected %d tasks executed, got %d", numTasks, counter)
	}
}
