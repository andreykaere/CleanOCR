package pool

import (
	"log"
	"sync"
	"sync/atomic"
)

type Job struct {
	task  func() error
	state string
}

type workerPool struct {
	workersCount atomic.Int64
	workersMax   int
	buffer       chan Job
	workersWG    sync.WaitGroup
	closed       bool
	stopFlag     chan bool
	mutex        sync.Mutex
}

func newWorkerPool(workersMax, bufferSize int) *workerPool {
	return &workerPool{
		workersMax: workersMax,
		buffer:     make(chan Job, bufferSize),
		stopFlag:   make(chan bool),
	}
}

func (wp *workerPool) handle(job *Job) {
	job.state = "running"
	if err := job.task(); err == nil {
		job.state = "done"
		return
	}
	job.state = "failed"
}

func (wp *workerPool) workerDone() {
	wp.workersCount.Add(-1)
	wp.workersWG.Done()
}

func (wp *workerPool) workerReadBuffer() bool {
	select {
	case job, ok := <-wp.buffer:
		if ok {
			wp.handle(&job)
			return true
		}
	default:
	}
	return false
}

func (wp *workerPool) worker(job *Job) {
	defer wp.workerDone()
	wp.handle(job)

	// Handle queued tasks until stopped
	for {
		select {
		case _, ok := <-wp.stopFlag:
			if !ok {
				return
			}
		default:
			if !wp.workerReadBuffer() {
				return
			}
		}
	}
}

func (wp *workerPool) Submit(task func() error) {
	if wp.closed {
		log.Println("Worker pool is closed for adding new tasks")
		return
	}
	job := Job{task: task}
	wp.mutex.Lock()

	// Check if able to start new worker or queue
	if int(wp.workersCount.Load()) < wp.workersMax {
		wp.workersCount.Add(1)
		wp.workersWG.Add(1)
		wp.mutex.Unlock()
		go wp.worker(&job)
	} else {
		wp.buffer <- job
		wp.mutex.Unlock()
		job.state = "queued"
	}
}

// Submit and wait for task to finish
func (wp *workerPool) SubmitWait(task func() error) {
	flag := make(chan struct{})
	wp.Submit(func() error {
		defer close(flag)
		return task()
	})

	<-flag
}

// Stop accepting new tasks and wait for running and queued jobs to finish
func (wp *workerPool) StopWait() {
	if wp.closed {
		panic("Can't stop WorkerPool, it's already stopped")
	}
	wp.closed = true
	wp.workersWG.Wait()
}

// Stop accepting new tasks and wait for running jobs to finish
func (wp *workerPool) Stop() {
	if wp.closed {
		panic("Can't stop WorkerPool, it's already stopped")
	}
	wp.closed = true
	close(wp.stopFlag)
	wp.workersWG.Wait()
}

func (wp *workerPool) Closed() bool {
	return wp.closed
}

type Pool interface {
	Submit(task func() error)
	SubmitWait(task func() error)
	Stop()
	StopWait()
	Closed() bool
}

func NewPool(workersMax, bufferSize int) *workerPool {
	return newWorkerPool(workersMax, bufferSize)
}
