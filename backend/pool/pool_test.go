package pool

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestSubmit(t *testing.T) {
	var count atomic.Int64
	pool := NewPool(20, 100)
	taskCount := 10000

	for i := 0; i < taskCount; i++ {
		pool.Submit(func() error {
			time.Sleep(time.Millisecond)
			count.Add(1)
			return nil
		})
	}
	pool.StopWait()

	if x := int(count.Load()); x != taskCount {
		t.Errorf("Expected %d, got %d", taskCount, x)
	}
}

func TestStop(t *testing.T) {
	var count atomic.Int64
	pool := NewPool(30, 90)
	taskCount := 100

	for i := 0; i < taskCount; i++ {
		pool.Submit(func() error {
			time.Sleep(700 * time.Millisecond)
			count.Add(1)
			return nil
		})
	}
	pool.Stop()

	if x := int(count.Load()); x >= taskCount {
		t.Errorf("Expected %d, got %d", taskCount, x)
	}
}
