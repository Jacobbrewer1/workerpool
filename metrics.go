package workerpool

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type asyncMetricHandler struct {
	muxIdleWorkers   *sync.RWMutex
	muxActiveWorkers *sync.RWMutex
	muxPendingTasks  *sync.RWMutex
}

var (
	// idleWorkers is the number of idle workers in the pool.
	idleWorkers = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "workerpool_idle_workers",
		Help: "Number of idle workers in the worker pool",
	})

	// activeWorkers is the number of active workers in the pool.
	activeWorkers = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "workerpool_active_workers",
		Help: "Number of active workers in the worker pool",
	})

	// pendingTasks is the number of tasks pending in the worker pool.
	pendingTasks = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "workerpool_pending_tasks",
		Help: "Number of tasks pending in the worker pool",
	})

	// taskDuration is the duration of the task.
	taskDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "workerpool_task_duration_seconds",
		Help: "Duration of the task",
	})
)

func newAsyncMetricHandler() *asyncMetricHandler {
	return &asyncMetricHandler{
		muxIdleWorkers:   new(sync.RWMutex),
		muxActiveWorkers: new(sync.RWMutex),
		muxPendingTasks:  new(sync.RWMutex),
	}
}

func (h *asyncMetricHandler) incIdleWorkers() {
	h.muxIdleWorkers.Lock()
	idleWorkers.Inc()
	h.muxIdleWorkers.Unlock()
}

func (h *asyncMetricHandler) decIdleWorkers() {
	h.muxIdleWorkers.Lock()
	idleWorkers.Dec()
	h.muxIdleWorkers.Unlock()
}

func (h *asyncMetricHandler) setIdleWorkers(value float64) {
	h.muxIdleWorkers.Lock()
	idleWorkers.Set(value)
	h.muxIdleWorkers.Unlock()
}

func (h *asyncMetricHandler) incActiveWorkers() {
	h.muxActiveWorkers.Lock()
	activeWorkers.Inc()
	h.muxActiveWorkers.Unlock()
}

func (h *asyncMetricHandler) decActiveWorkers() {
	h.muxActiveWorkers.Lock()
	activeWorkers.Dec()
	h.muxActiveWorkers.Unlock()
}

func (h *asyncMetricHandler) setActiveWorkers(value float64) {
	h.muxActiveWorkers.Lock()
	activeWorkers.Set(value)
	h.muxActiveWorkers.Unlock()
}

func (h *asyncMetricHandler) incPendingTasks() {
	h.muxPendingTasks.Lock()
	pendingTasks.Inc()
	h.muxPendingTasks.Unlock()
}

func (h *asyncMetricHandler) decPendingTasks() {
	h.muxPendingTasks.Lock()
	pendingTasks.Dec()
	h.muxPendingTasks.Unlock()
}

func (h *asyncMetricHandler) setPendingTasks(value float64) {
	h.muxPendingTasks.Lock()
	pendingTasks.Set(value)
	h.muxPendingTasks.Unlock()
}
