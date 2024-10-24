package workerpool

type WorkerOption func(pool *WorkerPool)

// WithTotalWorkers sets the total number of workers in the pool.
func WithTotalWorkers(workers int) WorkerOption {
	return func(pool *WorkerPool) {
		if workers < 1 {
			return
		}

		pool.totalWorkers = workers
	}
}

// WithMaxQueueLength sets the maximum number of tasks that can be scheduled.
func WithMaxQueueLength(length int) WorkerOption {
	return func(pool *WorkerPool) {
		if length < 1 {
			return
		}

		pool.maxQueueLength = length
	}
}

// WithDelayedStart sets the flag to indicate if the worker pool should start immediately.
//
// If delayedStart is set to true, the worker pool will start only when the first task is scheduled.
// This is useful when you want to start the worker pool only when you have tasks to schedule.
func WithDelayedStart(delayed bool) WorkerOption {
	return func(pool *WorkerPool) {
		pool.delayedStart = delayed
	}
}
