package bench

import "github.com/rcrowley/go-metrics"

var r metrics.Registry
var totalRequest float64
var errRequest float64
var latency metrics.Timer

func prepareMetrics() {
	r = metrics.NewRegistry()
	latency = metrics.GetOrRegisterTimer("latency", r)
}
