package workerpool

import (
	"errors"
	"runtime"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	// ErrQueueFull is returned when the worker pool queue is full.
	ErrQueueFull = errors.New("worker pool queue is full")

	// ErrWorkerPoolStopped is returned when the worker pool is stopped.
	ErrWorkerPoolStopped = errors.New("worker pool stopped")
)

type WorkerPool struct {
	// wg is the wait group for the worker pool.
	wg *sync.WaitGroup

	// done is the channel to signal the worker pool to stop.
	done chan struct{}

	// tasks is the channel to send tasks to the worker pool.
	tasks chan Runnable

	// idleWorkers is the number of idle workers in the pool.
	idleWorkers uint32

	// activeWorkers is the number of active workers in the pool.
	activeWorkers uint32

	// delayedStart is the flag to indicate if the worker pool should start immediately.
	//
	// If delayedStart is set to true, the worker pool will start only when the first task is scheduled.
	// This is useful when you want to start the worker pool only when you have tasks to schedule.
	delayedStart bool

	// started is the flag to indicate if the worker pool has started.
	started bool

	// totalWorkers is the total number of workers to deploy to the pool.
	totalWorkers uint32

	// maxQueueLength is the maximum number of tasks that can be scheduled.
	maxQueueLength uint32
}

func NewWorkerPool(opts ...WorkerOption) *WorkerPool {
	pool := &WorkerPool{
		totalWorkers:   uint32(runtime.NumCPU()),
		maxQueueLength: uint32(runtime.NumCPU() * 1000),
		wg:             new(sync.WaitGroup),
		done:           make(chan struct{}),
		idleWorkers:    0,
		activeWorkers:  0,
		delayedStart:   false,
		started:        false,
	}

	for _, opt := range opts {
		opt(pool)
	}

	if pool.totalWorkers == 0 {
		pool.totalWorkers = uint32(runtime.NumCPU())
	}
	if pool.maxQueueLength == 0 {
		pool.maxQueueLength = pool.totalWorkers * 1000
	}

	pool.tasks = make(chan Runnable, pool.maxQueueLength)

	if !pool.delayedStart {
		pool.start()
	}

	return pool
}

func (p *WorkerPool) start() {
	p.started = true
	idleWorkers.Set(float64(p.totalWorkers))
	for i := uint32(0); i < p.totalWorkers; i++ {
		p.wg.Add(1)
		go p.deployWorker()
	}
}

func (p *WorkerPool) deployWorker() {
	defer p.wg.Done()
	for {
		select {
		case <-p.done:
			return
		case task, ok := <-p.tasks:
			if !ok {
				return
			}

			p.runTask(task)
		}
	}
}

func (p *WorkerPool) runTask(task Runnable) {
	idleWorkers.Dec()
	activeWorkers.Inc()
	defer func() {
		idleWorkers.Inc()
		activeWorkers.Dec()
	}()

	pendingTasks.Set(float64(len(p.tasks)))
	defer pendingTasks.Set(float64(len(p.tasks)))

	t := prometheus.NewTimer(taskDuration)
	defer t.ObserveDuration()

	task.Run()
}

func (p *WorkerPool) Stop() {
	close(p.done)
	p.wg.Wait()
}

func (p *WorkerPool) isDone() bool {
	select {
	case <-p.done:
		return true
	default:
		return false
	}
}

// Schedule schedules a task to be executed by the worker pool.
//
// Note: This is a non-blocking operation. If the worker pool is stopped, it will return an error.
func (p *WorkerPool) Schedule(task Runnable) error {
	if p.isDone() {
		return ErrWorkerPoolStopped
	}

	if p.delayedStart && !p.started {
		p.start()
	}

	select {
	case <-p.done:
		return ErrWorkerPoolStopped
	case p.tasks <- task:
		return nil
	default:
		return ErrQueueFull
	}
}

// BlockingSchedule schedules a task to be executed by the worker pool.
//
// Note: This is a blocking operation. If the worker pool is stopped, it will return an error.
func (p *WorkerPool) BlockingSchedule(task Runnable) error {
	if p.isDone() {
		return ErrWorkerPoolStopped
	}

	select {
	case <-p.done:
		return ErrWorkerPoolStopped
	case p.tasks <- task:
		return nil
	}
}
