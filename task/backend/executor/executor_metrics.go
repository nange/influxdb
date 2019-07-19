package executor

import (
	"time"

	"github.com/influxdata/influxdb"
	"github.com/influxdata/influxdb/task/backend"
	"github.com/prometheus/client_golang/prometheus"
)

type ExecutorMetrics struct {
	totalRunsComplete *prometheus.CounterVec
	queueDelta        prometheus.Summary
	runDuration       prometheus.Summary
	errorsCounter     prometheus.Counter
	activeRunCount    prometheus.Collector
}

func NewExecutorMetrics(te *TaskExecutor) *ExecutorMetrics {
	const namespace = "task"
	const subsystem = "executor"

	return &ExecutorMetrics{
		totalRunsComplete: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "total_runs_complete",
			Help:      "Total number of runs completed across all tasks, split out by success or failure.",
		}, []string{"status"}),

		activeRunCount: NewRunCollector(te),

		queueDelta: prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace:  namespace,
			Subsystem:  subsystem,
			Name:       "run_queue_delta",
			Help:       "The duration in seconds between a run being due to start and actually starting.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		}),

		runDuration: prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace:  namespace,
			Subsystem:  subsystem,
			Name:       "run_duration",
			Help:       "The duration in seconds between a run starting and finishing.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		}),

		errorsCounter: prometheus.NewCounter(prometheus.CounterOpts{ // todo
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "errors_counter",
			Help:      "The number of errors that are thrown for a task.",
		}),
	}
}

// PrometheusCollectors satisfies the prom.PrometheusCollector interface.
func (em *ExecutorMetrics) PrometheusCollectors() []prometheus.Collector {
	return []prometheus.Collector{
		em.totalRunsComplete,
		em.activeRunCount,
		em.queueDelta,
		// return errors, duration, new runs active, etc.
	}
}

// any time prometheus asks, go and look at the length of the runs map for the task(s)
// instead of manually updating a counter with start/finish run
// except for queueDelta.Observe

// task executo with metrics function in the task executor
//

// StartRun store the delta time between when a run is due to start and actually starting.
func (em *ExecutorMetrics) StartRun(taskID influxdb.ID, queueDelta time.Duration) {
	em.queueDelta.Observe(queueDelta.Seconds())
}

// FinishRun adjusts the metrics to indicate a run is no longer in progress for the given task ID.
func (em *ExecutorMetrics) FinishRun(taskID influxdb.ID, status backend.RunStatus) { // add some more information
	em.totalRunsComplete.WithLabelValues(status.String()).Inc()
	// add another metric for how long it took to execute this run
	// em.
	// percentile of the slowest task runs
}
