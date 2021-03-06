// Package mdata stands for "managed data" or "metrics data" if you will
// it has all the stuff to keep metric data in memory, store it, and synchronize
// save states over the network
package mdata

import "github.com/raintank/met"

var (
	LogLevel int

	chunkCreate met.Count
	chunkClear  met.Count

	metricsTooOld met.Count

	memToIterDuration met.Timer
	persistDuration   met.Timer

	metricsActive met.Gauge
	gcMetric      met.Count // metrics GC
)

func InitMetrics(stats met.Backend) {
	chunkCreate = stats.NewCount("chunks.create")
	chunkClear = stats.NewCount("chunks.clear")

	metricsTooOld = stats.NewCount("metrics_too_old")

	memToIterDuration = stats.NewTimer("mem.to_iter_duration", 0)
	persistDuration = stats.NewTimer("persist_duration", 0)

	gcMetric = stats.NewCount("gc_metric")
	metricsActive = stats.NewGauge("metrics_active", 0)
}
