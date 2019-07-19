package executor

import (
	"github.com/prometheus/client_golang/prometheus"
)

type runCollector struct {
	totalActiveRunCountDesc  *prometheus.Desc
	activeRunCountByTaskDesc *prometheus.Desc
	manualRunsCount          *prometheus.Desc
	te                       *TaskExecutor
}

// NewRunCollector returns a collector which exports influxdb process metrics.
func NewRunCollector(te *TaskExecutor) prometheus.Collector {
	return &runCollector{
		totalActiveRunCountDesc: prometheus.NewDesc(
			"task_total_runs_active",
			"Total number of runs across all tasks that have started but not yet completed",
			nil,
			prometheus.Labels{},
		),
		activeRunCountByTaskDesc: prometheus.NewDesc(
			"task_runs_active_by_task",
			"Number of runs for a given task that have started but not yet completed",
			[]string{"taskID"},
			nil,
		),
		manualRunsCount: prometheus.NewDesc(
			"manual_runs_task",
			"Total number of manual runs scheduled to run by taskID",
			[]string{"taskID"},
			nil,
		),
		te: te,
	}
}

// Describe returns all descriptions of the collector.
func (r *runCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- r.totalActiveRunCountDesc
}

// Collect returns the current state of all metrics of the collector.
func (r *runCollector) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(r.totalActiveRunCountDesc, prometheus.GaugeValue, r.te.WorkersBusy())

	ch <- prometheus.MustNewConstMetric(r.activeRunCountByTaskDesc, prometheus.GaugeValue, float64(r.te.tcs.CurrentlyRunning()))

	// ch <- prometheus.MustNewConstMetric(r.totalActiveRunCountDesc, prometheus.CounterValue, float64(r.te.ManualRuns()))

	// map of task id to run count
	// helper function in task executor that gives the RunCount() metrics function the info it needs
	// mb: given a task id, get all the active runs?
}
